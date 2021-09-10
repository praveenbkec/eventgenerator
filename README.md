# EVENT GENERATOR

### Implement Event generator Service using microservice architecture and deploy minikube environment using docker and kubernetes tools

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

### PREREQUISITES  
* docker, minikube, bazel, go, kafka, patroni
* https://minikube.sigs.k8s.io/docs/start/
* https://helm.sh/docs/intro/install/

### COMMANDS

1. bazel run //:gazelle 
2. bazel run --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 producer/src/eventproducer:latest
3. docker run -it bazel/producer/src/eventproducer:latest


