#start AD
# closedloop
In this project there is presented an usage of Kubernetes operators to create closed loops. Here there are presented two closed loops: 
responsive - fast and deliberative - slow. Fast manages same of resources as cpu and memory and slow manages the parameters of a fast closed loop as eg. a prority in
serving resources, when simultaneous serving is needed.   

## Description
One closed loop which we treat as operator with closedloop name consists of three operators: monitoring, decision and execution. Separation of this functions is purely conventional. This makes easy to define a closedloop. The deliberative closedloop has added suffix "_d" to his name in order to make different names compare to the responsive closedloop. Then we have for responsive closedloops the following operators: closedloop_d, monitoring_d, decision_d and execution_d.

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

## Prerequisites
We use kubebuilder (https://book.kubebuilder.io/quick-start) to build and run our operators. This tool require the following (state on November 2023):

1. [go version v1.20.0+ ](https://go.dev/dl/)
2. [docker version 17.03+.](https://docs.docker.com/engine/install/) To run docker without sudo, follow [post-installation steps for Docker](https://docs.docker.com/engine/install/linux-postinstall/)
3. [kubectl version v1.11.3+.](https://kubernetes.io/docs/tasks/tools/) To run kubectl without sudo, change the ownership and permissions of ~/.kube
```sh
sudo mv /root/.kube $HOME/.kube # this will write over any previous configuration
sudo chown -R $USER $HOME/.kube
sudo chgrp -R $USER $HOME/.kube
```
4. Access to a Kubernetes v1.11.3+ cluster. For kubernetes cluster you can use [KIND](https://sigs.k8s.io/kind) described above

## Install kubebuilder
[Run the installation procedure.](https://book.kubebuilder.io/quick-start#installation) Run this as sudo user. Check instalation by
```sh
kubebuilder version
```

## Get the project
Use git to clone our project to a machine with above configuration. 

```sh
git clone [ link to git repository ]
```

Let assume that the project is cloned into closedloop directory. So we can compile and run
using bellow commands. This run controller locally. In order to run on k8s, build and deploy a docker container in k8s following instructions described bellow in  [Running on the cluster](#running-on-the-cluster).

```sh
cd closedloop
#Generate your CRD and a configuration file
make generate && make manifests && make install
#run controller localy
make run
#install CR instances of responsive closedloop   
kubectl apply -f config/samples/closedlooppooc_v1_closedloop3.yaml
#install CR instances of deliberative closedloop   
kubectl apply -f config/samples/closedlooppooc_d_v1_closedloop3.yaml
```
Also you need a RESTPod-Listen as proxy from managed system to closedloop, and our managed system RestSys which simple create data and send them (data-send) toRESTPod-Listen. Detailed descriprion in [C) Deploy the operator image and the two Managed Systems](#c-deploy-the-operator-image-and-the-two-managed-systems)

If you want to create you own closedloop project go to [Kubebuilder ClosedLoop on Minikube](#kubebuilder-closedloop-on-minikube)

#stop AD

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



---------------------------

## Kubebuilder ClosedLoop on Minikube

This Part explains how to reproduce this project starting from kubebuilder init (starting from scratch: initialisation of a project in kubebuilder, writing all code on your own, and running the code locally in kubebuilder, i.e., not in kubernetes cluster). Running the loop in real cluster is described in a separate section later on in this document.

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

# CONSOLIDTED ACTION SET
## This workflow is closer to what could be used in practice

### Note: The following description contains also details of manual configuration not mentioned before. They are needed to tune the data sender (data.go application) as we do not use DNS service for local name resolution.

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
#AD
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

   To be done each time for a newly created data-sender instance !!!
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

#AD
##### or with specific data declared manually in the run command
##### Note: single run of data.go results in sending a single "mesurement" message. Emulating a stream of mesurement messages requires multiple runs of data.go.

```sh
go run projects/data.go [cpu] [memory] [disk]
```

### 3. check CR changes

Use kubectl or k9s tool. Display available CRs

```sh
kubectl get crd

```

Next display CR which you want to see, for example:

```sh
kubectl -o yaml get monitoringv2s.closedlooppooc.closedloop.io
```


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

**********************************************************
# DEMONSTRATOR - TWO CASES: REACTIVE CLOSED LOOP and COMBINED REACTIVE & DELIBERATIVE LOOPS
**********************************************************
### Note: In this part of the guide, the steps illustrating the operation of the reactive and deliberative loops are described. The demo uses loop instances created as described in the previous part of this guide. Actually, two separate runs are presented: PART 1: for the reactive loop running in isolation, and PART 2: for the reactive and deliberative loops running in parallel and interacting with each other so that the deliberative component monitors the reactive loop and tunes its parameters according to deliberative loop policy. We remind that acording the the model of our loops described in the report (see a separate Orange-internal document) the parameters of the deliberative loop policy for a particular instance of that loop are specified in the corresponding CRD while the logic of the policy of the loop is hardcoded in source code of the deliberative loop operator.

The workflow of operations within the demo is presented in the figure below. The steps of the workflow are marked with consecutive integers, each step comprising one or two "operations" (symbolically represented as "messages" exchanged between particular functional blocks). The figure given below covers a complete demo workflow, i.e., two loops running in parallel. However, the actions relevant to PART 1 (isolated reactive loop) and PART 2 (both loops are interworking) are easily distinguishable in the figure and we will refer to respective steps in the descriptions that follows.

<img src="./images/1-closed-loop-demo-oct2023.svg" width="90%"></img>

**********************************************************
# REACTIVE CLOSED LOOP (isolated)
**********************************************************
#### This run illustrates the basic worflow within a loop. In particular, one can observe how particular modules engaged in the loop exchange information by modifying dedicated parts of appropriate custom resource (CR). Each change of the CR invokes respective reconciler that executes appropriate logic of a given function of the loop.
**********************************************************

### STEP 1 ###
The Managed System is modelled by the data-sender being a containerized application that sends a report to the Proxy-Pod. This report contains a triplet describing the usage of resources. As shown in the figure below, the triplet is now equal to {CPU:2 Memory:40 Disk:196}. (Note: running data.go application is described in section (H) of the CONSOLIDATED ACTION SET section above.) The containerized app data.go (models the Managed System) always sends a triplet, but in this demo only the tuple CPU and Mem is used; Disk is ignored as early as in the Proxy-Pod). Physically, this report is sent as a REST message to the Proxy-Pod, and in our case Proxy-Pod has forwarded a tuple {CPU=2 Mem=40} to the Monitoring function. All those operations are covered by **step 1** in the figure.

<img src="./images/data_send_1.png" width="50%"></img>

The latter step is confirmed by checking the values (CPU:2, Memory:40) in _Spec.Data_ of Monitoring CR (Custom Resource). This is shown in the screenshot presented below. Additionally, Monitoring operator complements the data tuple with a tag containing the time of receiving the message from Proxy-Pod. This can be seen in the screenshot as parameter _Spec.Time_ with the value _2023-11-27 11:1:32.386296_ which in subsequent steps will be passed on and recorded in all CRs taking part in this instance on the loop. This time tag serves as a unique identifier of message instance so that it is always possible to distinguish between different message instances even if the data carried in those messages (e.g., CPU/Mem/Disk usage) is the same. 

#### STEP 2 ####
Now consider the values received in the context of the monitoring policy applicable in our closed loop. In this case Cpu=2 which is lower than the CPU threshold set to 5 (see _MonitorinData-1-thresholdvalue: 5_ in the figuge below), and Memory=40 is compared to the memory threshold 50 (_MonitorinData-2-thresholdvalue: 50_). In our convention, _MonitorinData-x-thresholkind: inferior_ in the monitoring policy means that if, e.g., _cpu<MonitorinData-1-thresholdvalue_ then the state of CPU is considered to be "Low"; similar interpretation applies to Memory.

<img src="./images/monitoringv2_1.png" width="50%"></img>

According to the interpretation of the thresholds and the comparison conventions of the monitoring policy described above, obtaining the values of resource usage Cpu=2 and Memory=40 results in sending a notification form the Monitoring to the Decision function indicating "Low cpu, Low memory". This can be confirmed by inspecting the value of the field _Spec.Message_ in the Decision CR - see the figure below.

<img src="./images/decision_1.png" width="50%"></img>

#### STEP 3 ####
The value _PriorityRank.rank-1: cpu_ (flull name _Spec.Decisionpolicies.Priorityspec.Priorityrank.rank-1: cpu_) in Decision CR shown above indicates that Cpu has higher rank than Memory. Therefore, as both resource types have been reported as being in shotage ("Low" state for both of them) and only one can be scaled in a given iteration of the loop, it is CPU that is going to be scaled this time (see our Report for more detailed description of the scaling algorithm). Accordingly, "React to cpu" message is sent to Execution which can be verified by inspecting the value of parameter _Spec.Action_ in the Execution CR - see the figure below.

<img src="./images/execution_1.png" width="50%"></img>

#### STEP 4 ####
Step 4 belongs to the deliberative loop and will be referred to in PART 2 below.

#### STEP 5 and STEP 6 ####
Step 5 and step 6 cover subsequent reaction of Execution to receiving "React to cpu" request from Decision. In the demo, step 5 in Execution does not lead to sending explicit requestr/control message to the Managed System. As we mentioned in the Report, the demo focuses on showing the basic mechanisms feedback operations in the reactive loop have not been implemented (that's also why the data sender application data.go has to be triggered manually to emulate sending consecutive monitoring reports).

**********************************************************
# COMBINED REACTIVE & DELIBERATIVE LOOPS (loop interworking)
**********************************************************
#### This run illustrates the worflow in a setup where a deliberative loop monitors the operation of a reactive loop and in case when the reactive loop does not behave according to the expectation of (policy imposed by) the deliberative component that the latter . In particular, one can observe how particular modules engaged in the loop exchange information by modifying dedicated parts of appropriate custom resource (CR). The steps of each component (reactive loop, deliberative loop) are interleaved in such an order so that the presentation reflects the real flow of operations in the most natural way.
**********************************************************

##### _Naming convention: The components of the deliberative loop in our example cooperate based on similar pronciples as in the reactive loop described above. Therefore, the whole structure of the loop is replicated after the reactive loop. In the implementation, all components of this loop are named adding suffix "D" at the end on respective name._

## Reactive closed loop

In the following figure, we informatively show the master CR of the reactive loop. Important in the context of current experiment are fields _Spec.Decisionpolicies.Priorityspec.Priorityrank.rank-1: cpu_ and _Spec.Decisionpolicies.Priorityspec.Priorityrank.rank-2: memory_. They will be subject to changes based on the decision of the deliberative loop. The latter is a new component making the whole setup more "autonomous" (by observing and tunning the operation of the reactive loop). We will observe the changes of those fields in the course of loop operation.

<img src="./images/monitoringv2_3.png" width="50%"></img>

Starting the experiment: similary to the previous case of isolated reactive loop, the tuple CPU:5 Memory:28 is sent to Monitoring. This is shown in the figure below.

<img src="./images/data_send_3.png" width="50%"></img>

Because this time only emory is below respective threshold (memory shortage), the message "Low memory" is sent to Decision, which is confirmed by inspecting the field _Spec.Message_ in the Decision CR whose value is set to "Low memory" - see the figure below (we skipped Monitoring CR to shorten the presentation).

<img src="./images/decision_3.png" width="50%"></img>

 In reaction to receiving "Low memory", Decision sends "React to memory" notification to Execution. This can be confirmed by inspecting the valu of parameter _|Spec.ction: React to memory_ in Executiuon CR - see the figure below. Also, and according to the detailed description of the deliberative loop operation form section 5 of the Report, in parallel to triggering Execution, the indication of the scaled resource Metric=memory is sent by Decision to MonitoringD that runs in the deliberative loop. This will be shown in the next figure.

<img src="./images/execution_3.png" width="50%"></img>

## Reactive closed loop + Deliberative closedloop:
## Deliberative closed loop

In the MonitoringD CR (see the figure below), we can now check the value _Spec.Data.Metric: memory_ - this is what MonitoringD has just received from the reactive loop component Decision. Also, in the field _Spec.Time_, the value _2023-12-01 21:18:35.936615_ received from the reactive loop has been stored (again, it identifies a message, but also certain "threat" in the loop operation). This value (the tag) will next be sent to DecisionD where it will be also saved in a list _Spec.Data.Memory_ containing memory scaling times. - We will see that in the next figure that follows.

<img src="./images/monitoringd_3.png" width="50%"></img>

<img src="./images/decisiond_3.png" width="50%"></img>

# Reactive closed loop + Deliberative closedloop (before modifying Reactive closed loop):
## Reactive closed loop

Here we see a normal work of a deliberate closedloop. Because execution realizes "React to cpu", metric cpu is send to 
a deliberate closedloop to MonitoringD CR. 

<img src="./images/data_send_4.png" width="50%"></img>

<img src="./images/monitoringv2_4.png" width="50%"></img>

<img src="./images/decision_4.png" width="50%"></img>

<img src="./images/execution_4.png" width="50%"></img>

## Reactive closed loop + Deliberative closedloop (before modifying Reactive closed loop):
## Deliberative closed loop

Metric cpu is written to MonitoringD, but there is not decision yet in DecisionD, so Not metric is set in ExecutionD and 
nothing is send to a master reactive Closedloop

<img src="./images/monitoringd_4.png" width="50%"></img>

<img src="./images/decisiond_4.png" width="50%"></img>

<img src="./images/executiond_4.png" width="50%"></img>

<img src="./images/closedloop_4.png" width="50%"></img>

# Reactive closed loop + Deliberative closedloop (after modifying Reactive closed loop):
## Reactive closed loop

Now a new data is processed by reactive closedloop. React to cpu is triggered en metric cpu is send to a deliberative
MonitoringD

<img src="./images/data_send_5.png" width="50%"></img>

<img src="./images/monitoringv2_5.png" width="50%"></img>

<img src="./images/decision_5.png" width="50%"></img>

<img src="./images/execution_5.png" width="50%"></img>

<img src="./images/closedloop_4.png" width="50%"></img>

## Reactive closed loop + Deliberative closedloop (after modifying Reactive closed loop):
## Deliberative closed loop

But now a new time of metric cpu is sufficient to take a decision by DecisionD and in ExecutionD is Spec.Action=Increase rank and Spec.Metric=memory. This values are send to master reactive ClosedLoop 

<img src="./images/monitoringd_5.png" width="50%"></img>

<img src="./images/decisiond_5.png" width="50%"></img>

<img src="./images/executiond_5.png" width="50%"></img>

## Reactive closed loop + Deliberative closedloop (after modifying Reactive closed loop):
## Reactive closed loop after modification

And we see that Status.Increaserank=memory in ClosedLoop. This value is propagated to Monitoring and Decision for changing priority policy.

<img src="./images/closedloop_6.png" width="50%"></img>

<img src="./images/monitoringv2_6.png" width="50%"></img>

<img src="./images/decision_6.png" width="50%"></img>

# Reactive closed loop + outer closedloop (with new configured Reactive closed loop):
## Reactive closed loop

Here we see a work of both closedloop with a new priority policy set.

<img src="./images/data_send_7.png" width="50%"></img>

<img src="./images/monitoringv2_7.png" width="50%"></img>

<img src="./images/decision_7.png" width="50%"></img>

<img src="./images/execution_7.png" width="50%"></img>

## Reactive closedloop + outer closedloop (with new configured Reactive closed loop):
## Deliberative closedloop

<img src="./images/monitoringd_7.png" width="50%"></img>

<img src="./images/decisiond_7.png" width="50%"></img>

<img src="./images/executiond_7.png" width="50%"></img>




