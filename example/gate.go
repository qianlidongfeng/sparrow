package main

import (
	"github.com/qianlidongfeng/log"
	"github.com/qianlidongfeng/sparrow/gate"
	Msg "github.com/qianlidongfeng/sparrow/msg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"golang.org/x/net/context"
	pb "github.com/qianlidongfeng/sparrow/example/proto"
)

func OnNormal(gt *gate.GateServer,sid int64,msg []byte){
	tag:="NORMAL"
	gt.WriteMsg(sid,Msg.MakeMsg(tag,[]byte(string(msg)+" I am sparrow")))
}

func GRPC(gt *gate.GateServer,sid int64,msg []byte){
	conn, err := grpc.Dial("127.0.0.1:50052", grpc.WithInsecure())

	if err != nil {
		grpclog.Fatalln(err)
	}

	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)
	reqBody := new(pb.HelloRequest)
	reqBody.Message = string(msg)
	// 调用方法
	r, err := c.SayHello(context.Background(), reqBody)
	if err != nil {
		grpclog.Warning(err)
		return
	}
	tag:="GRPC"
	gt.WriteMsg(sid,Msg.MakeMsg(tag,[]byte(r.Message+" I am sparrow")))
}

func MQ(gt *gate.GateServer,sid int64,msg []byte){
	tag:="MQ"
	gt.WriteMsg(sid,Msg.MakeMsg(tag,[]byte(string(msg)+" I am sparrow")))
}

func main(){
	gateServer:=gate.NewServer()
	gateServer.RigsterMsg("normal",OnNormal)
	gateServer.RigsterMsg("grpc",GRPC)
	gateServer.RigsterMsg("mq",MQ)
	err :=gateServer.Run()
	if err != nil{
		log.Fatal(err)
	}
}
