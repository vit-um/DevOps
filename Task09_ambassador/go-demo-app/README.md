# Let's KIK!

#
## Runc and containerazed app https://developers.redhat.com/blog/2018/02/22/container-terminology-practical-introduction#containerized__components
	mkdir demo

vim Container
    
	FROM busybox
	CMD while true; do { echo -e 'HTTP/1.1 200 OK\n\nVersion: v1.0.0'; }|nc -vlp 8080;done
	EXPOSE 8080
 #   
	img build -t denvasyliev/demo:v1.0.0 -f Container .
	img unpack denvasyliev/demo:v1.0.0
	runc spec
	runc run demo
	vim config.json:
		terminal false
		"sh", "-c", "while true; do { echo -e 'HTTP/1.1 200 OK\n\nVersion: v1.0.0'; }|nc -vlp 8080;done"
		rootfs false
		"path": "/var/run/netns/runc"
## Network
	sudo brctl addbr runc0
	sudo ip link set runc0 up
	sudo ip addr add 192.168.10.1/24 dev runc0
	sudo ip link add name veth-host type veth peer name veth-guest
	sudo ip link set veth-host up
	sudo brctl addif runc0 veth-host
	sudo ip netns add runc
	sudo ip link set veth-guest netns runc
	sudo ip netns exec runc ip link set veth-guest name eth1
	sudo ip netns exec runc ip addr add 192.168.10.101/24 dev eth1
	sudo ip netns exec runc ip link set eth1 up
	sudo ip netns exec runc ip route add default via 192.168.10.1
	
    runc run demo
	curl 192.168.10.101:8080
	img push denvasyliev/demo:v1.0.0
	check size on dockerhub

# K8S DEPLOYMENT
	k create deploy demo --image denvasyliev/demo:v1.0.0

	k expose deploy demo --port 80 --type LoadBalancer --target-port 8080 
	k get svc -o wide -w
	tmux: 
		alias k=kubectl
		k get po -lapp -w

	LB=$(k get svc demo -o jsonpath="{..ingress[0].ip}")
	while true;do curl -s $LB;sleep 0.7;done
	vim Container
	#change to v2.0.0
	img build -t denvasyliev/demo:v2.0.0
	img push denvasyliev/demo:v2.0.0

	k set image deployment demo demo=denvasyliev/demo:v2.0.0 --record
	k rollout history deploy demo
	k rollout undo deploy demo --to-revision 1

# K8S CANARY
	k create deploy demo2 --image denvasyliev/demo:v2.0.0
	k get po --show-labels
	k get svc -o wide

	k scale deploy demo --replicas 10

	k label po -lapp run=demo
	k get svc -o wide
	k patch svc demo --patch '{"spec":{"selector": {"run": "demo"} }}'
	k get svc -o wide

	k patch svc demo --type json   -p='[{"op": "remove", "path": "/spec/selector/app"}]'
	k scale deploy demo2 --replicas 5
	k label po -lapp run=demo --overwrite

	k scale deploy demo --replicas 1
	k scale deploy demo --replicas 0
	k scale deploy demo2 --replicas 1


# API-GATEWAY
	git clone https://github.com/den-vasyliev/go-demo-app.git
	cd go-demo-app
	helm template ./helm --namespace demo --name demo |k apply -f -
	watch kubectl get po,svc -n demo

	k -n demo patch svc ambassador --patch '{"spec":{"type": "LoadBalancer"}}'
	DEMOLB=$(k -n demo get svc ambassador -o jsonpath="{..ingress[0].ip}")
	curl -XPOST --data '{"text":"test"}' $DEMOLB/ascii/

	watch -t -d -n 0.5 curl -s $DEMOLB/api/
	helm template ./helm --namespace demo --name demo2 --set image.tag=v2 --set app.version=v2 --set api.canary=30|k apply -f -

	wget -O /tmp/g.png https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png
	curl -F 'image=@/tmp/g.png' $DEMOLB/img/

	helm template ./helm --namespace demo --name demo2 --set image.tag=v2 --set app.version=v2 --set api.canary=80|k apply -f -
	helm template ./helm --namespace demo --name demo2 --set image.tag=v2 --set app.version=v2 --set api.canary=100|k apply -f -

	helm template ./helm --namespace demo --name demo --set ns.enable=false --set api-gateway.enable=false|k delete -f -

	helm template ./helm --namespace demo --name demo3 --set image.tag=v3 --set app.version=v3 --set api.canary=50 --set api.header=canary|k apply -f -

	watch kubectl get po,svc -n demo

	watch -t -d -n 0.5 curl -s $DEMOLB/api/ -Hx-mode:canary
	watch -t -d -n 0.5 curl -s $DEMOLB/api/ 

	helm template ./helm --namespace demo --name demo3 --set image.tag=v3 --set app.version=v3 --set api.canary=100 |k apply -f -
	helm template ./helm --namespace demo --name demo2 --set ns.enable=false --set api-gateway.enable=false|k delete -f -

	watch -t -d -n 0.5 curl -s $DEMOLB/api/ 
	OPEN_IN_BROWSER $DEMOLB/ml5/

