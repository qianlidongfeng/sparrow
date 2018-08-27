package main

import (
	"os"
	"os/signal"
	"fmt"
	"syscall"
	"github.com/qianlidongfeng/log"
	"github.com/qianlidongfeng/sparrow/gate"
	"github.com/golang/protobuf/proto"
	pdb "github.com/qianlidongfeng/sparrow/example/chatroom/proto"
)

func onChatMessage(gt *gate.GateServer,sid int64,msg []byte){
	m := `
公司发动机公司公司及公司广东省可根据数据可是感觉是gdsfkgjsgsjgsjgkldjgsdljgdklgjdsggs公司发动机公司公司及公司广东省可根据数据可是感觉是gdsfkgjsgsjgsjgkldjgsdljgdklgjdsggs公司发动机公司公司及公司广东省可根据数据可是感觉是gdsfkgjsgsjgsjgkldjgsdljgdklgjdsggs公司发动机公司公司及公司广东省可根据数据可是感觉是gdsfkgjsgsjgsjgkldjgsdljgdklgjdsggs公司发动机公司公司及公司广东省可根据数据可是感觉是gdsfkgjsgsjgsjgkldjgsdljgdklgjdsggs公司发动机公司公司及公司广东省可根据数据可是感觉是gdsfkgjsgsjgsjgkldjgsdljgdklgjdsggs公司发动机公司公司及公司广东省可根据数据可是感觉是gdsfkgjsgsjgsjgkldjgsdljgdklgjdsggs公司发动机公司公司及公司广东省可根据数据可是感觉是gdsfkgjsgsjgsjgkldjgsdljgdklgjdsggs
`
	gt.WriteMsg(sid,"chatmessage",[]byte(m))
}
func onJoinroom(gt *gate.GateServer,sid int64,msg []byte){
	gt.WriteMsg(sid,"joinroom",msg)
}

func mqCloseSession(gt *gate.GateServer,msg[]byte){
	info := &pdb.Replysession{}
	proto.Unmarshal(msg,info)
	sid:=info.SessionID
	gt.SessionClose(sid)
}
func mqReplysession(gt *gate.GateServer,msg[]byte){
	info := &pdb.Replysession{}
	proto.Unmarshal(msg,info)
	sid:=info.SessionID
	data:=info.Msg
	gt.WriteMsg(sid,"hehe",data)
}

func main(){
	gateServer:=gate.NewServer()
	gateServer.RigsterClientMsg("chatmessage",onChatMessage)
	gateServer.RigsterClientMsg("joinroom",onJoinroom)
	gateServer.RigsterMqMsg("closesession",mqCloseSession)
	gateServer.RigsterMqMsg("replysession",mqReplysession)
	err :=gateServer.Run()
	if err != nil{
		log.Fatal(err)
	}
}
