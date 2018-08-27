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
func register(msg []byte){
	info := &Msg.Msg2Server{}
	err:=proto.Unmarshal(msg,info)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(string(info.Msg))
	data,err:=proto.Marshal(&pdb.Replysession{SessionID:info.SessionID,Msg:[]byte("fuck you,操你妈")})
	if err != nil{
		fmt.Println(err)
	}
	pubqueue.Push(info.GateSubject,Msg.MakeMsg("replysession",data))
}

func login(msg []byte){
	info := &Msg.Msg2Server{}
	err:=proto.Unmarshal(msg,info)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(string(info.Msg))
	data,err:=proto.Marshal(&pdb.Replysession{SessionID:info.SessionID,Msg:[]byte("fuck you,操你妈")})
	if err != nil{
		fmt.Println(err)
	}
	pubqueue.Push(info.GateSubject,Msg.MakeMsg("closesession",data))
}

func main(){
	subQueue:=subcriber.NewSubQueue()
	subQueue.Register("register","server001",register)
	subQueue.Register("login","server001",login)
	subQueue.Init(uint64(4),"nats://localhost:4222")
	runtime.Goexit()
}


