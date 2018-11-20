package gate

import (
	"sync"
)

type Router struct{
	table map[string]func(gt *GateServer,sid int64,msg []byte)
	mu sync.Mutex
}

func NewRouter() Router{
	return Router{
		table:make(map[string]func(gt *GateServer,sid int64,msg []byte)),
	}
}

func (this *Router) registerMsg(tag string,handle func(gt *GateServer,sid int64,msg []byte)){
	this.mu.Lock()
	this.table[tag] = handle
	this.mu.Unlock()
}

func (this *Router)routeMsg(tag string,gt *GateServer,sid int64,msg []byte) bool{
	if handle,ok:=this.table[tag];ok{
		handle(gt,sid,msg)
		return true
	}
	return false
}

