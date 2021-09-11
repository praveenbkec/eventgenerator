package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

const (
	brokerAddress = "messaging-kafka.default.svc.cluster.local:9092"
	topic = "message-event"
	partion = 0
	batchSize = int(10e6)
)
func main() {
	consumeEventsNew()
}

type Event struct {
	Name string
	Dept string
	EmpID string
	Time string
}

func consumeEventsNew() {

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{brokerAddress},
		Topic:     topic,
		Partition: partion,
		MinBytes:  batchSize,
		MaxBytes:  batchSize,
	})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("unable to read message ", err)
			//break
		}
		fmt.Println("===================================== Event received ==========================================")
		fmt.Println(""+string(msg.Key)+ " : "+string(msg.Value))
		//var event map[string] interface{}
		//eventJson := json.Unmarshal(msg.Value, &event)
		//fmt.Println("eventJson",eventJson)
		eventObj:= &Event{}
		json.Unmarshal(msg.Value, eventObj)
		fmt.Println("Name:"+eventObj.Name+", Dept:"+eventObj.Dept+", EmpID:"+eventObj.EmpID+", Time:"+eventObj.Time)
	}

	//errC := r.Close()
	//if errC != nil {
	//	log.Fatal("unable to close reader ", errC)
	//}
}