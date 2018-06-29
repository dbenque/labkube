# Services

Now we know how to deploy our application. How should be target it?

The ```Service ``` resource is here to help. The ```Service``` resource contains the information relative to the port and protocol to use to consume the service inside a pod. To retrieve the eligible pods for a given service, we use a selector in the spec definition of the service

``` yaml
kind: Service
apiVersion: v1
metadata:
  name: labkubesvc
  labels:
    purpose: training
spec:
  selector:
    run: labkube
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

The selector is a set of ```key:value```, it is used by the enpoint controller to select all the pod which have a subset of labels that matches the selector.

The labels section of a pod is inside the metadata and comes from the template definition inside the deployment (and associated replicaSet). Note that all objects have metadata containing labels. The selector is only used against the labels of the pods, not the labels of the `deployment`.

```shell
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:                       # <-- Not these labels.
    run: labkube
  name: labkube
spec:
  replicas: 1
  selector:
    matchLabels:                # <-- Not these labels.
      run: labkube
  template:
    metadata:
      labels:                   # <-- YES these labels!
        run: labkube
        app: helloworld
    spec:
      containers:
      - image: dbenque/labkube:v1
        name: labkube
        ports:
        - containerPort: 8080  
```

| Match | Match | Match | Not Match | Not Match |
|-------|-------|-------|-------|-------|
|{run:labkube}|{app:helloworld}|{run:labkube, app:helloworld}|{run:other}|{target:prd, app:helloworld}|
|


Let's clean previous deployments and pods, and let's create new objects:
```shell
> kubectl delete deployments $(kubectl get deployments -ojsonpath={.items[*].metadata.name})
...
> kubectl delete pods $(kubectl get pods -ojsonpath={.items[*].metadata.name})
...

> kubectl get pods
No resources found.
```

Now let's create a brand new deployment that creates pods with the following labels:

```yaml
      labels:
        run: labkube
        instances: type1
```

Let's do that using file `deployment-1.yaml`:

```shell
> kubectl create -f deployment-1.yaml
deployment "labkube-1" created
```

This should trigger the creation of 2 pods:

``` shell
> kubectl get pods
NAME                         READY     STATUS    RESTARTS   AGE
labkube-1-556c647f86-24lgv   1/1       Running   0          24s
labkube-1-556c647f86-v28gr   1/1       Running   0          24s
```

Now let's create a service with the following selector:

```yaml
  selector:
    run: labkube
```

Let's do that using file `service-1.yaml`:
```yaml
kind: Service
apiVersion: v1
metadata:
  name: labkubesvc
  labels:
    purpose: training
spec:
  selector:
    run: labkube
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
```

```shell
> kubectl create -f service-1.yaml
service "labkubesvc" created
```

Now let's look at the resource that was effectively create:

```shell
> kubectl get service labkubesvc -oyaml
```
```yaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: 2018-06-29T09:27:09Z
  labels:
    purpose: training
  name: labkubesvc
  namespace: default
  resourceVersion: "56356"
  selfLink: /api/v1/namespaces/default/services/labkubesvc
  uid: 98ef90eb-7b7e-11e8-80ae-0800270f19f1
spec:
  clusterIP: 10.97.199.85
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
  selector:
    run: labkube
  sessionAffinity: None
  type: ClusterIP
status:
  loadBalancer: {}
