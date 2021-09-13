package rpc

import (
	"context"
	"fmt"
	pb "github.com/praveenbkec/eventgenerator/consumer/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {

}

func(s *server) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.GetEventResponse, error)  {
	fmt.Println("********** invoked getEvent ******** ")
	return &pb.GetEventResponse{}, nil
}

func(s *server) ListEvent(ctx context.Context, req *pb.ListEventRequest) (*pb.ListEventResponse, error)  {
	fmt.Println("********** invoked ListEvent ******** ")
	return &pb.ListEventResponse{}, nil
}


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