package main

import (
	"context"
	"encoding/json"
	"fmt"
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

type Event struct {
	Name string
	Dept string
	EmpID string
	Time string
}


func produceEvents() {

	events := []Event{
		{
			Name:  "praveen",
			Dept:  "IT",
			EmpID: "12345",
		},
		{
			Name:  "rajesh",
			Dept:  "TEST",
			EmpID: "23456",
		},
		{
			Name:  "suresh",
			Dept:  "TEST",
			EmpID: "34567",
		},
		{
			Name:  "vinaya",
			Dept:  "IT",
			EmpID: "45678",
		},
	}
	fmt.Println("********** inside producer ******** ")
	w := kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	i := 3
	for i>=0 {
		event := events[i]
		event.Time = time.Now().Format(time.RFC850)
		fmt.Printf("creating event %s\n", i)
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
		fmt.Printf("creating complete %s\n", i)
		time.Sleep(30 * time.Second)
		if i ==0 {
			fmt.Println("************** Repeating count ****************")
			i = 3
		} else {
			i--
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