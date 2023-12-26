# Task 17 Створення коду Terraform для Flux на kind_cluster

1. В `main.tf` змінимо модуль що відповідає за розгортання кластеру згідно з завданням. Але оберемо гілку модуля, що може працювати без створення файлу `kubeconfig` та використовую інший вид авторизації

```hcl
module "kind_cluster" {
  source = "github.com/den-vasyliev/tf-kind-cluster?ref=cert_auth"
}
```
2. Виконаємо ініціалізацію terraform:
```sh
✗ terraform init
Terraform has been successfully initialized!
```
- В процесі отримуємо помилку:
```sh
Error: Failed to query available provider packages
│ 
│ Could not retrieve the list of available versions for provider fluxcd/flux: locked provider registry.terraform.io/fluxcd/flux 1.2.1 does not match configured version constraint 1.0.0-rc.3; must use terraform init -upgrade to
│ allow selection of new versions
```
- Щоб уникнути ціє помилки робимо модуль локальним та змінюємо у файлі `terraform.tf` версію провайдера на `1.2.1`
```hcl
  required_providers {
    flux = {
      source  = "fluxcd/flux"
      version = "1.2.1"
    }
```
- Виконуємо оновлення провайдера командою:
```sh
terraform init -upgrade
```
- Перевіримо які модулі були створені щоб видалити повністю файл стану:
```sh
✗ terraform state list
module.flux_bootstrap.flux_bootstrap_git.this
module.kind_cluster.kind_cluster.this
✗ terraform state rm module.flux_bootstrap.flux_bootstrap_git.this
Removed module.flux_bootstrap.flux_bootstrap_git.this
Successfully removed 1 resource instance(s).
✗ terraform state rm module.kind_cluster.kind_cluster.this        
Removed module.kind_cluster.kind_cluster.this
Successfully removed 1 resource instance(s).
✗ kind get clusters            
kind-cluster
✗ kind delete clusters kind-cluster
Deleted clusters: ["kind-cluster"]
```


3. Перевіримо код 
```sh
✗ tf validate
Success! The configuration is valid.
```

4. Виконаємо початкові команду `terraform apply`.
```sh
✗ tf apply
Apply complete! Resources: 5 added, 0 changed, 0 destroyed.
```

5. Створені ресурси:
```sh
$ tf state list
module.flux_bootstrap.flux_bootstrap_git.this
module.github_repository.github_repository.this
module.github_repository.github_repository_deploy_key.this
module.kind_cluster.kind_cluster.this
module.tls_private_key.tls_private_key.this
```

6. Розміщення файлу в [bucket](https://console.cloud.google.com/storage/browser)  
Щоб розмістити файл state в бакеті, ви можете використовувати команду terraform init з опцією --backend-config. Наприклад, щоб розмістити файл state в бакеті Google Cloud Storage, ви можете виконати наступну команду:
```sh
# Створимо bucket:
$ gsutil mb gs://vit-secret
Creating gs://vit-secret/...

# Перевірити вміст диску:
$ gsutil ls gs://vit-secret
gs://vit-secret/terraform/
```
7. Як створити bucket [читаємо документацію](https://developer.hashicorp.com/terraform/language/settings/backends/gcs#example-configuration) та додаємо до основного файлу конфігурації наступний код:

```hcl
terraform {
  backend "gcs" {
    bucket  = "tf-state-prod"
    prefix  = "vit-secret"
  }
}
```
```sh
$ terraform init
$ tf show | more
```

8. Перевіримо список ns по стан поду системи flux:
```sh
✗ k get ns
NAME                 STATUS   AGE
default              Active   16m
flux-system          Active   15m
kube-node-lease      Active   16m
kube-public          Active   16m
kube-system          Active   16m
local-path-storage   Active   16m

✗ k get po -n flux-system
NAME                                       READY   STATUS    RESTARTS   AGE
helm-controller-69dbf9f968-qsgq9           1/1     Running   0          16m
kustomize-controller-796b4fbf5d-jxqdx      1/1     Running   0          16m
notification-controller-78f97c759b-c8vpr   1/1     Running   0          16m
source-controller-7bc7c48d8d-c8kxk         1/1     Running   0          16m
``` 
9. Для зручності встановимо [CLI клієнт Flux](https://fluxcd.io/flux/installation/)
```sh
✗ curl -s https://fluxcd.io/install.sh | bash
✗ flux get all
```

10. Додамо в репозиторій каталог `demo` та файл `ns.yaml` що містить маніфест довільного `namespace`  
```sh
$ k ai "маніфест ns demo"
✨ Attempting to apply the following manifest:

apiVersion: v1
kind: Namespace
metadata:
  name: demo
```
- Після зміни стану репозиторію контролер Flux їх виявить:
    - зробить git clone  
    - підготує артефакт   
    - виконає узгодження поточного стану IC   

У даному випадку буде створено `ns demo`:
```sh

✗ flux logs -f
2023-12-19T08:36:29.686Z info GitRepository/flux-system.flux-system - stored artifact for commit 'Create ns.yaml' 
2023-12-19T08:37:31.484Z info GitRepository/flux-system.flux-system - garbage collected 1 artifacts 

✗ k get ns 
NAME                 STATUS   AGE
default              Active   23m
demo                 Active   4s
```
Це був приклад як Flux може керувати конфігурацією ІС Kubernetes

11. Застосуємо CLI Flux для генерації маніфестів необхідних ресурсів:
```sh
$ git clone https://github.com/vit-um/flux-gitops.git
$ cd ../flux-gitops 
$ flux create source git kbot \
    --url=https://github.com/vit-um/kbot \
    --branch=main \
    --namespace=demo \
    --export > clusters/demo/kbot-gr.yaml
$ flux create helmrelease kbot \
    --namespace=demo \
    --source=GitRepository/kbot \
    --chart="./helm" \
    --interval=1m \
    --export > clusters/demo/kbot-hr.yaml
$ git add .
$ git commit -m"add manifest"
$ git push

$ flux logs -f
2023-12-19T08:58:45.061Z info GitRepository/flux-system.flux-system - stored artifact for commit 'add manifest' 
2023-12-19T08:58:45.466Z info Kustomization/flux-system.flux-system - server-side apply for cluster definitions completed 
2023-12-19T08:58:45.559Z info Kustomization/flux-system.flux-system - server-side apply completed 
2023-12-19T08:58:45.596Z info Kustomization/flux-system.flux-system - Reconciliation finished in 498.659581ms, next run in 10m0s 
2023-12-19T08:59:46.501Z info GitRepository/flux-system.flux-system - garbage collected 1 artifacts 
```

11. Перевіримо наявність пода з нашим PET-проектом та розберемо кластер:
```sh
$ k get po -n demo
NAME                         READY   STATUS             RESTARTS       AGE
kbot-helm-6796599d7c-sqwx7   0/1     CrashLoopBackOff   7 (100s ago)   12m
k describe po -n demo | grep Warning
  Warning  BackOff    4m35s (x47 over 14m)  kubelet            Back-off restarting failed container kbot in pod kbot-helm-6796599d7c-sqwx7_demo(401ca7a7-2b0c-4a27-b81c-e053936cd9ed)

$ tf destroy
$ tf state list 
```
