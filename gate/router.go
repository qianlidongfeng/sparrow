package gate

import (
	"sync"
)

type Router struct{
	clientTable map[string]func(gt *GateServer,sid int64,msg []byte)
	mqTable map[string]func(gt *GateServer,msg []byte)
	mu sync.Mutex
}

func NewRouter() Router{
	return Router{
		clientTable:make(map[string]func(gt *GateServer,sid int64,msg []byte)),
		mqTable:make(map[string]func(gt *GateServer,msg []byte)),
	}
}

func (this *Router) registerClientMsg(tag string,handle func(gt *GateServer,sid int64,msg []byte)){
	this.mu.Lock()
	this.clientTable[tag] = handle
	this.mu.Unlock()
}

func (this *Router)routeClientMsg(tag string,gt *GateServer,sid int64,msg []byte) bool{
	if handle,ok:=this.clientTable[tag];ok{
		handle(gt,sid,msg)
		return true
	}
	return false
}

func (this *Router) registerMqMsg(tag string,handle func(gt *GateServer,msg []byte)){
	this.mu.Lock()
	this.mqTable[tag] = handle
	this.mu.Unlock()
}

func (this *Router)routeMqMsg(tag string,gt *GateServer,msg []byte) bool{
	if handle,ok:=this.mqTable[tag];ok{
		handle(gt,msg)
		return true
	}
	return false
}