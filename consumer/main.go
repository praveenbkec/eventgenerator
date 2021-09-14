package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/praveenbkec/eventgenerator/consumer/proto"
	"github.com/praveenbkec/eventgenerator/consumer/src/eventconsumer"
	//"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/lib/pq"
	"github.com/praveenbkec/eventgenerator/consumer/src/rpc"
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
	//return http.ListenAndServe(address, allowCORS(mux))
	return http.ListenAndServe(address, mux)
}
//TODO Login interceptor example
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

//TODO check whether required
//allowCORS allows Cross Origin Resoruce Sharing from any origin.
//Don't do this without consideration in production systems.
//func allowCORS(h http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if origin := r.Header.Get("Origin"); origin != "" {
//			w.Header().Set("Access-Control-Allow-Origin", origin)
//			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
//				preflightHandler(w, r)
//				return
//			}
//		}
//		h.ServeHTTP(w, r)
//	})
//}

//func preflightHandler(w http.ResponseWriter, r *http.Request) {
//	headers := []string{"Content-Type", "Accept"}
//	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
//	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
//	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
//	//glog.Infof("preflight request for %s", r.URL.Path)
//	fmt.Println("preflight request for %s", r.URL.Path)
//	return
//}


