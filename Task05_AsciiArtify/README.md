# AsciiArtify Kubernetes Deployment Tools Evaluation. Demo.

## Етапи створення проекту з розробки нового програмного продукту 

1. [Concept](./doc/Concept.md) - Документ з концепцією по вибору рішення із задокументованим висновком.
2. [PoC](./doc/POC.md) - Довести технічну чи концептуальну життєздатність ідеї та концепції застосування ArgoCD в якості інструменту CD. 
3. [MVP](./doc/MVP.md) - Створити мінімальний функціональний продукт, який може вивести на ринок та отримати зворотний зв'язок від користувачів. В нашому випадку це демонстрація роботи застосунку AsciiArtify 


## Розгортання застосунку на Kubernetes за допомогою k3d.

1. Встановіть k3d за допомогою команди:
```bash
$ wget -q -O - https://raw.githubusercontent.com/rancher/k3d/main/install.sh | bash
```

2. Створіть новий кластер Kubernetes:
```bash
$ k3d cluster create asciiartify
```

3. Перевірте чи всі сервіси кластера працюють:
```bash
$ kubectl cluster-info
Kubernetes control plane is running at https://0.0.0.0:46155
CoreDNS is running at https://0.0.0.0:46155/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
Metrics-server is running at https://0.0.0.0:46155/api/v1/namespaces/kube-system/services/https:metrics-server:https/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

4. Розгорнемо приклад застосунку "hello-world" із вказаним контейнером (gcr.io/google-samples/hello-app:1.0). Розгортання в Kubernetes визначає, як розгортати та керувати одним чи кількома репліками контейнера в середовищі кластера.

```bash
$ kubectl create deployment hello-world --image=gcr.io/google-samples/hello-app:1.0


```
5. Експозиція Розгортання як Сервісу з доступом зовні:
```bash
$ kubectl expose deployment hello-world --type=LoadBalancer --port=8080

```
Ця команда створює службу (Service) для розгортання "hello-world". Опція --type=LoadBalancer означає, що цей сервіс буде доступний зовні кластера за допомогою зовнішнього балансувальника навантаження. Порт 8080 вказує на порт, на якому служба буде слухати запити.

6.  Отримання Інформації про Сервіс, Ця команда виводить інформацію про службу "hello-world", включаючи зовнішню IP-адресу, за якою можна отримати доступ до розгортання "hello-world".

```bash
$ kubectl get services hello-world
NAME          TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
hello-world   LoadBalancer   10.43.47.129   172.18.0.3    8080:30808/TCP   23m
```

7. Перевірте стан розгорнутого застосунку:
```bash
$ kubectl get pods
NAME                           READY   STATUS    RESTARTS   AGE
hello-world-5bc74c8b8d-xqcd4   1/1     Running   0          24m
```
Ви повинні побачити, що контейнер hello-world є в стані Running.

8. Виведіть вміст логів контейнера:

```bash
$ kubectl logs hello-world-5bc74c8b8d-xqcd4
2023/11/23 21:38:50 Server listening on port 8080
```
Таким чином, ви розгорнули простий "Hello World" застосунок на Kubernetes за допомогою k3d. Будь ласка, зауважте, що для використання цих команд потрібно мати встановлені kubectl та k3d, а також Docker для створення контейнерів.

9. Вбираємо після себе
```bash
# Зупиняємо запущені контейнери 
docker stop $(docker ps -aq)    
fbcf90482d30
259890caa301
9bb71cfa4fad

# Видаляємо контейнери
docker rm $(docker ps -aq)    
fbcf90482d30
259890caa301
9bb71cfa4fad

# Видалити всі іміджі
docker rmi $(docker images -q)

# Ви можете видалити k3d за допомогою наступної команди:
sudo rm /usr/local/bin/k3d
```

