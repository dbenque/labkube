# Deployment

Now that we know ```pods``` the following, we can ask the following questions:

- What if the pod that I have launched die?
- What if I want to run multiple instances of the same pod?
- What if I want to rollout a new version of my pod ?

The answer is ```deployment```

First delete all the pods inside you namespace:

``` shell
> kubectl delete pods $(kubectl get pods -ojsonpath={.items[*].metadata.name})
pod "labkube" deleted
pod "labkube-env" deleted
```

Let's use again the command ```kubectl run``` but with a different ```restartPolicy```

``` shell
> kubectl run labkube --image=dbenque/labkube:v1 --port=8080 --restart=Always
deployment "labkube" created
```

Note that this time a deployment was create, it is not a pod. Let's have a look to the definition of the object.

``` shell
> kubectl get deployment labkube -oyaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
  creationTimestamp: 2018-06-28T09:54:43Z
  generation: 1
  labels:
    run: labkube
  name: labkube
  namespace: default
  resourceVersion: "6373"
  selfLink: /apis/extensions/v1beta1/namespaces/default/deployments/labkube
  uid: 483b215a-7ab9-11e8-80ae-0800270f19f1
spec:
  replicas: 1
  selector:
    matchLabels:
      run: labkube
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        run: labkube
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
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
status:
  availableReplicas: 1
  conditions:
  - lastTransitionTime: 2018-06-28T09:54:43Z
    lastUpdateTime: 2018-06-28T09:54:43Z
    message: Deployment has minimum availability.
    reason: MinimumReplicasAvailable
    status: "True"
    type: Available
  observedGeneration: 1
  readyReplicas: 1
  replicas: 1
  updatedReplicas: 1
```

What about the pods? Pay attention to the name!

``` shell
> kubectl get pods
NAME                       READY     STATUS    RESTARTS   AGE
labkube-76c7c9754d-nnjpm   1/1       Running   0          10m
```

Let's have a look at the pod definition and lets focus on the field ```metadata.ownerReferences```

``` shell
> kubectl get pods -oyaml
apiVersion: v1
items:
- apiVersion: v1
  kind: Pod
  metadata:
    creationTimestamp: 2018-06-28T09:54:43Z
    generateName: labkube-76c7c9754d-
    labels:
      pod-template-hash: "3273753108"
      run: labkube
    name: labkube-76c7c9754d-nnjpm
    namespace: default
    ownerReferences:
    - apiVersion: extensions/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: ReplicaSet                              #Look on that line!
      name: labkube-76c7c9754d
      uid: 483bd32b-7ab9-11e8-80ae-0800270f19f1
    resourceVersion: "6371"
    selfLink: /api/v1/namespaces/default/pods/labkube-76c7c9754d-nnjpm
    uid: 483d57e7-7ab9-11e8-80ae-0800270f19f1
  spec:
    containers:
    ...
```

So the parent of the ```pod``` is not the ```deployment```, it is a replicaSet. Let's have a look at that ```replicaSet``` object.

``` shell
> kubectl get replicaset
NAME                 DESIRED   CURRENT   READY     AGE
labkube-76c7c9754d   1         1         1         17m
```

If you open the yaml defintion of the replicaset you will notice that it is really close to the deployment. It contains the number of replicas and the pod definition.

Purpose of the objects:

- The Pod run the container(s)
- The ReplicaSet object is going to be used by kubernetes to ensure that the relevant number of pods are running.
- The Deployment help to transition from one replicaSet to another (check the spec.strategy section in the definition)

Depending on the modification that you do on the Deployment object the existing replicaSet will be modified or a new one will be create. If a new one is create then the deployment strategy is applied to transition from one definition to another.

Let's delete the current deployment and create one with the resource definition in file ```deployment-1.yaml```

```shell
> kubectl delete deployment labkube
deployment "labkube" deleted

> kubectl create --record -f deployment-1.yaml
deployment "labkube" created
```

Open a new shell and run the following command to monitor what is happening at pod level:

``` shell
> kubectl get pods -w
NAME                       READY     STATUS    RESTARTS   AGE
labkube-76c7c9754d-ks88c   1/1       Running   0          1m
....
```

Open another new shell and run the following command to monitor what is happening at replicaSet level:
``` shell
> kubectl get replicaset -w
NAME                 DESIRED   CURRENT   READY     AGE
labkube-76c7c9754d   1         1         1         1m
....
```

Now let's change the number of replica to 2, by applying a new definition:
```
> kubectl apply --record -f deployment-2.yaml
deployment "labkube" configured
```

In the screen with the pods events you should see that a new pod is created.

