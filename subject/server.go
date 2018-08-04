package subject
import(
	"flag"
	"net/http"
	"github.com/qianlidongfeng/log"
	"os"
	"path/filepath"
	"io/ioutil"
	"encoding/json"
)
var(
	port *string
	dir string
	rmSubjects map[string]interface{}
	mSubjects map[string]interface{}
)


func init(){
	var err error
	port =flag.String("p","2580","subject server port")
	flag.Parse()
	dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadFile(dir+"/subject_conf/subject.json")
	if err != nil{
		log.Fatal(err)
	}
	mSubjects=make(map[string]interface{})
	rmSubjects=make(map[string]interface{})
	err = json.Unmarshal(data, &mSubjects)
	if err != nil{
		log.Fatal(err)
	}
	for k,v := range mSubjects{
		rmSubjects[v.(string)]=k
	}
}

type Server struct{
	port string
}


func NewServer() *Server{
	return &Server{*port}
}

func (this *Server) Run(){
	router := NewRouter()
	http.ListenAndServe(":"+this.port, router)
}