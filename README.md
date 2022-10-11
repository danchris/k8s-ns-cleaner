# Kubernetes Namespace Cleaner
A program written in go that looks for namespaces without active resources and delete them. Default and system namespaces are excluded.

## Usage
```sh
./main --help
  -kubeconfig string
        (optional) absolute path to the kubeconfig file (default "$HOME/.kube/config")
```
