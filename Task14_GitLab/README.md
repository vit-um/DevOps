# Task14 Міграція Pipeline в gitlab cicd

## Підготовка середовища розробки
1. Виявляється що [GitLab](https://docs.gitlab.com/ee/install/) може бути розгорнуто локально з усіма його можливостями та інструментами. 
- З усіх можливих варіантів обираємо вже звичний нам [GitLab Helm chart](https://docs.gitlab.com/charts/) 
- Підготуємо кластер для розгортання пакету хелм
```sh
$ k3d cluster create gitlab-cluster --servers 1 --agents 3 --kubeconfig-update-default
$ kubectl cluster-info
Kubernetes control plane is running at https://0.0.0.0:41913
CoreDNS is running at https://0.0.0.0:41913/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
Metrics-server is running at https://0.0.0.0:41913/api/v1/namespaces/kube-system/services/https:metrics-server:https/proxykubectl cluster-info

$ k config set-context --current --namespace=default

$ helm repo add gitlab https://charts.gitlab.io/
"gitlab" has been added to your repositories
```
- Install GitLab
```sh
$ helm install gitlab gitlab/gitlab \
  --set global.hosts.domain=smart-home.dp.ua \
  --set certmanager-issuer.email=uanetsvilatiy@gmail.com

NAME: gitlab
LAST DEPLOYED: Wed Dec 13 00:56:51 2023
NAMESPACE: default
STATUS: deployed
REVISION: 1
NOTES:
=== CRITICAL
The following charts are included for evaluation purposes only. They will not be supported by GitLab Support
for production workloads. Use Cloud Native Hybrid deployments for production. For more information visit
https://docs.gitlab.com/charts/installation/index.html#use-the-reference-architectures.
- PostgreSQL
- Redis
- Gitaly
- MinIO

=== NOTICE
The minimum required version of PostgreSQL is now 13. See https://gitlab.com/gitlab-org/charts/gitlab/-/blob/master/doc/installation/upgrade.md for more details.

=== NOTICE
You've installed GitLab Runner without the ability to use 'docker in docker'.
The GitLab Runner chart (gitlab/gitlab-runner) is deployed without the `privileged` flag by default for security purposes. This can be changed by setting `gitlab-runner.runners.privileged` to `true`. Before doing so, please read the GitLab Runner chart's documentation on why we
chose not to enable this by default. See https://docs.gitlab.com/runner/install/kubernetes.html#running-docker-in-docker-containers-with-gitlab-runners
Help us improve the installation experience, let us know how we did with a 1 minute survey:https://gitlab.fra1.qualtrics.com/jfe/form/SV_6kVqZANThUQ1bZb?installation=helm&release=16-6
```
- Локальній машині не вистачило ресурсів для нормальної роботи серверу. 

2. Доопрацювати:
- Розгорнути робоче середовище на сервері: 
https://www.youtube.com/watch?v=8r5tF9TZ3wU&t=150s
- Запустити пейплайн  
https://www.youtube.com/watch?v=jAIhhULc7YA  
https://www.youtube.com/watch?v=phlsVGysQSw  