```

Note that Kubernetes as associated a `clusterIP` to the service. This is the IP that we should target from inside the cluster to target the pod selected by the service.

Let's put ourselves inside the cluster with a shell, for that let's use a light image that as curl:

```shell
> kubectl run -i -t --rm shell --image=appropriate/curl --restart=Never --command -- /bin/sh
If you don't see a command prompt, try pressing enter.
/ #
```

Now that we are inside that we are inside the cluster let's target the service clusterIP

```shell
/ # curl 10.97.199.85/hello
Hello I am pod labkube-1-7b494bb868-9h75d
Welcome to kubernetes lab.
```

In another shell you can check the logs of the pod:

```shell
> kubectl logs labkube-1-7b494bb868-9h75d -f
2018/06/29 09:53:13 Server listening on port 8080...
2018/06/29 09:53:17 request on URI /ready
...
2018/06/29 09:53:20 request on URI /hello
...
```

As a client application I know that I need to target service named `labkubesvc` but I cannot predict its `clusterIP` that was dynamically assigned by Kubernetes. Should I call Kubernetes to know that `clusterIP`? That is a possibility but that would imply that your application is given the authorization to call Kubernetes API... that is not super cool: for security reasons, and also because applications will end up DDOS the api server of kubernetes. IN act Kubernetes as configured a dns inside the cluster for you to access the service via its name. Try the following:

```shell
/ # curl labkubesvc/hello
Hello I am pod labkube-1-7b494bb868-9h75d
Welcome to kubernetes lab.
```

Try to target the service multiple time using the curl command (with the clusterIP or the dnsname). You should notice that the traffic is loadbalanced between the 2 pods that are associated to the service:

```shell
/ # curl labkubesvc/hello
Hello I am pod labkube-1-7b494bb868-kp7jp
Welcome to kubernetes lab.
/ # curl labkubesvc/hello
Hello I am pod labkube-1-7b494bb868-9h75d
Welcome to kubernetes lab.
/ # curl labkubesvc/hello
Hello I am pod labkube-1-7b494bb868-kp7jp
Welcome to kubernetes lab.
/ # curl labkubesvc/hello
Hello I am pod labkube-1-7b494bb868-kp7jp
Welcome to kubernetes lab.
/ # curl labkubesvc/hello
Hello I am pod labkube-1-7b494bb868-9h75d
Welcome to kubernetes lab.
```

When you create a `Service` with a non empty `spec.Selector`, Kubernetes creates another object call `Endpoint` with the same name. Let's have a look at it:

```shell
> kubectl get endpoints
NAME         ENDPOINTS                         AGE
labkubesvc   172.17.0.6:8080,172.17.0.7:8080   31m

> kubectl get endpoints labkubesvc -oyaml
```
```yaml
apiVersion: v1
kind: Endpoints
metadata:
  creationTimestamp: 2018-06-29T09:53:27Z
  labels:
    purpose: training
  name: labkubesvc
  namespace: default
  resourceVersion: "57630"
  selfLink: /api/v1/namespaces/default/endpoints/labkubesvc
  uid: 457d84d1-7b82-11e8-80ae-0800270f19f1
subsets:
- addresses:
  - ip: 172.17.0.6
    nodeName: minikube
    targetRef:
      kind: Pod
      name: labkube-1-7b494bb868-kp7jp
      namespace: default
      resourceVersion: "57615"
      uid: 3c99770b-7b82-11e8-80ae-0800270f19f1
  - ip: 172.17.0.7
    nodeName: minikube
    targetRef:
      kind: Pod
      name: labkube-1-7b494bb868-9h75d
      namespace: default
      resourceVersion: "57618"
      uid: 3c9a03a2-7b82-11e8-80ae-0800270f19f1
  ports:
  - port: 8080
    protocol: TCP

```

This endpoint objects contains the list of IPs that matches the service selector. Remember the `ReadinessProbe` ? What if the service was not ready?

Let's modify the content of one pod to make it not ready. The code associated to the readiness probe of our example application check that the file "/ready" exists. Here is the code:
```go
if _, err := os.Stat("/ready"); os.IsNotExist(err) {
    w.WriteHeader(500)
    log.Printf("--> readiness probe failed")
    return
}
```

This empty file was added at in the image:
```docker
FROM busybox

RUN ["touch", "/ready"]
```

Let's 2 things now:
- Monitor the endpoints, in a dedicated shell:
```shell
> kubectl get endpoints -w
NAME         ENDPOINTS                         AGE
labkubesvc   172.17.0.6:8080,172.17.0.7:8080   46m
...
```
- Modify the pod content to change its readiness status:
```shell
> kubectl exec labkube-1-7b494bb868-hdfcx -- mv /ready /notready
```

Check how often the probe is running inside you pod; after that period you should see changes in the endpoints list once you have done the `/notready` modification:

```shell
> kubectl get pods -ojsonpath='{range .items[*]}{.metadata.name} : {.spec.containers[*].readinessProbe.periodSeconds}s{"\n"}{end}'
labkube-1-7b494bb868-hdfcx : 30s
labkube-1-7b494bb868-xjbgl : 30s
shell : s
```

When one of the pod is not ready check the endpoint resource content:

```shell
> kubectl get endpoints -oyaml
```
```yaml
> kubectl get endpoints labkubesvc -oyaml
apiVersion: v1
kind: Endpoints
metadata:
  creationTimestamp: 2018-06-29T09:53:27Z
  labels:
    purpose: training
  name: labkubesvc
