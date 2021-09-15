package main

import (
	"context"
	"fmt"
	pb "github.com/praveenbkec/eventgenerator/consumer/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	//GetEvent("12345")
	ListEvents()
}



func GetEvent(empId string)  {
	conn, err :=  grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	c := pb.NewEventGeneratorSvcClient(conn)
	req := &pb.GetEventRequest{EmpID: empId}
	resp, err := c.GetEvent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Event)
}

func ListEvents()  {
	conn, err :=  grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	c := pb.NewEventGeneratorSvcClient(conn)
	resp, err := c.ListEvent(context.Background(), &pb.ListEventRequest{})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response from server %v", resp.Events)
}
