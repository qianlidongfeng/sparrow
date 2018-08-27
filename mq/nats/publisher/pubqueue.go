package publisher

import (
	"github.com/qianlidongfeng/sparrow/mq"
	"github.com/qianlidongfeng/log"
	"time"
)

type PubQueue struct{
	dataQue chan mq.PubData
}

func NewPubQueue() PubQueue{
	return PubQueue{}
}

func (this *PubQueue) Init(publisherNum uint64,pubAddr string){
	this.dataQue = make(chan mq.PubData,1024*1024)
	var i uint64
	for ;i<publisherNum;i++{
		go func(){
			publisher := NewPublisher()
			err:=publisher.Connect(pubAddr)
			if err != nil{
				log.Warn(err)
				return
			}
			for{
				pubData:=<-this.dataQue
				subject :=pubData.Subject
				msg:=pubData.Msg
				for{
					if publisher.IsConnect(){
						err:=publisher.Pub(subject,msg)
						if err != nil{
							log.Warn(err)
						}
					}else if publisher.IsReconnecting(){
						time.Sleep(time.Second)
						continue
					}
					break
				}
			}
		}()
	}
}

func (this *PubQueue)Push(subject string,msg []byte){
	this.dataQue<-mq.PubData{Subject:subject,}
}
