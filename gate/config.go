package gate

import (
	"flag"
	"io/ioutil"
	"encoding/json"
	"github.com/qianlidongfeng/log"
	"path/filepath"
	"os"
)

var(
	g_config Config
)

func init(){
	g_config = Config{}
	g_config.Init()
}

type Config struct{
	ServerID int64
	TcpPort string
	WebsocketPort string
	MaxConnection uint64
	WirteQueLen uint64
	MaxMsgLen uint32
	BigEndian bool
	Distributed bool
	PublisherNum uint64
	SubcriberNum uint64
	PublishAddr string
	SubcribAddr string
}


func (this *Config)Init(){
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Warn(err)
	}
	configFile:=flag.String("c",dir+"/gate_config.json","config file path")
	flag.Parse()
	data,err := ioutil.ReadFile(*configFile)
	if err != nil{
		log.Warn(err)
	}
	var mConfig map[string]interface{}
	err = json.Unmarshal(data,&mConfig)
	if err != nil{
		log.Warn(err)
	}
	if serverID,ok:= mConfig["ServerID"];ok{
		this.ServerID = int64(serverID.(float64))
	}else{
		this.ServerID=1
	}
	if port,ok:=mConfig["TcpPort"];ok{
		this.TcpPort = port.(string)
	}else{
		this.TcpPort = "5956"
	}
	if port,ok:=mConfig["WebsocketPort"];ok{
		this.WebsocketPort = port.(string)
	}else{
		this.WebsocketPort = "5957"
	}
	if maxConnection,ok:=mConfig["MaxConnection"];ok{
		this.MaxConnection = uint64(maxConnection.(float64))
	}else{
		this.MaxConnection = 10000
	}
	if wirteQueLen,ok:=mConfig["WirteQueLen"];ok{
		this.WirteQueLen = uint64(wirteQueLen.(float64))
	}else{
		this.WirteQueLen = 300
	}
	if maxMsgLen,ok:=mConfig["MaxMsgLen"];ok{
		this.MaxMsgLen = uint32(maxMsgLen.(float64))
	}else{
		this.MaxMsgLen = 4096
	}
	if bigEndian,ok:=mConfig["BigEndian"];ok{
		this.BigEndian = bigEndian.(bool)
	}else{
		this.BigEndian = true
	}
	if distributed,ok:=mConfig["Distributed"];ok{
		this.Distributed=distributed.(bool)
	}else{
		this.Distributed=false
	}
	if publisherNum,ok:=mConfig["PublisherNum"];ok{
		this.PublisherNum = uint64(publisherNum.(float64))
	}else{
		publisherNum=8
	}
	if subcriberNum,ok:=mConfig["SubcriberNum"];ok{
		this.SubcriberNum = uint64(subcriberNum.(float64))
	}else{
		this.SubcriberNum=8
	}
	if publishAddr,ok := mConfig["PublishAddr"];ok{
		this.PublishAddr = publishAddr.(string)
	}else{
		this.PublishAddr = "nats://localhost:4222"
	}
	if subcribAddr,ok := mConfig["SubcribAddr"];ok{
		this.SubcribAddr = subcribAddr.(string)
	}else{
		this.SubcribAddr = "nats://localhost:4222"
	}
}