```shell
NAME                       READY     STATUS    RESTARTS   AGE
...
labkube-76c7c9754d-6tsf5   0/1       Pending   0         0s
labkube-76c7c9754d-6tsf5   0/1       Pending   0         0s
labkube-76c7c9754d-6tsf5   0/1       ContainerCreating   0         0s
labkube-76c7c9754d-6tsf5   1/1       Running   0         2s
```

In the screen with the replicaSets events you should see that the replication control has been update, the count are modified.

```shell
NAME                 DESIRED   CURRENT   READY     AGE
labkube-76c7c9754d   1         1         1         22s
labkube-76c7c9754d   2         1         1         1m
labkube-76c7c9754d   2         1         1         1m
labkube-76c7c9754d   2         2         1         1m
labkube-76c7c9754d   2         2         2         1m
``` 

Now let's modify the definition of the pod. This will trigger the creation of a new replicaSet and a pod rolling update:

```shell 
> kubectl apply --record -f deployment-3.yaml
deployment "labkube" configured
```

We can clearly see the rolling update sequence on the replicaSet counters:
``` shell
NAME                 DESIRED   CURRENT   READY     AGE
...
labkube-76c7c9754d   2         2         2         1m       # Initial state
labkube-7c7fbbc857   1         0         0         0s       # New Replication controller created
labkube-7c7fbbc857   1         0         0         0s       
labkube-76c7c9754d   1         2         2         9m       
labkube-7c7fbbc857   2         0         0         0s       
labkube-76c7c9754d   1         2         2         9m
labkube-7c7fbbc857   2         1         0         0s       # +1 new pods
labkube-76c7c9754d   1         1         1         9m       # -1 old pods  (but the new one is not yet ready!)
labkube-7c7fbbc857   2         1         0         0s
labkube-7c7fbbc857   2         2         0         0s       # +1 new pods
labkube-7c7fbbc857   2         2         1         0s       
labkube-76c7c9754d   0         1         1         9m      
labkube-76c7c9754d   0         1         1         9m
labkube-76c7c9754d   0         0         0         9m       # -1 old pods
labkube-7c7fbbc857   2         2         2         0s
```

You can see that during the rolling update the Ready count goes down to 1.

- First we have not defined what qualifies our pod as "Ready": this is the purpose of a `readiness` probe that we can define at container level. Check documentation [here](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/). Here is a probe for our container (check in file deployment-4.yaml):

``` yaml
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 3
          periodSeconds: 3
```

- Second we have not tuned our rolling update strategy. We can set some parameters to be sure that we will have always 2 pods running at any point in time. Check the documentation [here](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy). Here is a strategy to have alwasy 2 pods (check in file deployment-4.yaml):

```yaml
        strategy:
            rollingUpdate:
            maxSurge: 1
            maxUnavailable: 0
            type: RollingUpdate
```

Let's perform an update with the new deployment to check the counters:

```shell
> kubectl apply --record -f deployment-4.yaml
deployment "labkube" configured
```

Now the sequence of counter always show at least 2 ready pod at any time:

```shell
NAME                 DESIRED   CURRENT   READY     AGE
labkube-7c7fbbc857   2         2         2         23m          # initial state
labkube-86c45b695    1         0         0         0s           # rreation of new replicaSet
labkube-86c45b695    1         0         0         0s
labkube-86c45b695    1         1         0         0s
labkube-86c45b695    1         1         1         6s           # +1 new pod ready
labkube-7c7fbbc857   1         2         2         24m
labkube-86c45b695    2         1         1         6s
labkube-7c7fbbc857   1         2         2         24m          
labkube-86c45b695    2         1         1         6s
labkube-7c7fbbc857   1         1         1         24m          # -1 old pod ready
labkube-86c45b695    2         2         1         6s
labkube-86c45b695    2         2         2         9s           # +1 new pod ready
labkube-7c7fbbc857   0         1         1         24m
labkube-7c7fbbc857   0         1         1         24m
labkube-7c7fbbc857   0         0         0         24m          # -1 old pod 

```

A deployment can `pause/resume` a deployment. It is also possible to `rollback` a deployment:

```shell
> kubectl rollout undo deployment/labkube
deployment "labkube"
```

Pay attention 2 consecutive ```rollout undo``` takes you back to the initial state before the first ```rollout undo```. You can undo to a dedicated version using the flag ```--to-revision```. The revision number can be found in the history.The history is available if the deployment was create and updated with the `--record` flag.

```shell
> kubectl rollout history deploy/labkube
deployments "labkube"
REVISION  CHANGE-CAUSE
1         kubectl apply --record=true --filename=deployment-2.yaml
2         kubectl apply --record=true --filename=deployment-3.yaml
3         kubectl apply --record=true --filename=deployment-4.yaml
```