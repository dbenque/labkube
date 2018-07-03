# Pod

Deploy your first pod. Almost like ```docker run```, use ```kubectl run```

``` shell
kubectl run --help
```

let's run our first pod

``` shell
> kubectl run labkube --image=dbenque/labkube:v1 --port=8080 --restart=Never
pod "labkube" created
```

let's have a look at pods

``` shell
> kubectl get pods
NAME      READY     STATUS    RESTARTS   AGE
labkube   1/1       Running   0          3m
```

let's get some more details like the IP with extended view

``` shell
> kubectl get pods -owide
NAME      READY     STATUS    RESTARTS   AGE       IP           NODE
labkube   1/1       Running   0          4m        172.17.0.6   minikube
```

let's get some more information about the status

``` shell
> kubectl describe pod labkube
Name:         labkube
Namespace:    default
Node:         minikube/192.168.99.100
Start Time:   Thu, 28 Jun 2018 10:38:50 +0200
Labels:       run=labkube
Annotations:  <none>
Status:       Running
IP:           172.17.0.6
Containers:
  labkube:
    Container ID:   docker://bddbea28eb9e8fcbb3960e8e507137992586d5cd8099898a56050f4b3e3052c6
    Image:          dbenque/labkube:v1
    Image ID:       docker-pullable://dbenque/labkube@sha256:fb5fd5dd05c4f63d88eca049786e7b2627c8bf63cecf23f8b86f39adaeec28c3
    Port:           8080/TCP
    State:          Running
      Started:      Thu, 28 Jun 2018 10:38:51 +0200
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from default-token-vtgmb (ro)
Conditions:
  Type           Status
  Initialized    True
  Ready          True
  PodScheduled   True
Volumes:
  default-token-vtgmb:
    Type:        Secret (a volume populated by a Secret)
    SecretName:  default-token-vtgmb
    Optional:    false
QoS Class:       BestEffort
Node-Selectors:  <none>
Tolerations:     <none>
Events:
  Type    Reason                 Age   From               Message
  ----    ------                 ----  ----               -------
  Normal  Scheduled              7m    default-scheduler  Successfully assigned labkube to minikube
  Normal  SuccessfulMountVolume  7m    kubelet, minikube  MountVolume.SetUp succeeded for volume "default-token-vtgmb"
  Normal  Pulled                 7m    kubelet, minikube  Container image "dbenque/labkube:v1" already present on machine
  Normal  Created                7m    kubelet, minikube  Created container
  Normal  Started                7m    kubelet, minikube  Started container
```

what logs?

``` shell
kubectl logs labkube
```

can we enter the container?

``` shell
> kubectl exec -t -i labkube /bin/sh
/ # ls
bin      dev      etc      home     labkube  proc     root     sys      tmp      usr      var
/ # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /labkube
    8 root      0:00 /bin/sh
   15 root      0:00 ps
/ # exit
```

let's have a look at the pod definition

``` shell
> kubectl get pod -oyaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: 2018-06-28T08:38:50Z
  labels:
    run: labkube
  name: labkube
  namespace: default
  resourceVersion: "2828"
  selfLink: /api/v1/namespaces/default/pods/labkube
  uid: aeb023ca-7aae-11e8-80ae-0800270f19f1
spec:
  containers:
  - image: dbenque/labkube:v1
    imagePullPolicy: IfNotPresent
    name: labkube
    ports:
    - containerPort: 8080
      protocol: TCP
    resources: {}
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: default-token-vtgmb
      readOnly: true
  dnsPolicy: ClusterFirst
  nodeName: minikube
  restartPolicy: Never
  schedulerName: default-scheduler
  securityContext: {}
  serviceAccount: default
  serviceAccountName: default
  terminationGracePeriodSeconds: 30
  volumes:
  - name: default-token-vtgmb
    secret:
      defaultMode: 420
      secretName: default-token-vtgmb
status:
...
...
  hostIP: 192.168.99.100
  phase: Running
  podIP: 172.17.0.6
  qosClass: BestEffort
  startTime: 2018-06-28T08:38:50Z
```

let's delete our pod now

``` shell
> kubectl delete pod labkube
pod "labkube" deleted
```

Now let's play directly with pod object definition. Have a look at file pod-1.yaml.
All the values that have disapeared compare to what we have seen with ```kubectl get pod labkube -oyaml``` are either the default values or status or runtime value.

``` shell
> cat pod-1.yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: labkube
  name: labkube
spec:
  containers:
  - image: dbenque/labkube:v1
    name: labkube
    ports:
    - containerPort: 8080


> kubectl create -f pod-1.yaml
pod "labkube" created
```

Let's create another pod with a modified definition to introduce an environment variable.
Have a look at file pod-2.yaml and inject it.

``` shell
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: labkube
  name: labkube-env
spec:
  containers:
  - image: dbenque/labkube:v1
    env:
    - name: MY_LABKUBE_VAR
      value: "Hello from the environment"
    name: labkube
    ports:
    - containerPort: 8080


> kubectl create -f pod-1.yaml
pod "labkube-env" created
```

Now we should have 2 pods

``` shell
> kubectl get pods
NAME          READY     STATUS    RESTARTS   AGE
labkube       1/1       Running   0          15m
labkube-env   1/1       Running   0          39s
```

Let's have a look at their environment

``` shell
> kubectl exec labkube -- /bin/sh -c "cat /proc/1/environ | tr '\0' '\n'"
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=labkube
KUBERNETES_PORT_443_TCP_ADDR=10.96.0.1
KUBERNETES_SERVICE_HOST=10.96.0.1
KUBERNETES_SERVICE_PORT=443
KUBERNETES_SERVICE_PORT_HTTPS=443
KUBERNETES_PORT=tcp://10.96.0.1:443
KUBERNETES_PORT_443_TCP=tcp://10.96.0.1:443
KUBERNETES_PORT_443_TCP_PROTO=tcp
KUBERNETES_PORT_443_TCP_PORT=443
HOME=/root

> kubectl exec labkube-env -- /bin/sh -c "cat /proc/1/environ | tr '\0' '\n'"
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=labkube-env
MY_LABKUBE_VAR=Hello from the environment
KUBERNETES_PORT_443_TCP_PORT=443
KUBERNETES_PORT_443_TCP_ADDR=10.96.0.1
KUBERNETES_SERVICE_HOST=10.96.0.1
KUBERNETES_SERVICE_PORT=443
KUBERNETES_SERVICE_PORT_HTTPS=443
KUBERNETES_PORT=tcp://10.96.0.1:443
KUBERNETES_PORT_443_TCP=tcp://10.96.0.1:443
KUBERNETES_PORT_443_TCP_PROTO=tcp
HOME=/root
```

Also something noticeable here: the HOSTNAME environment variable is set to the name of the pod.

You can directly edit some sections of the resource. For exampe you can add an annotation. The following command will open the editor configured for your environment:

``` shell
> kubectl edit pod labkube
```
