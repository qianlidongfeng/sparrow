package gate

import (
	"net"
	"github.com/qianlidongfeng/log"
	"github.com/qianlidongfeng/netserver"
	"github.com/qianlidongfeng/sparrow/mq/nats/publisher"
	"github.com/qianlidongfeng/sparrow/mq/nats/subcriber"
	"github.com/qianlidongfeng/sparrow/uuid"
	Msg "github.com/qianlidongfeng/sparrow/msg"
	"github.com/golang/protobuf/proto"
	"strings"
	"sync"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)


type GateServer struct{
	conns map[int64]Conn
	router Router
	maxConnectionlimiter chan struct{}
	tcpServer netserver.TcpServer
	websocketServer netserver.WebsocketServer
	gateSubject string
	snowflake *uuid.Snowflake
	pubQueue publisher.PubQueue
	subQueue subcriber.SubQueue
	OnConnect func(sid int64)
	OnDisconnect func(sid int64)
	mu sync.Mutex
}

func NewServer() GateServer{
	server:=GateServer{
		conns:make(map[int64]Conn),
		router:NewRouter(),
		maxConnectionlimiter:make(chan struct{},g_config.MaxConnection),
		tcpServer:netserver.NewTcpServer(),
		websocketServer:netserver.NewWebsocketServer(),
		pubQueue:publisher.NewPubQueue(),
		subQueue:subcriber.NewSubQueue(),
		OnConnect:func(sid int64){},
		OnDisconnect:func(sid int64){},
		}
	return server
}


func (this *GateServer) Run() error{
	var err error
	this.snowflake,err=uuid.NewSnowflake(g_config.ServerID)
	if err != nil{
		return err
	}
	this.gateSubject=strings.ToUpper(fmt.Sprintf("%x",this.snowflake.Generate()))
	if g_config.Distributed{
		this.subQueue.Register(this.gateSubject,this.gateSubject,func(msg []byte){
			tagLen := msg[0]
			tag := msg[1:tagLen+1]
			this.router.routeMqMsg(string(tag),this,msg[1+tagLen:])
		})
		this.subQueue.Register("00000000",this.gateSubject,func(msg []byte){
			tagLen := msg[0]
			tag := msg[1:tagLen+1]
			this.router.routeMqMsg(string(tag),this,msg[1+tagLen:])
		})
		this.pubQueue.Init(g_config.PublisherNum,g_config.PublishAddr)
		this.subQueue.Init(g_config.SubcriberNum,g_config.SubcribAddr)
	}
	this.tcpServer.OnAccept = this.accept
	go func(){
		err := this.tcpServer.Listen(":"+g_config.TcpPort)
		if err != nil{
			log.Warn(err)
		}
	}()
	this.websocketServer.OnAccept = this.accept
	go func(){
		err := this.websocketServer.Listen(":"+g_config.WebsocketPort)
		if err != nil{
			log.Warn(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt,os.Kill,syscall.SIGTERM)
	sig := <-c
	fmt.Printf("sparrow gate closing down (signal: %v)\n", sig)
	this.ShutDown()
	return nil
}


func (this *GateServer)RigsterClientMsg(tag string,handle func(gt *GateServer,sid int64,msg []byte)){
	this.router.registerClientMsg(tag,handle)
}

func (this *GateServer)RigsterMqMsg(tag string,handle func(gt *GateServer,msg []byte)){
	this.router.registerMqMsg(tag,handle)
}

func (this *GateServer)accept(conn net.Conn){
	this.limiterAcquire()
	localAddr:=conn.LocalAddr()
	port := strings.Split(localAddr.String(),":")[1]
	var c Conn
	if port == g_config.TcpPort{
		c = &TcpConn{connection:conn,sendQue:make(chan []byte,g_config.WirteQueLen)}
	}else if port == g_config.WebsocketPort{
		c = &WebsocketConn{connection:conn,sendQue:make(chan []byte,g_config.WirteQueLen)}
	}else{
		conn.Close()
		this.limiterReleaseOne()
		return
	}
	sessionID := this.snowflake.Generate()
	this.mu.Lock()
	this.conns[sessionID]=c
	this.mu.Unlock()
	go func (conn Conn){
		c.initWrite()
	}(c)
	this.OnConnect(sessionID)
	go func(sid int64){
		this.readMsg(sid)
		this.SessionClose(sid)
	}(sessionID)
}

func (this *GateServer) readMsg(sid int64){
	conn := this.conns[sid]
	for{
		msgData,err:= conn.read()
		if err != nil{
			break
		}
		tagLen:=msgData[0]
		tag := msgData[1:tagLen+1]
		if g_config.Distributed&&!this.router.routeClientMsg(string(tag),this,sid,msgData[1+tagLen:]){
			err := this.Publish(string(tag),sid,msgData[1+tagLen:])
			if err != nil{
				log.Warn(err)
			}
		}
	}
}

func (this *GateServer) Publish(subject string,sid int64,msg[]byte) error{
	msg2server,err:= proto.Marshal(&Msg.Msg2Server{GateSubject:this.gateSubject, SessionID:sid, Msg:msg})
	if err != nil{
		return err
	}
	this.pubQueue.Push(subject, msg2server)
	return nil
}

func (this *GateServer) WriteMsg(sid int64,tag string,msg []byte){
	conn,ok := this.conns[sid]
	if ok{
		tagLen := len(tag)
		msgLen := len(msg)
		length := tagLen+msgLen+1
		data:=make([]byte,length)
		data[0]=byte(tagLen)
		copy(data[1:],[]byte(tag))
		copy(data[1+tagLen:],msg)
		err:=conn.write(data)
		if err != nil{
			log.Warn(err)
		}
	}
}


func (this *GateServer) SessionClose(sid int64){
	conn,ok:= this.conns[sid]
	if ok{
		this.mu.Lock()
		delete(this.conns,sid)
		this.mu.Unlock()
		conn.close()
		this.OnDisconnect(sid)
		this.limiterReleaseOne()
	}
}

func (this *GateServer) RemoteAddr(sid int64) net.Addr{
	conn,ok:= this.conns[sid]
	if ok{
		return conn.remoteAddr()
	}else{
		return nil
	}
}

func (this *GateServer) limiterAcquire(){
	this.maxConnectionlimiter <- struct{}{}
}

func (this *GateServer) limiterReleaseOne(){
	<-this.maxConnectionlimiter
}

func (this *GateServer) ShutDown(){
	this.tcpServer.Close()
	this.websocketServer.Close()
	for k,_ := range(this.conns){
		this.SessionClose(k)
	}
}

