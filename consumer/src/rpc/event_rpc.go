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
	listen, err := net.Listen("tcp", ":10000")
	if err != nil {
		log.Fatal("Unable to listen on port 10000", listen)
	}
	gRpcSvr := grpc.NewServer()
	pb.RegisterEventGeneratorSvcServer(gRpcSvr, &server{})
}