package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/ausrasul/hashgen"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/nickschuch/sherlock/common/storage"
)

var (
	cliLogLines       = kingpin.Flag("log-lines", "Number of log lines to capture").Default("100").OverrideDefaultFromEnvar("LOG_LINES").Int64()
	cliS3Bucket       = kingpin.Flag("s3-bucket", "Name of the S3 bucket to store data").Default("").OverrideDefaultFromEnvar("S3_BUCKET").String()
	cliRegion         = kingpin.Flag("s3-region", "Region of the S3 bucket to store data").Default("ap-southeast-2").OverrideDefaultFromEnvar("S3_REGION").String()
	cliPrometheusPort = kingpin.Flag("prometheus-port", "Prometheus metrics port").Default(":9000").OverrideDefaultFromEnvar("METRICS_PORT").String()
	cliPrometheusPath = kingpin.Flag("prometheus-path", "Prometheus metrics path").Default("/metrics").OverrideDefaultFromEnvar("METRICS_PATH").String()
	cliSlackKey       = kingpin.Flag("slack-key", "Slack key to use for authentication").Default("").OverrideDefaultFromEnvar("SLACK_KEY").String()
	cliSlackChannel   = kingpin.Flag("slack-channel", "Slack channel to use for posting updates").Default("general").OverrideDefaultFromEnvar("SLACK_CHANNEL").String()
	cliClusterName    = kingpin.Flag("cluster-name", "Cluster name to use for Slack notifications").Default("").OverrideDefaultFromEnvar("CLUSTER_NAME").String()
)

func main() {
	kingpin.Parse()

	log.Info("Starting Watson")

	log.Info("Connecting to Kubernetes")

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	log.Info("Start Prometheus metrics")

	go metrics(*cliPrometheusPort, *cliPrometheusPath)

	log.Info("Watching for Pod changes")

	watchlist := cache.NewListWatchFromClient(kubeClient.CoreV1().RESTClient(), "pods", v1.NamespaceAll, fields.Everything())

	_, eController := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		time.Minute*30,
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: podUpdate,
		},
	)

	eController.Run(wait.NeverStop)
}

// Helper function for responding to pod updates.
func podUpdate(oldObj, newObj interface{}) {
	var (
		oldPod = oldObj.(*v1.Pod)
		newPod = newObj.(*v1.Pod)
	)

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	for _, newContainer := range newPod.Status.ContainerStatuses {
		oldContainer, err := restarts(oldPod.Status.ContainerStatuses, newContainer.Name)
		if err != nil {
			log.Errorf("Failed to get container restarts: %s", err)
			continue
		}

		if newContainer.RestartCount > oldContainer.RestartCount {
			go func(pod *v1.Pod, container v1.ContainerStatus) {
				err := push(kubeClient, pod, container)
				if err != nil {
					log.Infof("Failed sending trace: %s", err)
				}
			}(oldPod, oldContainer)
		}
	}
}

// Helper function for pushing data to a S3 bucket.
func push(kubeClient *kubernetes.Clientset, pod *v1.Pod, container v1.ContainerStatus) error {
	var (
		id       = fmt.Sprintf("%s/%s/%s", pod.Namespace, pod.Name, container.Name)
		store    = storage.New(*cliRegion, *cliS3Bucket)
		incident = hashgen.Get(24)
	)

	log.With("id", id).Info("Requesting object")

	meta, err := getObject(*pod)
	if err != nil {
		return fmt.Errorf("Failed to get object: %s", err)
	}

	log.With("id", id).Info("Sending object to storage")

	store.Write(pod.Namespace, pod.Name, container.Name, incident, fileObject, meta)
	if err != nil {
		panic(err)
	}

	log.With("id", id).Info("Requesting log")

	logs, err := getLogs(kubeClient, pod.Namespace, pod.Name, container.Name)
	if err != nil {
		return fmt.Errorf("Failed to get container logs: %s", err)
	}

	log.With("id", id).Info("Sending logs to storage")

	store.Write(pod.Namespace, pod.Name, container.Name, incident, fileLogs, logs)
	if err != nil {
		return fmt.Errorf("Failed to store container logs: %s", err)
	}

	log.With("id", id).Info("Requesting events")

	events, err := getEvents(kubeClient, pod)
	if err != nil {
		return fmt.Errorf("Failed to get container events: %s", err)
	}

	log.With("id", id).Info("Sending events to storage")

	store.Write(pod.Namespace, pod.Name, container.Name, incident, fileEvents, events)
	if err != nil {
		return fmt.Errorf("Failed to store container events: %s", err)
	}

	// Check if we need to send a message to Slack about this incident.
	if *cliSlackKey != "" && *cliSlackChannel != "" {
		err := notifySlack(*cliSlackKey, *cliSlackChannel, *cliS3Bucket, *cliClusterName, pod.Namespace, pod.Name, container.Name, incident)
		if err != nil {
			return fmt.Errorf("Failed to send Slack notification: %s", err)
		}
	}

	return nil
}

func metrics(port, path string) {
	http.Handle(path, promhttp.Handler())
	log.Fatal(http.ListenAndServe(port, nil))
}
