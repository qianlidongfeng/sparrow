package main

import (
	"github.com/qianlidongfeng/sparrow/mq/nats/subcriber"
	"github.com/qianlidongfeng/sparrow/mq/nats/publisher"
	Msg "github.com/qianlidongfeng/sparrow/msg"
	"github.com/qianlidongfeng/log"
	"runtime"
)
var pubqueue publisher.PubQueue
func init(){
	pubqueue.Init(uint64(4),"nats://localhost:4222")
}
func OnOther(msg []byte){
	gateID,sid ,data,err:=Msg.UnMarshalMqMsgToServer(msg)
	if err!= nil{
		log.Fatal(err)
	}
	tag:="mq"
	data,err=Msg.MakeMqMsgToGate(sid,tag,data)
	if err != nil{
		return
	}
	pubqueue.Publish(gateID,data)
}


func main(){
	subQueue:=subcriber.NewSubQueue()
	subQueue.Register("other","server001",OnOther)
	subQueue.Init(uint64(4),"nats://localhost:4222")
	runtime.Goexit()
}


