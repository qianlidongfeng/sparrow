package main

import (
	"github.com/qianlidongfeng/sparrow/mq/nats/subcriber"
	"github.com/qianlidongfeng/sparrow/mq/nats/publisher"
	Msg "github.com/qianlidongfeng/sparrow/msg"
	"github.com/golang/protobuf/proto"
	pdb "github.com/qianlidongfeng/sparrow/example/chatroom/proto"
	"fmt"
	"runtime"
)
var pubqueue publisher.PubQueue
func init(){
	pubqueue.Init(uint64(4),"nats://localhost:4222")
}
func chat(msg []byte){
	info := &Msg.Msg2Server{}
	err:=proto.Unmarshal(msg,info)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(string(info.Msg))
	tag:="respone"
	content:="I am sparrow"
	tagLen := len(tag)
	contentLen := len(content)
	length := tagLen+contentLen+1
	data:=make([]byte,length)
	data[0]=byte(tagLen)
	copy(data[1:],[]byte(tag))
	copy(data[1+tagLen:],content)
	respone,err:=proto.Marshal(&pdb.Replysession{SessionID:info.SessionID,Msg:data})
	if err != nil{
		fmt.Println(err)
		return
	}
	pubqueue.Push(info.GateSubject,Msg.MakeMqMsg("mqMessage",respone))
}


func main(){
	subQueue:=subcriber.NewSubQueue()
	subQueue.Register("chat","server001",chat)
	subQueue.Init(uint64(4),"nats://localhost:4222")
	runtime.Goexit()
}


