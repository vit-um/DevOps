#  Задача 10. 
Підготувати `helm template` з необхідними параметрами для розгортання нової версії мікросервісу.

Here are the steps to reproduce the problem:

1. Clone the application repository using the following command:
```sh
$ git clone --depth=1 https://github.com/den-vasyliev/go-demo-app.git
$ cd go-demo-app 
```
2. Install the application using Helm:
```sh
$ helm install current-version ./helm
NAME: current-version
LAST DEPLOYED: Tue Dec  5 12:29:46 2023
NAMESPACE: default
STATUS: deployed
REVISION: 1
```
3. This will deploy a boundle of microservices, databases, message broker and api-gateway with service called "ambassador".
```sh
$ k get po --show-labels  
NAME                                        READY   STATUS    RESTARTS   AGE    LABELS
current-version-data-6ffb8f8645-skm7r       1/1     Running   0          5m5s   app=current-version-data,pod-template-hash=6ffb8f8645,version=v4
current-version-api-dbc75f7b6-m9tpr         1/1     Running   0          5m5s   app=current-version-api,pod-template-hash=dbc75f7b6,version=v4
current-version-img-68849cbb77-jzsqd        1/1     Running   0          5m5s   app=current-version-img,pod-template-hash=68849cbb77,version=v4
current-version-ascii-76779ffbfd-hwnvv      1/1     Running   0          5m5s   app=current-version-ascii,pod-template-hash=76779ffbfd,version=v4
current-version-nats-box-7cc65b6579-p8kfd   1/1     Running   0          5m5s   app=current-version-nats-box,pod-template-hash=7cc65b6579
current-version-front-6d5bbd46bb-7dk7t      1/1     Running   0          5m5s   app=current-version-front,pod-template-hash=6d5bbd46bb,version=3.0.1
cache-858575fc54-lz7lk                      1/1     Running   0          5m5s   app=cache,pod-template-hash=858575fc54
db-7968646c85-672g8                         1/1     Running   0          5m5s   app=db,pod-template-hash=7968646c85
current-version-nats-0                      3/3     Running   0          5m5s   app.kubernetes.io/instance=current-version,app.kubernetes.io/name=nats,controller-revision-hash=current-version-nats-5c95bdb489,statefulset.kubernetes.io/pod-name=current-version-nats-0
ambassador-54bb4484c-59qfq                  1/1     Running   0          5m5s   pod-template-hash=54bb4484c,service=ambassador
```

4. Forward the service port to your local machine:

```sh
$ k get svc ambassador
NAME         TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
ambassador   NodePort   10.43.217.50   <none>        80:31857/TCP   40m

$ k port-forward svc/ambassador 8080:80&
Handling connection for 8080

$ lsof -i :8080
COMMAND   PID USER   FD   TYPE  DEVICE SIZE/OFF NODE NAME
kubectl 40427 root    8u  IPv4 1092380      0t0  TCP localhost:http-alt (LISTEN)
kubectl 40427 root    9u  IPv6 1093197      0t0  TCP ip6-localhost:http-alt (LISTEN)                                   
```

5. Test the application by running the following command:

```sh
$ curl localhost:8080/api/
k8sdiy-api:599e1af#     
```

6. You should get the current version of the api microservice in the response: k8sdiy-api:599e1af

Finally, try to set new version and deploy a microservice using the following command:

helm template new-version ./helm -s templates/api-deploy.yaml --set image.tag=build-802e329

If you encounter an "ContainerCreating" message for Kubernetes pod, that is the first issue we need to resolve.

Second request. The new version of the api microservice should be deployed at the same endpoint as the current version but not available publicly. Only qa team should be able to access it.