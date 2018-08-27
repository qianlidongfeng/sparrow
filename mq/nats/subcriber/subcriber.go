package subcriber
import(
	"github.com/nats-io/go-nats"
)
type Subcriber struct{
	nc *nats.Conn
	urls string
}

func NewSubcriber() Subcriber{
	return Subcriber{}
}

func (this *Subcriber) Connect(urls string) error{
	var err error
	this.nc,err=nats.Connect(urls)
	return err
}

func (this *Subcriber) IsConnect() bool{
	return this.nc.IsConnected()
}

func (this *Subcriber) Close(){
	this.nc.Close()
}

func (this *Subcriber) IsReconnecting() bool{
	return this.nc.IsReconnecting()
}

func (this *Subcriber) Register(subject string,handle func([]byte)){
	this.nc.Subscribe(subject,func(msg *nats.Msg){
		handle(msg.Data)
	})
}

func (this *Subcriber) RegisterGroup(subject string,group string,handle func([]byte)){
	this.nc.QueueSubscribe(subject,group,func(msg *nats.Msg){
		handle(msg.Data)
	})
}