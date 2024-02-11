# closedloop
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

```sh
kubectl apply -f config/samples/WhatYouWantTo
```

2. Build image :

```sh
make docker-build IMG=controller:latest 
```

3. Save image as file to then send it to minikube
```sh
docker save -o ./savedimage controller:latest
```
then on minikube : 
```sh
docker load -i savedimage
```


4. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=controller:latest
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy IMG=controller:latest
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.


---------------------------

## Kubebuilder ClosedLoop on Minikube

This Part explain how to reproduce this project starting from kubebuilder init.

I. Init your Project 

```sh
username:~/closedloop$ kubebuilder init --domain closedloop.io --repo closedloop //Init folder
username:~/closedloop$ kubebuilder create api --group closedlooppooc --version v1 --kind ClosedLoop    //create API and Controller
username:~/closedloop$ kubebuilder create api --group closedlooppooc --version v1 --kind Monitoring    //create API and Controller
username:~/closedloop$ kubebuilder create api --group closedlooppooc --version v1 --kind Decision      //create API and Controller
username:~/closedloop$ kubebuilder create api --group closedlooppooc --version v1 --kind Execution     //create API and Controller
username:~/closedloop$ kubebuilder create api --group closedlooppooc --version v1 --kind Monitoringv2  //create API and Controller
```
II. Complet your API SPEC

Go to the folder : api/yourVersion and complete all the _types.go file to describe your CR Spec and Status

III. Generate your CRD and configuration file based on what you did on the _types.go files

```sh
username:~/closedloop$ make generate && make manifests && make install
```

IV. Complete the logic of your controller

Complete code of the controller files in the "/controllers" folder

V. Run your Project to test it localy (This is not like in production, refer to VII)

```sh
username:~/closedloop$make run
```

VI. Create your CR Ressources

Complete/fill in the files on /config/samples as a example
and excecute the command:
```sh
username:~/closedloop$ kubectl apply -f config/samples/closedlooppooc_v1_closedloop.yaml //(Example)
```

VII. Deploy your Operator like in production

Excecute the commands:

```sh
username:~/closedloop$ make docker-build IMG=controller:latest && docker save -o ./savedimage controller:latest
```

For Minikube : 
    Transfert the savedimage file to your minikube VM and build it : example 

Run From Minikube (ssh) to retreive from where your build the image

```ssh
$scp Username@IP:/Path/To/savedimage ./  // Copy the file in local
$docker load -i savedimage               // Load the Image in Minkube
```

Run on the Kubebuilder Host to Deploy your Operator, RBAC file, ..) :

```ssh
username:~/closedloop$ make deploy IMG=controller:latest
```

VIII. Load the Proxy Pod 

Run From Minikube (ssh) 

```ssh
$scp username@IP:/Path/To/closedloop/RESTPod-Listen/* ./ && docker build -t restpod:latest . //This will retreive and build the image needed for the proxy pod
```

VIV. Deploy the 2 Managed Systems

1) Exporter : 

Run From Minikube (ssh) 

```ssh
$scp username@IP:/Path/To/closedloop/exporter/* ./ && docker build -t exporter . //This will retreive and build the image needed for the exporter
```
Run on the Kubebuilder Host

```ssh
username:~/closedloop$ kubectl create -f ./exporter/exporter.yaml //This will create the exporter
```

2) PodToPushData to Proxy-Pod : 

Run From Minikube (ssh) 

```ssh
$scp username@IP:/Path/To/closedloop/RESTSys/* ./ && docker build -t data-send:latest . //This will retreive and build the image needed for the POdToPushData to Proxy-Pod
```
Run on the Kubebuilder Host

```ssh
username:~/closedloop$ kubectl create -f ./RESTSys/data-send-deployment.yaml //This will create the PodToPushData
```

# CONSOLIDTED ACTION SET - workflow typically used in practice (by DB)

### Note: This description also contains details of manual configuration not mentioned before. They are needed to tune data sender as we do not use DNS service for local name resolution.

**********************************************************
# UNINSTALL/UNDEPLOY
**********************************************************

##### UnDeploy the controller from the cluster:

```sh
make undeploy IMG=controller:latest
```

##### Delete the CRDs from the cluster:

```sh
make uninstall
```

**********************************************************
# DO NEEDED STUFF TO RUN THE DEMO
**********************************************************

We assume all code has already been provided

**********************************************************
## A) Generate controller code and artifacts (CDRs)
**********************************************************

#### 1.on kubebuilder node, from closedloop directory ~/.../closedloop/

```sh
make generate && make manifests && make install
make docker-build IMG=controller:latest && docker save -o ./savedimage controller:latest
```

#### 2. ssh to minikube (you sh to the master node): create operator image and load the image; check images

```sh
scp minikube@10.0.2.15://home/minikube/demos/closedloop-ad/closedloop/savedimage ./
docker load -i savedimage
docker image list
```

**********************************************************
## B) Run the master operator (closedloop)
**********************************************************

### on Kubebuilder Host, to Deploy your Operator, RBAC file, ...

```sh
~/.../closedloop/make deploy IMG=controller:latest
```

***********************************************************
## C) Deploy the operator image and the two Managed Systems
***********************************************************

##### Note: PodToPushData and Proxy-Pod together correspond to (represent) one of the two managed systems while exporter represents the second managed system.

### 1. load the code, generate the Proxy-Pod

