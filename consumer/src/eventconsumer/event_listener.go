package eventconsumer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

func ConsumeEvents() {
	fmt.Println("********* Starting Consumer ************")
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
		fmt.Println("\n ===================================== Event received ==========================================")
		fmt.Println(""+string(msg.Key)+ " : "+string(msg.Value))
		processEvent(msg)
	}

	//errC := r.Close()
	//if errC != nil {
	//	log.Fatal("unable to close reader ", errC)
	//}
}

func processEvent(msg kafka.Message) {
	eventObj := &EventRequest{}
	json.Unmarshal(msg.Value, eventObj)
	fmt.Println("Name:" + eventObj.Name + ", Dept:" + eventObj.Dept + ", EmpID:" + eventObj.EmpID + ", Time:" + eventObj.Time)
	eventMgmt := EventMgmtStruct{}
	ctx := context.Background()
	eventFromDb, err := eventMgmt.GetEvent(ctx, eventObj)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(eventFromDb, " eventFromDb")
	if eventFromDb != nil {
		fmt.Println("==== Update Event call ===== ")
		eventMgmt.UpdateEvent(ctx, eventObj)
	} else {
		fmt.Println("==== Create Event call ===== ")
		eventMgmt.CreateEvent(ctx, eventObj)
	}
}

