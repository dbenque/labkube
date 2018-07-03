# Lab Kubernetes presentation

## Configure your environment

- If you don't already have a recent version of `kubectl` (>=1.9.x) you needs to follow the instruction on this page: https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-kubectl
- Configure `kubectl` for accessing the managed kubernetes lab cluster:
  - download the `<namespace>-<user>.kubeconfig.yaml` file with the link provided by the trainer.
  - If you don't already have a `~/.kube/conf` file you can simply move or create a symlink of the file `ls -s <namespace>-<user>.kubeconfig.yaml ~/.kube/conf`. Then you should be able to access the kubernetes API server.
  - Else you can use the provided kubernetes config file without modifing yours by adding the parameter `--kubeconfig` when your are using the commandline: `kubectl --kubeconfig=./<namespace>-<user>.kubeconfig.yaml`.
- If you want to access the kubernetes dashboard, you need to run the `kubectl proxy` command, then go to the URI: `http://localhost:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#!/overview?namespace=<namespace>` (replace `<namespace>` in the URI by your namespace name).

## Architecture

API-Server

Client

namespaces

TODO: Slide
