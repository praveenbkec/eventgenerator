package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/praveenbkec/eventgenerator/consumer/proto"
	"github.com/praveenbkec/eventgenerator/consumer/pkg/eventconsumer"
	//"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/lib/pq"
	"github.com/praveenbkec/eventgenerator/consumer/pkg/rpc"
	"google.golang.org/grpc"
	"net/http"
)

const (
	grpcPort = "10000"
	serverAddress = ":8080"
)

var (
	getEndpoint  = flag.String("get", "localhost:"+grpcPort, "endpoint of YourService")
	postEndpoint = flag.String("post", "localhost:"+grpcPort, "endpoint of YourService")

	swaggerDir = flag.String("swagger_dir", "template", "path to the directory which contains swagger definitions")
)

func init() {
	rpc.RegisterGrpcServer()
}

func main() {
	fmt.Println("getEndpoint"+*getEndpoint)
	go eventconsumer.ConsumeEvents()
	RunGrpcGateway(serverAddress)
}

func RunGrpcGateway(address string, opts ...runtime.ServeMuxOption) error {
	fmt.Println("********* Starting gRPC Gateway ************")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	gw, err := newGateway(ctx, opts...)
	if err != nil {
		return err
	}
	mux.Handle("/", gw)
	return http.ListenAndServe(address, mux)
}
//TODO Login/Log interceptor example
func newGateway(ctx context.Context, opts ...runtime.ServeMuxOption) (http.Handler, error) {
	fmt.Println("********* Creating newGateway ************")
	mux := runtime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	fmt.Println("********* RegisterEventGeneratorSvcHandlerFromEndpoint get ************")
	err := pb.RegisterEventGeneratorSvcHandlerFromEndpoint(ctx, mux, *getEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}

	fmt.Println("********* RegisterEventGeneratorSvcHandlerFromEndpoint post ************")
	err = pb.RegisterEventGeneratorSvcHandlerFromEndpoint(ctx, mux, *postEndpoint, dialOpts)
	if err != nil {
		return nil, err
	}

	return mux, nil
}



