package subcriber
import(
	"github.com/qianlidongfeng/log"
)
type subInfo struct{
	g string
	h func([]byte)
}
type SubQueue struct{
	mSubInfo map[string]subInfo
}

func NewSubQueue() SubQueue{
	return SubQueue{
		mSubInfo:make(map[string]subInfo),
	}
}

func (this *SubQueue) Register(subject string,group string,handle func([]byte)){
	this.mSubInfo[subject]=subInfo{
		g:group,
		h:handle,
	}
}

func (this *SubQueue) Init(subcriberNum uint64,subAddr string){
	var i uint64
	for ;i<subcriberNum;i++{
		subcriber:=NewSubcriber()
		err:=subcriber.Connect(subAddr)
		if err != nil{
			log.Warn(err)
		}
		for k,v:=range(this.mSubInfo){
			subcriber.RegisterGroup(k,v.g,v.h)
		}
	}
}