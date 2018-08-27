package publisher

import (
	"github.com/nats-io/go-nats"
)
type Publisher struct{
	nc *nats.Conn
	urls string
}

func NewPublisher() Publisher{
	return Publisher{}
}


func (this *Publisher) Connect(urls string) error{
	var err error
	this.nc,err=nats.Connect(urls)
	return err
}


func (this *Publisher) IsConnect() bool{
	return this.nc.IsConnected()
}

func (this *Publisher) IsReconnecting() bool{
	return this.nc.IsReconnecting()
}

func (this *Publisher) Pub(subject string,msg []byte) error{
	this.nc.Publish(subject,msg)
	if err := this.nc.LastError(); err != nil {
		return err
	}
	return nil
}