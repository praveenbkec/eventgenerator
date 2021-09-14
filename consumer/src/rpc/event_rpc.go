package rpc

import (
	"context"
	"fmt"
	pb "github.com/praveenbkec/eventgenerator/consumer/proto"
	"github.com/praveenbkec/eventgenerator/consumer/src/eventconsumer"
	"google.golang.org/grpc"
	"log"
	"net"
)


type server struct {

}

func(s *server) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.GetEventResponse, error)  {
	fmt.Println("********** invoked getEvent ******** ")
	eventMgmt := eventconsumer.EventMgmtStruct{}
	eventsResp, err := eventMgmt.GetEvent(ctx, &eventconsumer.EventRequest{EmpID: req.EmpID})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(eventsResp)
	getEventResp:= &pb.GetEventResponse{
		Event: &pb.Event{
			EmpID: eventsResp.EmpID,
			Name: eventsResp.Name,
			Dept: eventsResp.Dept,
			Time: eventsResp.Time,
		},
	}
	return getEventResp, nil
}

func(s *server) ListEvent(ctx context.Context, req *pb.ListEventRequest) (*pb.ListEventResponse, error)  {
	fmt.Println("********** invoked ListEvent ******** ")
	eventMgmt := eventconsumer.EventMgmtStruct{}
	eventsDB, err := eventMgmt.ListEvent(ctx, &eventconsumer.EventRequest{EmpID: ""})
	events := []*pb.Event{}
	for _, eventFromDb := range eventsDB {
		eventPb := &pb.Event{
			EmpID: eventFromDb.EmpID,
			Name: eventFromDb.Name,
			Dept: eventFromDb.Dept,
			Time: eventFromDb.Time,
		}
		events = append(events, eventPb)
	}
	if err != nil {
		return nil, err
	}
	fmt.Println(events)
	return &pb.ListEventResponse{Events: events}, nil
}

// TODO tls certs

func RegisterGrpcServer() {
	fmt.Println("********** RegisterGrpcServer ******** ")
	addr := ":10000"
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Unable to listen on port 10000", listen)
	}
	gRpcSvr := grpc.NewServer()
	fmt.Println("********** RegisterEventGeneratorSvcServer ******** ")
	pb.RegisterEventGeneratorSvcServer(gRpcSvr, &server{})
	// Register reflection service on gRPC server.
	//reflection.Register(gRpcSvr)
	fmt.Println("********** GrpcServer gRpcSvr.Serve before ******** ")
	//if err := gRpcSvr.Serve(listen); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
	// Serve gRPC Server
	fmt.Println("Serving gRPC on https://", addr)
	go func() {
		log.Fatal(gRpcSvr.Serve(listen))
	}()
	fmt.Println("********** GrpcServer Listening on 10000 ******** ")
}