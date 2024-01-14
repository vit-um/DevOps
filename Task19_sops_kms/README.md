# Task19 Реалізуйте повну схему SOPS-KMS-FLUX для Kubernetes secret маніфесту, що містить TELE_TOKEN для вашого телеграм бота.

## Coding Session. Terraform + Flux + SOPS

Знайомство з концепцією тераформ модулів, розгорнемо набір інструментів що реалізують повний повний автоматичний цикл на базі GitOps та Kubernetes.
- Terraform створить  Kubernetes cluster та розгорне Flux 
- Flux почне узгоджувати стан ІС та застосунків базуючись на джерелі у GitHub
- GitHub в свою чергу також буде створено за допомогою Terraform

1. Підготуємо файли для автоматичного розгортання ІС:
- Перейменовуємо файл terraform.tfvars.example на terraform.tfvars
- Визначаємо змінні що в ньому завдані

2. Створимо secret для токена
- Призначимо змінній TELE_TOKEN значення яке ми отримали при створенні боту
```sh
$ read -s TELE_TOKEN
# експортуємо його 
$ export TELE_TOKEN
```
- Створимо секрети за допомогою `kubectl create secret` з опцією `dry-run`, що згенерує yaml маніфест. 
```sh
$ k -n demo create secret generic kbot --from-literal=token=$TELE_TOKEN --dry-run=client -o yaml > secret.yaml
```

3. Виконаємо ініціалізацію terraform та розгортання:
```sh
✗ terraform init
Terraform has been successfully initialized!

✗ terraform validate
Success! The configuration is valid.

✗ terraform apply -auto-approve
```

4. Крманди що можуть знадобитись для дігностики проблем:
```sh
✗ kubectl config set-context --current --namespace=demo
✗ kubectl get secrets
✗ kubectl get deploy
✗ kubectl get po
✗ kubectl describe po kbot-66dc58657-tmlpn | grep Warn

✗ gcloud kms keys list --location global --keyring sops-flux  

✗ flux logs -f
✗ flux get all
```

5. Не вийщло здолати помилку №409, тому при кожному новому розгортанні слід використовувати новий `keyring` в наступному модулі: 
```hcl
module "kms" {
  source             = "terraform-google-modules/kms/google"
  version            = "2.2.3"
  project_id         = var.GOOGLE_PROJECT
  keyring            = "sops-flux5"
  location           = "global"
  keys               = ["sops-key-flux"]
  prevent_destroy    = false
}
```
- Не забуваємо декативувати створені [ключі](https://console.cloud.google.com/security/kms/keyrings?project=vit-um) бо вони коштують гроші: 


6. Якщо застосунок не розгорнувся, видаліть сікрети, після чого можна додати новий без шифрування:
```sh
$ k apply -f secret.yaml 
$ k get secrets -n demo
```

7. Видаляємо IC. Помилки tf destroy 
```sh

module.gke_cluster.google_container_node_pool.this: Destruction complete after 4m15s
╷
│ Error: Could not delete Flux
│ 
│ could not clone git repository: unable to clone: repository not found: git repository: 'https://github.com/vit-um/flux-gitops.git'
╵
✗ tf state list
module.flux_bootstrap.flux_bootstrap_git.this
module.gke_cluster.data.google_client_config.current
module.gke_cluster.data.google_container_cluster.main
module.gke_cluster.google_container_cluster.this

✗ tf state rm module.flux_bootstrap.flux_bootstrap_git.this
Removed module.flux_bootstrap.flux_bootstrap_git.this
Successfully removed 1 resource instance(s).
Releasing state lock. This may take a few moments...

✗ tf destroy
Destroy complete! Resources: 1 destroyed.
```


## Оцінювання:

3 бали (junior): secret-маніфест генерується та шифрується мануально

7 балів (middle): secret-маніфест генерується та шифрується за допомогою пайплайну GitHub Actions, токен зберігається у GitHub Secret

10 балів (senior): secret-маніфест генерується та шифрується за допомогою пайплайну GitHub Actions, токен зберігається у GCP Secret Manager