Note: PodToPushData and Proxy-Pod work together to feed respective instance of a closed loop with monitoring data (by their design and the instantiation process, both of them correspond to one common instance of closed loop). PodToPushData generates random numbers for CPU, RAM and Disk usage and sends them to the Proxy-Pod. Proxy-Pod runs Python Simple HTTP Server that receives (PUT) the requests form PodToPushData Pod and resends them to the closed loop by accessing and modifying the value of parameter Data (and also Time) in the spec section of the Monitoring Custom Resource. This custom resource represents a given instance of the closed loop. Changing the value of Data/Time parameter pair triggers the reconciliation loop of the Monitoring operator thereby propelling the whole closed loop to run.

##### run from kubebuilder host

```sh
scp minikube@10.0.2.15://home/minikube/demos/closedloop-ad/closedloop/RESTPod-Listen/* ./ && 
socker build -t restpod:latest .
```

### 2. load the code, generate and create the exporter

Note: exporter is a Deployment running nginx web server togehter with a Python script that cyclically generates random values for the usage of CPU, RAM and Disk and writes then into the index.html of the server. The server can then be queried (GET) for the contents of the index page. However, currently we do not use exporter in our demos.

##### run from kubebuilder host

```sh
scp minikube@10.0.2.15://home/minikube/demos/closedloop-ad/closedloop/exporter/* ./ && docker build -t exporter .
```

##### run on kubebuilder host

```sh
kubectl create -f ./exporter/exporter.yaml
```

### 3. prepare the image for the data-sender Pod (i.e., PodToPushData that sends data to the Proxy-Pod) and create the data-sender Pod (PodToPushData)

##### run from ssh/minikube

```sh
scp minikube@10.0.2.15://home/minikube/demos/closedloop-ad/closedloop/RESTSys/* ./ && docker build -t data-send:latest .
```

### 4. create CRB (Cluster Role Binding) to allow the ProxyPod to write to the Monitoring CR

Below, we create a CRB (Cluster Role Binding) to allow ProxyPod accessing (i.e., editing) the Monitoring CR (with somewhat confusing name of the CR being closedloop-v2-monitoring-xyz...).

##### run from kubebuilder host

```sh
~/demos/closedloop/RESTPod-Listen$ kubectl apply -f .
```

**********************************************************
## D) Create the closed loop (all resources recursively)
**********************************************************

##### run from kubebuilder host

```sh
kubectl apply -f config/samples/closedlooppooc_v1_closedloop3.yaml
```

##### run for the deliberative (second) loop

```sh
kubectl apply -f config/samples/closedlooppooc_d_v1_closedloop3.yaml
```

**********************************************************
## E) Monitor pod's log (update to give your manager name)
**********************************************************

```sh
kubectl logs -f -n closedloop-system closedloop-controller-manager-7d9bf7cffd-b4g7n
```

**********************************************************
## F) Run data sender deployment (emulates the managed system as the source of events)
**********************************************************

##### run from kubebuilder host

```sh
kubectl create -f ./RESTSys/data-send-deployment.yaml
```

**********************************************************
## G) Update the hosts file in the data-sender deployment Pod

   To be done each time for a newly run data-sender instance !!!
**********************************************************

##### from ssh/minikube

look for POST message and notice the ProxyPod service name (for DNS resolution) in the form: closedloop-v2-monitoring-deployment-service.com:80

```sh
cat data.go
ip a
```

take note of the eth0 IP address above - this is the k8s node address to be used in the NodePort service type for the ProxyPod

##### assume the address is 192.168.49.2

(alternatively to the above, you can simply run "$ minikube ip" on the minikube/kubebuilder host)

##### on kubebuilder, login to the data-sender Pod to set the NodePort IP address for the ProxyPod service (remember to adjust the name of your data-send-deployment Pod)

```sh
kubectl get pods -A ==> check the name of data-sender Pod
kubectl exec --stdin --tty data-send-deployment-6c9f7dd689-qstdr -- /bin/bash
```

##### in the data-sender Pod, insert a DNS entry in the hosts file (adjust the address to your environment)

```sh
nano /etc/hosts
```

  and add a line as follows:
```sh
192.168.49.2 closedloop-v2-monitoring-deployment-service.com
```

  Note: closedloop-v2-monitoring-deployment-service.com is the FQDN of the ProxyPod service as hardcoded in the data-sender Pod program. If one sets a local DNS server able to resolve that FQDN onto the minikube node IP address than the above change is not needed. Configuring the receiver of the monitoring data is always specific and can be troublesome. Future work could focus on integrating with Prometheus, etc. But for now we are fine with workarounds as the one above.

**********************************************************
## H) Generate events using data-sender to test loop operation

   For visibility reasons, it is recommended to open 3 terminals of k9s and in each of them observe (Ctrl-D) the spec section of the custom resource Monitoring2, Decision and Excecution, respectively. One then will be able to easily trace the change of the spec properties involved in the message flow. Leverage on the use of Kubernetes ecosystem tools!
**********************************************************

### 1. remember to have run from kubebuilder host

```sh
~/demos/closedloop/$ kubectl apply -f ./RESTPod-Listen/
```

### 2. repeat multiple times

##### on the kubebuilder host - open shell on the data-sender Pod (adjust the data-sender name!!!)

```sh
kubectl exec --stdin --tty data-send-deployment-6c9f7dd689-qstdr -- /bin/bash
```

##### on the data-sender Pod

```sh
go run projects/data.go
```
##### or with specific data

```sh
go run projects/data.go [cpu] [memory] [disk]
```
