package main

import (
	"context"
	"fmt"
	"github.com/odpf/shield/integration/fixtures/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

const (
	grpcBackendPort = 13877
	grpcProxyPort   = grpcBackendPort + 1
)

func startTestGRPCServer(port int, greetSrv helloworld.GreeterServer) error {
	s := grpc.NewServer()
	defer s.Stop()
	reflection.Register(s)

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		panic(err)
	}

	helloworld.RegisterGreeterServer(s, greetSrv)
	return s.Serve(lis)
}

type greetServer struct{}

func (s *greetServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *greetServer) StreamExample(in *helloworld.StreamExampleRequest, server helloworld.Greeter_StreamExampleServer) error {
	for i := 0; i < 10; i++ {
		if err := server.Send(&helloworld.StreamExampleReply{Data: fmt.Sprintf("foo-%d", i)}); err != nil {
			panic(err)
		}
	}
	return nil
}

type BasicAuthentication struct {
	Token string
}

func (a *BasicAuthentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", a.Token),
	}, nil
}

func (a *BasicAuthentication) RequireTransportSecurity() bool {
	return false
}
