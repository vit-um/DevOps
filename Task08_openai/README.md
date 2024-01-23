# Встановіть та налаштуйте kubectl-ai плагін для створення ШІ Recommended YAML manifests

## Task steps
1. Create an [API key](https://platform.openai.com/account/api-keys)  
`key name: name: kubectl-ai`  

2. Install and configure [the kubectl-ai plugin](https://github.com/sozercan/kubectl-ai)
```sh
$ wget https://github.com/sozercan/kubectl-ai/releases/download/v0.0.11/kubectl-ai_linux_amd64.tar.gz
$ tar -zxvf kubectl-ai_linux_amd64.tar.gz
$ mv kubectl-ai /usr/local/bin/
$ chmod +x /usr/local/bin/kubectl-ai

$ kubectl plugin list                                                                                
The following compatible plugins are available:
/root/.krew/bin/kubectl-krew
/root/.krew/bin/kubectl-ns
/usr/local/bin/kubectl-ai
/usr/local/bin/kubectl-kubeplugin

$ nano ~/.zshrc
export OPENAI_API_KEY="***************************"

$ source ~/.zshrc

$ export OPENAI_DEPLOYMENT_NAME="gpt-4"
```

3. Practice writing and testing prompts on a local cluster
```yaml
$ k ai "get status of master node" --require-confirmation=false
✨ Attempting to apply the following manifest:

apiVersion: v1
kind: Pod
metadata:
  name: node-status-check
  namespace: default
spec:
  containers:
  - name: check-node-status
    image: bitnami/kubectl:latest
    command: ['/bin/bash', '-c', 'kubectl get node master -o jsonpath="{.status}"'] 
    resources:
      requests:
        cpu: 100m
        memory: 100Mi
  restartPolicy: OnFailure

$ mkdir yaml
$ kubectl ai "create an nginx deployment with 3 replicas" --require-confirmation=false > yaml/app.yaml
```

4. The resulting manifest in the yaml directory in the root of the repository.

| NAME                        | PROMPT                             | DESCRIPTION                                                              | EXAMPLE                                     |
|-----------------------------|------------------------------------|--------------------------------------------------------------------------|---------------------------------------------|
| app.yaml                    | Create Application Config          | YAML to define the basic schema of a Kubernetes application              | [app.yaml](yaml/app.yaml)                 |
| app-livenessProbe.yaml      | Add Liveness Probe                 | YAML to define a liveness probe for your application                    | [app-livenessProbe.yaml](yaml/app-livenessProbe.yaml) |
| app-readinessProbe.yaml     | Add Readiness Probe                | YAML to define a readiness probe for your application                   | [app-readinessProbe.yaml](yaml/app-readinessProbe.yaml) |
| app-volumeMounts.yaml       | Configure Volume Mounts            | YAML to define and configure storage volumes for your application       | [app-volumeMounts.yaml](yaml/app-volumeMounts.yaml) |
| app-cronjob.yaml            | Create Cron Job                    | YAML to define a cron job within your application                       | [app-cronjob.yaml](yaml/app-cronjob.yaml) |
| app-job.yaml                | Create a Job                       | YAML to define a job within your application                            | [app-job.yaml](yaml/app-job.yaml) |
| app-multicontainer.yaml     | Set Up Multi-container Pods        | YAML to define a pod that runs more than one container                  | [app-multicontainer.yaml](yaml/app-multicontainer.yaml) |
| app-resources.yaml          | Configure Resource Usage           | YAML to configure resource requests and limits for your application     | [app-resources.yaml](yaml/app-resources.yaml) |
| app-secret-env.yaml         | Set Up Secrets as Env Variables    | YAML to define environment variables using secrets                      | [app-secret-env.yaml](yaml/app-secret-env.yaml) |


