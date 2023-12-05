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

```sh
$ helm template new-version ./helm -s templates/api-deploy.yaml --set image.tag=build-802e329 > new-version-manifest.yaml

$ kubectl apply -f new-version-manifest.yaml
deployment.apps/new-version-api created

helm template new-version ./helm -s templates/api-deploy.yaml --set image.tag=build-802e329 | kubectl apply -f -

$ k get po                                                                                                      
NAME                                        READY   STATUS              RESTARTS   AGE
new-version-api-57578dfdfb-zj7vc            0/1     ContainerCreating   0          4m30s

$ k logs new-version-api-57578dfdfb-zj7vc 
Error from server (BadRequest): container "api" in pod "new-version-api-57578dfdfb-zj7vc" is waiting to start: ContainerCreating

$ k describe pod new-version-api-7bd8447c4f-zzvxw
Events:
  Type     Reason       Age               From               Message
  ----     ------       ----              ----               -------
  Normal   Scheduled    67s               default-scheduler  Successfully assigned default/new-version-api-7bd8447c4f-zzvxw to k3d-k3d-cluster-429-agent-2
  Warning  FailedMount  3s (x8 over 67s)  kubelet            MountVolume.SetUp failed for volume "data" : configmap "new-version-configmap" not found

$ k get events | grep Warning
117s        Warning   FailedMount         pod/new-version-api-7bd8447c4f-zzvxw             MountVolume.SetUp failed for volume "data" : configmap "new-version-configmap" not found

$ helm ls
NAME            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART           APP VERSION
current-version default         1               2023-12-05 13:36:53.894866152 +0200 EET deployed        helm-0.1.0      1.0    

$ k get configmaps
$ k create configmap new-version-configmap --from-literal=key=demo -n default

$ k describe pod new-version-api-57578dfdfb-8nbg8 | grep Warning
  Warning  Failed     23m (x8 over 24m)     kubelet            Error: secret "new-version-secret" not found

$ k get secrets 
$ helm template new-version ./helm -s templates/secret.yaml 
$ kubectl create secret generic new-version-secret --namespace=default --from-literal=license=MTIzNDU=

$ k label pod new-version-api-57578dfdfb-x522w app=current-version-api --overwrite=true
$ k delete deployment current-version-api
$ k delete po current-version-api-dbc75f7b6-7cwbq


```

If you encounter an "ContainerCreating" message for Kubernetes pod, that is the first issue we need to resolve.

Second request. The new version of the api microservice should be deployed at the same endpoint as the current version but not available publicly. Only qa team should be able to access it.

```sh
$ k delete deployments new-version-api 
deployment.apps "new-version-api" deleted
```

