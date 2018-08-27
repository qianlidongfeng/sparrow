package gate

import (
	"net"
	"io"
	"github.com/juju/errors"
	"encoding/binary"
)

var(
	EndianGet16 func([]byte) uint16
	EndianGet32 func([]byte) uint32
	EndianPut16 func(b []byte, v uint16)
	EndianPut32 func(b []byte, v uint32)
)

func init(){
	if g_config.BigEndian{
		EndianGet16 = binary.BigEndian.Uint16
		EndianGet32 = binary.BigEndian.Uint32
		EndianPut16 = binary.BigEndian.PutUint16
		EndianPut32 = binary.BigEndian.PutUint32
	}else{
		EndianGet16 = binary.LittleEndian.Uint16
		EndianGet32 = binary.LittleEndian.Uint32
		EndianPut16 = binary.LittleEndian.PutUint16
		EndianPut32 = binary.LittleEndian.PutUint32
	}
}

type Conn interface{
	read() (msgData []byte,err error)
	write([]byte) error
	initWrite()
	remoteAddr() net.Addr
	close()
}

type TcpConn struct{
	connection net.Conn
	sendQue chan []byte
}

func (this *TcpConn) remoteAddr() net.Addr{
	return this.connection.RemoteAddr()
}

func (this *TcpConn) read() (msgData []byte,err error){
	bufMsgLen := make([]byte,4)
	var msgLen uint32
	_,err=io.ReadFull(this.connection,bufMsgLen)
	if err != nil{
		return nil,err
	}
	msgLen=EndianGet32(bufMsgLen)
	if msgLen>g_config.MaxMsgLen{
		return nil,errors.New("message is too long")
	}
	msgData = make([]byte, msgLen)
	if _, err := io.ReadFull(this.connection, msgData); err != nil {
		return nil,err
	}
	return
}

func (this *TcpConn) initWrite(){
	for{
		select{
		case data:= <-this.sendQue:
			if data==nil{
				return
			}
			this.connection.Write(data)
		}
	}
}

func (this *TcpConn) close(){
	this.sendQue<-nil
	this.connection.Close()
}

func (this *TcpConn) write(msg []byte) error{
	msgLen := len(msg)
	data:=make([]byte,4+msgLen)
	EndianPut32(data,uint32(msgLen))
	copy(data[4:],msg)
	select{
	case this.sendQue<-data:
		break
	default:
		return errors.New("sendque is full")
	}
	return nil
}

type WebsocketConn struct{
	connection net.Conn
	sendQue chan []byte
}

func (this *WebsocketConn) remoteAddr() net.Addr{
	return this.connection.RemoteAddr()
}

func (this *WebsocketConn) read() (msgData []byte,err error){
	b:=make([]byte,1)
	_,err=io.ReadFull(this.connection,b)
	if err != nil{
		return
	}
	fin:=((b[0] >> 7) & 1) != 0
	rsv1 := ((b[0] >> 6) & 1) != 0
	rsv2 := ((b[0] >> 5) & 1) != 0
	rsv3 := ((b[0] >> 4) & 1) != 0
	if rsv1 || rsv2 || rsv3{
		err = errors.New("websocket srv Error")
		return
	}
	opCode:=b[0] & 0x0f
	if opCode==0x8{
		err = errors.New("client closed")
		return
	}
	_,err = io.ReadFull(this.connection,b)
	if err != nil{
		return
	}
	mask := (b[0] & 0x80) != 0
	if mask == false{
		err = errors.New("closeStatusProtocolError")
		return
	}
	var payloadLen int64
	b[0] &= 0x7f
	switch {
	case b[0] <= 125: // Payload length 7bits.
		payloadLen = int64(b[0])
	case b[0] == 126: // Payload length 7+16bits
		payloadLenBuf:=make([]byte,2)
		_,err = io.ReadFull(this.connection,payloadLenBuf)
		if err !=nil{
			return
		}
		payloadLen = int64(uint16(payloadLenBuf[1]) | uint16(payloadLenBuf[0])<<8)
		//payloadLen = int64(binary.LittleEndian.Uint16(payloadLenBuf))
	case b[0] == 127: // Payload length 7+64bits
		payloadLenBuf:=make([]byte,8)
		_,err = io.ReadFull(this.connection,payloadLenBuf)
		if err !=nil{
			return
		}
		payloadLen = int64(uint64(payloadLenBuf[7]) | uint64(payloadLenBuf[6])<<8 |
			uint64(payloadLenBuf[5])<<16 | uint64(payloadLenBuf[4])<<24 |
				uint64(payloadLenBuf[3])<<32 | uint64(payloadLenBuf[2])<<40 |
					uint64(payloadLenBuf[1])<<48 | uint64(payloadLenBuf[0]&0x7f)<<56)
	}
	maskingByte := make([]byte, 4)
	_,err = io.ReadFull(this.connection,maskingByte)
	if err != nil{
		return
	}
	payloadData := make([]byte, payloadLen)
	_,err = io.ReadFull(this.connection,payloadData)
	if err != nil{
		return nil,err
	}
	for i := int64(0); i < payloadLen; i++ {
		payloadData[i] = payloadData[i] ^ maskingByte[i % 4]
	}
	if fin{
		msgData = payloadData
		if uint32(len(msgData)) > g_config.MaxMsgLen{
			err = errors.New("message is too long")
		}
		return
	}
	nextData, err := this.read()
	if err != nil{
		return
	}
	msgData = append(msgData, nextData...)
	return
}

func (this *WebsocketConn) write(msg []byte) error{
	length := len(msg)
	var header []byte
	var b byte
	b |= 0x80
	b |=2
	header = append(header, b)
	b = 0
	lengthFields := 0
	switch {
	case length <= 125:
		b |= byte(length)
	case length < 65536:
		b |= 126
		lengthFields = 2
	default:
		b |= 127
		lengthFields = 8
	}
	header=append(header, b)
	for i := 0; i < lengthFields; i++ {
		j := uint((lengthFields - i - 1) * 8)
		b = byte((length >> j) & 0xff)
		header = append(header, b)
	}
	select{
	case this.sendQue<-append(header,msg...):
		break
	default:
		return errors.New("sendque is full")
	}
	return nil
}

func (this *WebsocketConn) initWrite(){
	for{
		select{
		case data:= <-this.sendQue:
			if data==nil{
				return
			}
			this.connection.Write(data)
		}
	}
}


func (this *WebsocketConn) close(){
	this.sendQue<-nil
	this.connection.Close()
}