package cmd

import (
	"fmt"
	"time"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"github.com/pkg/errors"
	"github.com/google/uuid"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/nickschuch/sherlock/utils"
	"github.com/nickschuch/sherlock/storage"
	"github.com/nickschuch/sherlock/storage/types"
	"github.com/nickschuch/sherlock/utils/notification"
)

type cmdWatson struct {
	clusterName string
	lines int64
	storage string
	region   string
	bucket   string
	slackKey string
	slackChannel string
}

func (cmd *cmdWatson) run(c *kingpin.ParseContext) error {
	log.Info("Starting Watson")

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	log.Info("Starting to watch for Pod changes")

	watchlist := cache.NewListWatchFromClient(kubeClient.CoreV1().RESTClient(), "pods", corev1.NamespaceAll, fields.Everything())

	_, eController := cache.NewInformer(
		watchlist,
		&corev1.Pod{},
		time.Minute*30,
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				var (
					oldPod = oldObj.(*corev1.Pod)
					newPod = newObj.(*corev1.Pod)
				)

				config, err := rest.InClusterConfig()
				if err != nil {
					panic(err)
				}

				kubeClient, err := kubernetes.NewForConfig(config)
				if err != nil {
					panic(err)
				}

				if utils.IsIgnored(newPod) {
					log.Info(fmt.Sprintf("Skipping ignored pod %s/%s", newPod.ObjectMeta.Namespace, newPod.ObjectMeta.Name))
				} else {
					for _, newContainer := range newPod.Status.ContainerStatuses {
						oldContainer, err := utils.HasRestarts(oldPod.Status.ContainerStatuses, newContainer.Name)
						if err != nil {
							log.Errorf("Failed to get container restarts: %s", err)
							continue
						}

						if newContainer.RestartCount > oldContainer.RestartCount {
							go func(pod *corev1.Pod, container corev1.ContainerStatus) {
								err := push(kubeClient, pod, container, *cmd)
								if err != nil {
									log.Infof("Failed sending trace: %s", err)
								}
							}(oldPod, oldContainer)
						}
					}
				}
			},
		},
	)

	eController.Run(wait.NeverStop)

	return nil
}

// Helper function for pushing data to a S3 bucket.
func push(kubeClient *kubernetes.Clientset, pod *corev1.Pod, container corev1.ContainerStatus, params cmdWatson) error {
	client, err := storage.New(params.storage, params.region, params.bucket)
	if err != nil {
		return errors.Wrap(err, "failed to get storage client")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "failed to generate UUID")
	}

	log.With("id", id).Info("Requesting object")

	object, err := getObject(*pod)
	if err != nil {
		return errors.Wrap(err,"failed to get object")
	}

	log.With("id", id).Info("Requesting log")

	logs, err := getLogs(kubeClient, pod.Namespace, pod.Name, container.Name, params.lines)
	if err != nil {
		return errors.Wrap(err,"failed to get logs")
	}

	events, err := getEvents(kubeClient, pod)
	if err != nil {
		return errors.Wrap(err,"failed to get events")
	}

	_, err = client.Put(types.PutParams{
		Incident: types.Incident{
			ID: id.String(),
			Created: time.Now(),
			Cluster: params.clusterName,
			Namespace: pod.Namespace,
			Pod: pod.Name,
			Container: container.Name,
			Clues: []types.Clue{
				{
					Name: "EVENT",
					Content: events,
				},
				{
					Name: "LOG",
					Content: string(logs),
				},
				{
					Name: "POD",
					Content: string(object),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	// Check if we need to send a message to Slack about this incident.
	if params.slackKey != "" && params.slackChannel != "" {
		err := notification.Slack(notification.SlackParams{
			Key: params.slackKey,
			Channel: params.slackChannel,
			Bucket: params.bucket,
			Cluster: params.clusterName,
			Namespace: pod.Namespace,
			Pod: pod.Name,
			Container: container.Name,
			ID: id.String(),
		})
		if err != nil {
			return errors.Wrap(err, "failed to send notification")
		}
	}

	return nil
}

// Watson declares the "watson" sub command.
func Watson(app *kingpin.Application) {
	c := new(cmdWatson)

	cmd := app.Command("watson", "Watson runs on the cluster and stores all the clues").Action(c.run)
	cmd.Flag("cluster-name", "Cluster name to use for Slack notifications").Default("").Envar("CLUSTER_NAME").StringVar(&c.clusterName)
	cmd.Flag("log-lines", "Number of log lines to capture").Default("100").Envar("LOG_LINES").Int64Var(&c.lines)
	cmd.Flag("storage", "Type of storage which has our incidents").Default("s3").Envar("SHERLOCK_STORAGE").StringVar(&c.storage)
	cmd.Flag("region", "Region of the S3 bucket to store data").Default("ap-southeast-2").Envar("SHERLOCK_REGION").StringVar(&c.region)
	cmd.Flag("bucket", "Name of the S3 bucket to store data").Required().Envar("SHERLOCK_BUCKET").StringVar(&c.bucket)
	cmd.Flag("slack-key", "Slack key to use for authentication").Default("").Envar("SLACK_KEY").StringVar(&c.slackKey)
	cmd.Flag("slack-channel", "Slack channel to use for posting updates").Default("general").Envar("SLACK_CHANNEL").StringVar(&c.slackChannel)
}

func getEvents(kubeClient *kubernetes.Clientset, pod *corev1.Pod) (string, error) {
	list, err := kubeClient.CoreV1().Events(pod.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	var events []string

	for _, event := range list.Items {
		events = append(events, fmt.Sprintf("%s - %s - %s - %s", event.CreationTimestamp, event.Type, event.Reason, event.Message))
	}

	return strings.Join(events, "\n"), nil
}


func getLogs(kubeClient *kubernetes.Clientset, namespace, pod, container string, lines int64) ([]byte, error) {
	opts := &corev1.PodLogOptions{
		Container: container,
		Previous:  true,
		TailLines: &lines,
	}

	return kubeClient.CoreV1().Pods(namespace).GetLogs(pod, opts).DoRaw()
}

func getObject(pod corev1.Pod) ([]byte, error) {
	return yaml.Marshal(pod)
}
