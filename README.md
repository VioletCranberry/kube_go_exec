# kube_go_exec

The application is written in [GO](https://golang.org/) with purpose of having an easy way to run custom commands inside pods of k8s cluster and being notified about results via Slack. Supported `golang` version is currently `1.16`.

Notifications are handled via [Slack API](https://github.com/slack-go/slack) library and kubernetes [client-go](https://github.com/kubernetes/client-go) is used for talking with [kubernetes](http://kubernetes.io/) cluster.

## Dependencies

The application can be configured via environmenal variables or by supplying `app.env` file into the root directory:

- `SLACK_TOKEN` - slack [access](https://api.slack.com/authentication/token-types) token.
- `SLACK_CHANNEL_ID` - slack channel ID.
- opt. `KUBERNETES_MASTER_URL` - k8s api url
- opt. `KUBERNETES_BEARER_CONFIG` - k8s secret token (see `rbac.yaml`)
- opt. `KUBERNETES_NAMESPACE` - k8s namespace where pod is running.
- `KUBERNETES_POD_LABEL` - k8s pod label to filter pods by.
- opt. `KUBERNETES_CONTAINER` - k8s pod container to exec commands into.
- `KUBERNETES_POD_EXEC` - command to be send to the pod.

The application will try to define whether it is running inside or ouside of the cluster and load / authenticate with appropriate `kube` config.

`service account`, `role` and `rolebinding` are usually needed for the application to work (see `KUBERNETES_BEARER_CONFIG` env. variable) - `rbac.yaml` already includes everything, just apply it on the cluster you are planning to work on and get service account token:

```
kubectl apply -f rbac.yaml
kubectl describe serviceaccount kube-go-exec -n kube-system | grep Tokens
kubectl describe secret <token> -n kube-system
```

## Building and running

```
GOOS=linux GOARCH=amd64 go build
```

This will cross-compile the application and create a binary `kube_go_exec` that you can run:

```
./kube_go_exec
```

## HELM

`kube-go-exec` contains [HELM](https://helm.sh) chart available for deployment:

```
helm install kube-go-exec kube-go-exec
```

In such case `RBAC` resources mentioned above will be created automatically based on `values.yaml` settings and jobs will be executed as `cronjob`. In this case strict `serviceAccount.namespace` <-> `KUBERNETES_NAMESPACE` variable and role-scope mapping applies. You can also supply your custom `RBAC` settings. Variables passed to `secrets` create appropriate resources and secrets will be automatically attached to every `cronjob`.

Alternatively you can install the chart with [helmfile](https://github.com/roboll/helmfile):
```
helmfile -f example_helmfile.yaml apply
```
