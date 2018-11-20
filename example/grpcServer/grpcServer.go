package main

import (
	"net"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	Address = "127.0.0.1:50052"
)

var HelloService = helloService{}

type helloService struct{}


func (this helloService) SayHello(ctx context.Context,in *HelloRequest)(*HelloReply,error){
	resp := new(HelloReply)
	resp.Message = in.Message
	return resp,nil
}


func main(){
	listen,err:=net.Listen("tcp",Address)
	if err != nil{
		grpclog.Fatalf("failed to listen: %v", err)
	}
	s:=grpc.NewServer()
	RegisterHelloServer(s,HelloService)
	s.Serve(listen)
	fmt.Println("1")
}