# DEMO APP API
	curl -XPOST --data '{"text":"test"}' $DEMOLB/ascii/
	curl -F 'image=@/tmp/g.png' $DEMOLB/img/
	OPEN_IN_BROWSER $DEMOLB/ml5/

# EFK

	https://github.com/upmc-enterprises/elasticsearch-operator

## HELM packages

	helm repo add es-operator https://raw.githubusercontent.com/upmc-enterprises/elasticsearch-operator/master/charts/
	helm fetch es-operator/elasticsearch-operator	
	helm fetch es-operator/elasticsearch	

## Elasticsearch Operator

	k create ns logging	

	helm template --name elasticsearch-operator elasticsearch-operator-0.1.3.tgz --set rbac.enabled=True --namespace logging | k create -n logging -f -	

	k -n logging logs -f	
	k get po -n logging -w	

## Elasticsearch Cluster

	helm template --name=elasticsearch elasticsearch-0.1.5.tgz \
	--set clientReplicas=1 \
	--set masterReplicas=1 \
	--set dataReplicas=2 \
	--set dataVolumeSize=10Gi \
	--set kibana.enabled=True \
	--set cerebro.enabled=True \
	--set zones="{europe-west4-b}" \
	--set storage.type=pd-standard \
	--set storage.classProvisioner=kubernetes.io/gce-pd \
	--namespace logging|k -n logging apply -f -

## Access Cluster
	k port-forward kibana-.... 5601 -n logging

## Fluentbit
	https://docs.fluentbit.io/manual/installation/kubernetes/

	k create -f https://raw.githubusercontent.com/fluent/fluent-bit-kubernetes-logging/master/fluent-bit-service-account.yaml
	
	k create -f https://raw.githubusercontent.com/fluent/fluent-bit-kubernetes-logging/master/fluent-bit-role.yaml
	
	k create -f https://raw.githubusercontent.com/fluent/fluent-bit-kubernetes-logging/master/fluent-bit-role-binding.yaml

	k create -f https://raw.githubusercontent.com/fluent/fluent-bit-kubernetes-logging/master/output/elasticsearch/fluent-bit-configmap.yaml

	k create -f https://raw.githubusercontent.com/fluent/fluent-bit-kubernetes-logging/master/output/elasticsearch/fluent-bit-ds.yaml

# ISTIO

	git clone https://github.com/banzaicloud/istio-operator.git -b release-1.0
	cd istio-operator
	make deploy
	kubectl create -n istio-system -f config/samples/istio_v1beta1_istio.yaml
	k get po -n istio-system -w
	kubectl label namespace demo istio-injection=enabled
	wget https://github.com/istio/istio/releases/download/1.0.5/istio-1.0.5-osx.tar.gz -O -|tar zxf - istio-1.0.5/bin/istioctl
	mv istio-1.0.5/bin/istioctl /usr/bin/istioctl
	alias i=/usr/bin/istioctl
	i version

	k get -n istio-system MeshPolicy default
	cat <<EOF | kubectl apply -f -
	apiVersion: "networking.istio.io/v1alpha3"
	kind: "DestinationRule"
	metadata:
	  name: "default"
	  namespace: "default"
	spec:
	  host: "*.local"
	  trafficPolicy:
	    tls:
	      mode: ISTIO_MUTUAL
	EOF


	helm template ./helm --namespace demo --name demo \
	--set app.version=v3 \
	--set api-gateway.enable=false \
	|i kube-inject -f - \
	|k apply -n demo -f -

	k get po -n demo -w


# SETUP KNATIVE
	
	apt-get install zsh unzip tmux kubectx -y && \
	sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)" 

	wget https://github.com/cppforlife/knctl/releases/download/v0.1.0/knctl-linux-amd64 -O kn && \
	chmod +x kn&&cp kn /usr/bin/kn&&ln -s /usr/bin/kn /usr/bin/knctl

	curl -L https://git.io/getLatestIstio | sh -
	cp istio-1.0.6/bin/istioctl /usr/local/bin
	rehash

	kn install

### Broken version
	kn deploy --service k8s-art --image docker.io/denvasyliev/k8s-art:v3.0.1
	kn curl --service k8s-art
	http://k8s-art.demo.k8s-diy.xyz

### Worked verison
	kn deploy --service k8s-art --image docker.io/denvasyliev/k8s-art:v3.0.1-k

	kn curl --service k8s-art

	http://k8s-art.demo.k8s-diy.xyz

	kn pod list --service k8s-art
	kn revision list --service k8s-art
	kn logs -f --service k8s-art

	https://github.com/cppforlife/knctl/tree/master/docs/cmd

### local registry
	git clone https://github.com/triggermesh/knative-local-registry.git
	cd knative-local-registry
	k create ns registry
	k apply -f templates

### docker-hub credentials
	export KNCTL_NAMESPACE=demo
	kn basic-auth-secret create -s docker-reg1 --docker-hub --username NAME --password 'PASS'
	kn service-account create -a serv-acct1 -s docker-reg1
### source-to-url
	kn deploy \
    	--service k8s-art \
	    --git-url https://github.com/den-vasyliev/go-demo-app \
	    --git-revision master \
	    --service-account serv-acct1 \
	    --image index.docker.io/denvasyliev/k8s-art-k 


