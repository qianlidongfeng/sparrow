package msg
func MakeMsg(tag string,data []byte) []byte{
	tagLen := len(tag)
	dataLen := len(data)
	b:=make([]byte,tagLen+dataLen+1)
	b[0]=byte(tagLen)
	copy(b[1:],tag)
	copy(b[1+tagLen:],data)
	return b
}