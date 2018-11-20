package main

import (
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	pb "github.com/qianlidongfeng/sparrow/example/proto"
)

const (
	Address = "127.0.0.1:50052"
)

var HelloService = helloService{}

type helloService struct{}


func (this helloService) SayHello(ctx context.Context,in *pb.HelloRequest)(*pb.HelloReply,error){
	resp := new(pb.HelloReply)
	resp.Message = in.Message
	return resp,nil
}


func main(){
	listen,err:=net.Listen("tcp",Address)
	if err != nil{
		grpclog.Fatalf("failed to listen: %v", err)
	}
	s:=grpc.NewServer()
	pb.RegisterHelloServer(s,HelloService)
	s.Serve(listen)
}