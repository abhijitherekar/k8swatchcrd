# CREDITS: 
1. Kubernetes community and kubernetes github page of client-go, apiextensions
2. Sample controller youtube video by rook, for example on codegen.sh

# k8swatchcrd

K8swatchcrd is an application, which will create a controller
and watch for the Pod resources that have been created in the cluster.

Also, it creates a additional CRD which keeps the count of the following 
and stores them in the etcd :

	1.Total number of pods created from when the cluster was created.
	2.Current number of pods in the cluster that are running.

The CRD name is monitor, the apigroup name is : monitors.k8spodcnt.io

example to fetch the data:

kubectl get monitors.k8spodcnt.io podcount -o yaml

# STEPS TO RUN THE APP

# 1
**** Apply the configmap before applying the POD yaml ****

kubectl apply -f k8swatchcrd-config.yaml


	#kubectl get configmap
	NAME          DATA      AGE
	k8swatchcrd   1         10s

**** Apply the configmap before applying the POD yaml ****

# 2
**** Apply the yaml to the Cluster with RBAC rules ****
The yaml file creates the following:
1. RBAC rules.
2. ServiceAccount named k8swatchcrd which has access to CRD's and k8s core API.
3. Pod yaml with the ServiceAccount k8swatchcrd so that it can access the
	k8s core API and CRD's.
4. Also, it has a configmap volume mounted which tells which k8s resource has
	to be monitored.
5. The Pod container image is pulled from my docker hub repo.

# kubectl apply -f k8swatchcrd.yaml

	# kubectl get pods
	NAME          READY     STATUS    RESTARTS   AGE
	k8swatchcrd   1/1       Running   0          26s

which will create a POD running our
docker image.
**** Apply the docker image to the Cluster****

# 3
**** Now check the CRD with the updates of the POD in cluster ****
# kubectl get monitor podcount -o yaml
	apiVersion: k8spodcnt.io/v1
	kind: Monitor
	metadata:
	  clusterName: ""
	  creationTimestamp: 2018-07-01T08:07:39Z
	  generation: 1
	  name: podcount
	  namespace: default
	  resourceVersion: "681"
	  selfLink: /apis/k8spodcnt.io/v1/namespaces/default/monitors/podcount
	  uid: d2b8a92d-7d05-11e8-9534-000c29da5c7d
	spec:
	  name: podcrd
	status:
	  allpodcnt: 10
	  currpodcnt: 10

**** Now check the CRD with the updates of the POD in cluster ****

# 4
Now if you want to access the data through the REST API
**** Checking the CRD through the REST API ****
# Start the kubectl proxy
	and then  curl '127.0.0.1:<PORT>/apis/k8spodcnt.io/v1/monitors'
**** Checking the CRD through the REST API ****

# 5
**** Creating RBAC rules for the users ****
Currently, only the ADMIN has the access to the cluster but, if you want 
to add more users and create RBAC rules for the users, then do the 
following:

	openssl genrsa -out user1.key 2048

	openssl req -new -key user1.key -out user1.csr -subj "/CN=user1/O=crdgroup"

	openssl x509 -req -in user1.csr -CA /etc/kubernetes/pki/ca.crt -CAkey /etc/kubernetes/pki/ca.key -CAcreateserial -out user1.crt

	openssl x509 -in user1.crt -text

	kubectl config set-credentials user --client-certificate=/home/newusers/user1.crt --client-key=/home/newusers/user1.key

	kubectl config set-context user1@kubernetes --cluster=kubernetes --user=user1

	kubectl config use-context user1@kubernetes

	kubectl config use-context kubernetes-admin@kubernetes

	kubectl create clusterrolebinding crdgroup-admin-binding --clusterrole=cluster-admin --group=crdgroup
****

# 6 If the User wants to play around with the code and build his own images
#	below are the commands
**** DOCKER BUILD **** 
	The dokcer build process:
	To make a docker image run the follwoing command:
	make binary-image

	This will create the follwoing image
	:k8swatchcrd# docker images
	REPOSITORY               TAG                 IMAGE ID            CREATED             SIZE
	herekar/k8swatchcrd      1.0                 ec54c069f9af        9 minutes ago       28.3 MB
**** DOCKER BUILD ****

# 7 Unit test
**** UNIT TEST END-to-END TEST****
	RUN: "make test"

	or cd "opt/controller";go test

	The Controller_test.go has the END to test

**** UNIT TEST END-to-END TEST****

# 8 Future work
	1. Make this work in HA mode using leader election algorithm
	2. Make a CLeaner based on the CRD.
	3. Add a SLACK notification.
	4. Created a HTTPS endpoint to listen on some events and run some events.
	