subsets:
- addresses:
  - ip: 172.17.0.7
    nodeName: minikube
    targetRef:
      kind: Pod
      name: labkube-1-7b494bb868-xjbgl
      namespace: default
      resourceVersion: "59514"
      uid: f7ccf07c-7b87-11e8-80ae-0800270f19f1
  notReadyAddresses:                                    # <- New section was created
  - ip: 172.17.0.6
    nodeName: minikube
    targetRef:
      kind: Pod
      name: labkube-1-7b494bb868-hdfcx
      namespace: default
      resourceVersion: "60598"
      uid: f7c89a22-7b87-11e8-80ae-0800270f19f1
  ports:
  - port: 8080
    protocol: TCP
```

You should notice that a `notReadyAddresses` section was created in the resource. This means that the pod is still selected by the service, but no traffic will be sent to it because of its readiness status.

Note that you can also use the `describe` command to get informations about objects. Try:
```shell
> kubectl describe service labkubesvc
Name:              labkubesvc
Namespace:         default
Labels:            purpose=training
Annotations:       <none>
Selector:          run=labkube
Type:              ClusterIP
IP:                10.97.199.85
Port:              <unset>  80/TCP
TargetPort:        8080/TCP
Endpoints:         172.17.0.10:8080,172.17.0.6:8080
Session Affinity:  None
Events:            <none>


> kubectl describe endpoints labkubesvc
Name:         labkubesvc
Namespace:    default
Labels:       purpose=training
Annotations:  <none>
Subsets:
  Addresses:          172.17.0.10,172.17.0.6
  NotReadyAddresses:  <none>
  Ports:
    Name     Port  Protocol
    ----     ----  --------
    <unset>  8080  TCP

Events:  <none>

```

## Exercise 1
In case you have modified them during the lab, re-apply the definition of deployment-1 and the service.
```shell
> kubectl apply -f deployment-1.yaml
...
> kubectl apply -f service-1.yaml
...
```

Let's make some cleanup by deleting all pods. Kubernetes will recreate them with in initial status.

```shell
> kubectl delete pods $(kubectl get pods -ojsonpath={.items[*].metadata.name})
...
```

Let's now add a second deployment, using the file deployment-2.yaml:
```shell
> kubectl apply -f deployment-2.yaml
deployment "labkube-2" created
```


Now your service should target all the pods created by the 2 deployments. You can check this using the following curl command:
```shell
/ # curl labkubesvc/mydeployment
Hello I am pod labkube-1-7b494bb868-lwq4t
Welcome to kubernetes lab.
MY_DEPLOYMENT environment variable is set to: My deployment is labkube-1
/ # curl labkubesvc/mydeployment
Hello I am pod labkube-1-7b494bb868-lwq4t
Welcome to kubernetes lab.
MY_DEPLOYMENT environment variable is set to: My deployment is labkube-1
/ # curl labkubesvc/mydeployment
Hello I am pod labkube-2-7b74446546-2w7jm
Welcome to kubernetes lab.
MY_DEPLOYMENT environment variable is set to: My deployment is labkube-2
/ # curl labkubesvc/mydeployment
Hello I am pod labkube-2-7b74446546-2w7jm
Welcome to kubernetes lab.
MY_DEPLOYMENT environment variable is set to: My deployment is labkube-2
/ # curl labkubesvc/mydeployment
Hello I am pod labkube-1-7b494bb868-prnnj
Welcome to kubernetes lab.
MY_DEPLOYMENT environment variable is set to: My deployment is labkube-1
```

- Trainer: "Create a service that only target the pods of the second deployment"
- Trainer: "Show the endpoints of the new service and curl the new service to validate your setup"


## Exercise 2

- User: "Can we target directly the pod instead of using the service clusterIP (or dnsname)?"

- Trainer: "Yes you can by using the IP of the pod. That could be interesting for investigation or development purposes. You should not do that with your regular application. Let's do it for the fun:"

```shell
> kubectl get pods -owide
NAME                         READY     STATUS    RESTARTS   AGE       IP           NODE
labkube-1-7b494bb868-9h75d   1/1       Running   0          24m       172.17.0.7   minikube
labkube-1-7b494bb868-kp7jp   1/1       Running   0          24m       172.17.0.6   minikube
shell                        1/1       Running   0          22m       172.17.0.8   minikube
```

Let's take the IP and use it with the curl:

```shell
/ # curl 172.17.0.7/hello
curl: (7) Failed to connect to 172.17.0.7 port 80: Connection refused
```

- User: "It is not working! Why did you say it would work?"
- Trainer: "It is going to work: look at your service definition and fix your command to be able to directly target the pod."