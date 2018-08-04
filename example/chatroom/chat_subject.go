package main
import (
	"github.com/qianlidongfeng/sparrow/subject"
)

func main(){
	server := subject.NewServer()
	server.Run()
}