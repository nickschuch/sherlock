package cmd

import (
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
	"github.com/pkg/errors"
	"github.com/google/uuid"

	"github.com/nickschuch/sherlock/storage"
	"github.com/nickschuch/sherlock/storage/types"
)

const (
	exampleEvents = `2018-01-02 07:08:14 +0000 UTC - Warning - MissingClusterDNS - kubelet does not have ClusterDNS IP configured and cannot create Pod using "ClusterFirst" policy. Falling back to DNSDefault policy.
2018-01-27 07:05:58 +0000 UTC - Normal - Killing - Killing container with id docker://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:pod "prod-3396783933-n4sls_app(xxxxxxxxxxxxxxxxxxxxxxxxx)" container "app" is unhealthy, it will be killed and re-created.
2018-01-27 07:05:58 +0000 UTC - Normal - Pulling - pulling image "example/app:1.1.8"
2018-01-27 07:06:01 +0000 UTC - Normal - Pulled - Successfully pulled image "example/app:1.1.8"
2018-01-27 07:06:01 +0000 UTC - Normal - Created - Created container with id xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
2018-01-27 07:06:01 +0000 UTC - Normal - Started - Started container with id xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
2018-01-27 07:07:58 +0000 UTC - Normal - Killing - Killing container with id docker://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:pod "prod-xxxxxxxxx-xxxxx_app(xxxxxxxxxxxxxxxxxxxxxxxxx)" container "app" is unhealthy, it will be killed and re-created.
2018-01-27 07:08:01 +0000 UTC - Normal - Created - Created container with id xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
2018-01-27 07:08:01 +0000 UTC - Normal - Started - Started container with id xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
2018-01-27 07:09:58 +0000 UTC - Normal - Killing - Killing container with id docker://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:pod "prod-xxxxxxxxx-xxxxx_app(xxxxxxxxxxxxxxxxxxxxxxxxx)" container "app" is unhealthy, it will be killed and re-created.
2018-01-27 07:10:03 +0000 UTC - Normal - Created - Created container with id xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
2018-01-27 07:10:03 +0000 UTC - Normal - Started - Started container with id xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
2018-01-27 07:11:58 +0000 UTC - Normal - Killing - Killing container with id docker://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:pod "prod-xxxxxxxxx-xxxxx_app(xxxxxxxxxxxxxxxxxxxxxxxxx)" container "app" is unhealthy, it will be killed and re-created.
2018-01-27 07:12:01 +0000 UTC - Normal - Created - Created container with id xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
2018-01-27 07:12:01 +0000 UTC - Normal - Started - Started container with id xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
2018-01-27 07:13:58 +0000 UTC - Normal - Killing - Killing container with id docker://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:pod "prod-xxxxxxxxx-xxxxx_app(xxxxxxxxxxxxxxxxxxxxxxxxx)" container "app" is unhealthy, it will be killed and re-created.
2018-01-27 07:14:01 +0000 UTC - Normal - Created - Created container with id xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
2018-01-27 07:14:01 +0000 UTC - Normal - Started - Started container with id xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
2018-01-27 07:15:58 +0000 UTC - Normal - Killing - Killing container with id docker://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:pod "prod-xxxxxxxxx-xxxxx_app(xxxxxxxxxxxxxxxxxxxxxxxxx)" container "app" is unhealthy, it will be killed and re-created.
2018-01-27 07:16:01 +0000 UTC - Normal - Created - Created container with id xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
2018-01-27 07:16:01 +0000 UTC - Normal - Started - Started container with id xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
	exampleLogs = `AH00558: apache2: Could not reliably determine the server's fully qualified domain name, using 241.0.173.30. Set the 'ServerName' directive globally to suppress this message
[ 2018-01-27 07:14:01.8145 14/7f13ea5a6780 age/Wat/WatchdogMain.cpp:1281 ]: Starting Passenger watchdog...
[ 2018-01-27 07:14:01.8315 17/7fd3bddcd780 age/Cor/CoreMain.cpp:1070 ]: Starting Passenger core...
[ 2018-01-27 07:14:01.8316 17/7fd3bddcd780 age/Cor/CoreMain.cpp:245 ]: Passenger core running in multi-application mode.
[ 2018-01-27 07:14:01.8359 17/7fd3bddcd780 age/Cor/CoreMain.cpp:820 ]: Passenger core online, PID 17
[ 2018-01-27 07:14:01.8460 25/7f9708334780 age/Ust/UstRouterMain.cpp:529 ]: Starting Passenger UstRouter...
[ 2018-01-27 07:14:01.8467 25/7f9708334780 age/Ust/UstRouterMain.cpp:342 ]: Passenger UstRouter online, PID 25
[ 2018-01-27 07:14:01.8490 25/7f97016f9700 age/Ust/UstRouterMain.cpp:422 ]: Signal received. Gracefully shutting down... (send signal 2 more time(s) to force shutdown)
[ 2018-01-27 07:14:01.8491 25/7f9708334780 age/Ust/UstRouterMain.cpp:492 ]: Received command to shutdown gracefully. Waiting until all clients have disconnected...
[ 2018-01-27 07:14:01.8491 25/7f97016f9700 Ser/Server.h:464 ]: [UstRouter] Shutdown finished`
	examplePod = `metadata:
  labels:
    env: prod
  name: prod-xxxxxxx-xxxx
  namespace: example
spec:
  containers:
  - name: app
    image: example/app:1.1.8
    imagePullPolicy: Always
    ports:
    - containerPort: 80
      protocol: TCP
status:
  conditions:
  - lastProbeTime: null
    lastTransitionTime: 2018-01-02T07:08:12Z
    status: "True"
    type: Initialized
  - lastProbeTime: null
    lastTransitionTime: 2018-01-02T07:08:19Z
    status: "True"
    type: Ready
  - lastProbeTime: null
    lastTransitionTime: 2018-01-02T07:08:12Z
    status: "True"
    type: PodScheduled
  containerStatuses:
  - containerID: docker://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    image: example/app:1.1.8
    imageID: docker-pullable://example/app3@sha256:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    lastState:
      terminated:
        containerID: docker://xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
        exitCode: 137
        finishedAt: 2018-01-27T07:13:58Z
        reason: Error
        startedAt: null
    name: app
    ready: true
    restartCount: 5
    state:
      running:
        startedAt: 2018-01-27T07:14:01Z
  hostIP: 10.0.0.10
  phase: Running
  podIP: 10.0.0.30
  qosClass: Burstable
  startTime: 2018-01-02T07:08:12Z`
)

type cmdDummy struct {
	storage string
	region   string
	bucket   string
}

func (cmd *cmdDummy) run(c *kingpin.ParseContext) error {
	client, err := storage.New(cmd.storage, cmd.region, cmd.bucket)
	if err != nil {
		return errors.Wrap(err, "failed to get client")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return errors.Wrap(err, "failed to generate UUID")
	}

	_, err = client.Put(types.PutParams{
		Incident: types.Incident{
			ID: id.String(),
			Created: time.Now(),
			Cluster: "example.com",
			Namespace: "test",
			Pod: "dev",
			Container: "app",
			Clues: []types.Clue{
				{
					Name: "EVENT",
					Content: exampleEvents,
				},
				{
					Name: "LOG",
					Content: exampleLogs,
				},
				{
					Name: "POD",
					Content: examplePod,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// Dummy declares the "dummy" sub command.
func Dummy(app *kingpin.Application) {
	c := new(cmdDummy)

	cmd := app.Command("dummy", "Create a dummy incident").Action(c.run)
	cmd.Flag("storage", "Type of storage which has our incidents").Default("s3").Envar("SHERLOCK_STORAGE").StringVar(&c.storage)
	cmd.Flag("region", "Region of the S3 bucket to store data").Default("ap-southeast-2").Envar("SHERLOCK_REGION").StringVar(&c.region)
	cmd.Flag("bucket", "Name of the S3 bucket to store data").Required().Envar("SHERLOCK_BUCKET").StringVar(&c.bucket)
}
