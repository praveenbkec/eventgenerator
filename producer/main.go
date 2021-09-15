package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/praveenbkec/eventgenerator/producer/pkg"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func init() {
	createTopic()
}

func main() {
	fmt.Println("***** Starting producer ******")
	produceEvents()
}

const (
	brokerAddress = "messaging-kafka-0.messaging-kafka-headless.default.svc.cluster.local:9092"
	topic = "message-event"
	partion = 0
)

// Event = Name:XXXXX,Dept=OSS,EmplD:1234, Time=21-7-2021 21:00:10
func produceEvents() {
	employess := []string{"12345", "23456", "34567", "45678"}
	fmt.Println("********** inside producer ******** ")
	w := kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	j:=0
	for _ = range time.NewTicker(30 * time.Second).C {
		empEventProducer := &pkg.EmployeeAccessEvent{
			EventType: pkg.EmployeeAccessEventConst,
			EmpID:     employess[j],
		}
		event, _ := empEventProducer.ProduceEvent()
		log.Printf("============= Creating event ============= %v\n", j)
		if j == 3 { j=0 } else { j++ }
		fmt.Println(event)
		eventJson, errMarshall := json.Marshal(event)
		if errMarshall != nil {
			log.Fatal("unable to marshall event to json object")
		}
		err := w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte("Event"),
				Value: eventJson,
			},
		)
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}
	}
	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func createTopic()  {
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokerAddress, topic, partion)
	if err != nil {
		panic(err)
	}
	// close the connection because we won't be using it
	conn.Close()
}