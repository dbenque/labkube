# Configuration

Configuration of an application may be done througt:
- Environment Variables
- Files
- Services

Grabbing your configuration using a Service is something that you can do using the Service Discovery as described in lab4

Environment variables can directly be set inside the pod defintion at container level. We had an example with HOSTNAME in the first lab:

```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: labkube
  name: labkube-env
spec:
  containers:
  - image: dbenque/labkube:v1
    env:                                    # <-- section for environment variables
    - name: MY_LABKUBE_VAR
      value: "Hello from the environment"
    name: labkube
    ports:
    - containerPort: 8080
```

It is possible to inject in environment variables some values that comes from different sources ( see the api [here](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#envvarsource-v1-core) ):
- fieldRef: pod information. You have access to the invariant information such as:
  - metadata.name
  - metadata.namespace
  - spec.nodeName
  - spec.serviceAccountName
  - status.hostIP
  - status.podIP
- resourceFieldRef: information about limits and requests of memory and cpu
- configMapKeyRef: values associated to a key i a configmap
- secretKeyRef: values associated to a key i a configmap

Note that the `metadta.name` value is already injected in the environment variable `HOSTNAME`

## Exercise 1

Modify the pod definition and create a pod containing a variable NAMESPACE.

## Exercise 2

Modify the pod definition and create a pod containing a variable with a label. Does it work? why?


