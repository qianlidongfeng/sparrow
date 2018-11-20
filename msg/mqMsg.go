package msg

import (
	"github.com/golang/protobuf/proto"
)
func MakeMsg(tag string,data []byte) []byte{
	tagLen := len(tag)
	dataLen := len(data)
	b:=make([]byte,tagLen+dataLen+1)
	b[0]=byte(tagLen)
	copy(b[1:],tag)
	copy(b[1+tagLen:],data)
	return b
}

func MakeMqMsgToServer(gateSubject string,sid int64,data []byte)([]byte,error){
	return proto.Marshal(&MsgToServer{GateSubject:gateSubject, SessionID:sid, Data:data})
}

func UnMarshalMqMsgToServer(msg []byte)(gateSubject string,sid int64,data []byte,err error){
	info := &MsgToServer{}
	err=proto.Unmarshal(msg,info)
	gateSubject=info.GateSubject
	sid = info.SessionID
	data = info.Data
	return
}

func MakeMqMsgToGate(sid int64,tag string,data []byte)([]byte,error){
	return proto.Marshal(&MsgToGate{SessionID:sid,Tag:tag,Data:data})
}

func UnMarshalMqMsgToGate(msg[]byte)(sid int64,tag string,data []byte,err error){
	info := &MsgToGate{}
	err = proto.Unmarshal(msg,info)
	sid = info.SessionID
	tag = info.Tag
	data = info.Data
	return
}
