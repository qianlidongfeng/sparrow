package main
import(
	"net"
	"encoding/binary"
	"io"
	"fmt"
)



func main(){
	conn,err:=net.Dial("tcp","127.0.0.1:5956")
	if err != nil {
		panic(err)
	}
	tag:="grpc"
	//tag:="message2"
	data:="hello"
	taglen:=len(tag)
	datalen:=len(data)
	msg:=make([]byte,taglen+datalen+5)
	binary.LittleEndian.PutUint32(msg, uint32(len(data)+len(tag)+1))
	//binary.BigEndian.PutUint32(msg, uint32(len(data)+len(tag)+1))
	msg[4]=byte(taglen)
	copy(msg[5:],tag)
	copy(msg[5+len(tag):],data)
	conn.Write(msg)
	lenBuf:=make([]byte,4)
	io.ReadFull(conn,lenBuf)
	var length uint32
	length = binary.LittleEndian.Uint32(lenBuf)
	msg=make([]byte,length)
	io.ReadFull(conn,msg)
	taglen=int(msg[0])
	tag=string(msg[1:taglen+1])
	data=string(msg[1+taglen:])
	fmt.Println(string(tag))
	fmt.Println(string(data))
}
