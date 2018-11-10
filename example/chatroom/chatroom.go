package main

import (
	"github.com/qianlidongfeng/log"
	"github.com/qianlidongfeng/sparrow/gate"
	"github.com/golang/protobuf/proto"
	pdb "github.com/qianlidongfeng/sparrow/example/chatroom/proto"
)

func onLocalMessage(gt *gate.GateServer,sid int64,msg []byte){
	tag:="local"
	content:="I am sparrow"
	tagLen := len(tag)
	contentLen := len(content)
	length := tagLen+contentLen+1
	data:=make([]byte,length)
	data[0]=byte(tagLen)
	copy(data[1:],[]byte(tag))
	copy(data[1+tagLen:],content)
	gt.WriteMsg(sid,data)
}


func onMqMessage(gt *gate.GateServer,msg[]byte){
	info := &pdb.Replysession{}
	proto.Unmarshal(msg,info)
	sid:=info.SessionID
	data:=info.Msg
	gt.WriteMsg(sid,data)
}

func main(){
	gateServer:=gate.NewServer()
	gateServer.RigsterLocalMsg("localMessage",onLocalMessage)
	gateServer.RigsterMqMsg("mqMessage",onMqMessage)
	err :=gateServer.Run()
	if err != nil{
		log.Fatal(err)
	}
}
