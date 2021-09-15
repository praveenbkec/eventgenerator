# EVENT GENERATOR

### Problem: Implement Event generator Service using microservice architecture and deploy minikube environment using docker and kubernetes tools

1. One service produces kafka events in constant interval say every 30 seconds
2. Events can be employee swiping the card to gain access or train ticket punching. Example: Event = Name:XXXXX,Dept=OSS,EmplD:1234, Time=21-7-2021 21:00:10
3. Other service consumes the kafka events, stores  in database(based on unique fields) if record exists updates it.
4. Expose an endpoint/api to view the events (single and bulk)

   Database: Casandra/ mariadb / any preferred
   Messaging: Kafka

   Language: Go / Java / C / C++ / C#

   Non-Function Requirements:
* Documentation using Swagger
* Performance Metrics
* Docker-compose tests
  The code should be compiled and deployable in minikube

### Solution: This solution is implemented using following tech stack
### PREREQUISITES  
* docker, minikube, bazel, go, kafka, postgres/patroni

### Steps to install/run Infra components
* https://minikube.sigs.k8s.io/docs/start/
* https://helm.sh/docs/intro/install/
* https://docs.bazel.build/versions/main/install-os-x.html (brew install bazel)
* follow below steps to install postgres and kafka
  - helm repo add bitnami https://charts.bitnami.com/bitnami
  - helm install messaging bitnami/kafka --set autoCreateTopicsEnable=true,deleteTopicEnable=true
  - helm install db bitnami/postgresql --set postgresqlUsername=admin,postgresqlPassword=admin,postgresqlDatabase=event_db

    
    
### Steps to run
1. git clone project using
    ```
    git clone https://github.com/praveenbkec/eventgenerator.git
   ```
2. run below bazel command to resolve depenedencies
    ```
    bazel run //:gazelle
    ```
   
3. build event producer docker
    ``` 
    bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 producer:latest
    ```
4. install producer helm chart
    ```
    helm install eventproducer producer/deploy/eventproducer
    ```
3. build event consumer docker
    ```
    bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 consumer:latest
    ```   

6. install consumer helm chart
    ```
    helm install eventconsumer consumer/deploy/eventconsumer
    ```
   
7. check all running pods
```
kmpraveen@kmpraveen-mbp eventgenerator % kubectl get po

NAME                             READY   STATUS    RESTARTS      AGE
db-postgresql-0                  1/1     Running   0             14h
eventconsumer-69bbd4c6df-kb48r   1/1     Running   0             28m
eventproducer-5f8df5c4fb-gmjwl   1/1     Running   0             13h
messaging-kafka-0                1/1     Running   2 (15h ago)   15h
messaging-kafka-client           1/1     Running   0             3d20h
messaging-zookeeper-0            1/1     Running   0             15h
```

8. check logs of a pod
```
kmpraveen@kmpraveen-mbp eventgenerator % kubectl logs -f eventconsumer-69bbd4c6df-kb48r 
===================================== Event received ==========================================
Event : {"Name":"praveen","Dept":"IT","EmpID":"12345","Time":"2021-09-11 18:53:30.0136598 +0000 UTC m=+1.003192301"}
```

9. Generate Swagger/Open API json using commad
```bigquery
bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 consumer/proto:event_swagger
```
10. Run Benchmark tests
```bigquery
bazel run bechmark:go_default_test -- -test.bench=.
```
### COMMANDS

1. bazel run //:gazelle
2. bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 producer:latest
3. docker run -it bazel/producer:latest
4. helm install eventproducer producer/deploy/helm-chart

### NOTE:
1. generated protobuf files from bazel-bin folder copied/linked to proto folder, required only if you want to change proto.
2. swagger json will be generated will present on bazel-bin folder to be copied and view in swagger-editor here https://editor.swagger.io/
3. while running services on minikube to access httpserver or gRPC server on eventconsumer service need expose port using below command
```bigquery
kubectl --namespace default port-forward $POD_NAME 8080:$CONTAINER_PORT
```
  