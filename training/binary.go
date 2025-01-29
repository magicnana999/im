package training

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func TestBinary() {

	{
		var writeBuf bytes.Buffer
		//var readBuf bytes.Buffer

		binary.Write(&writeBuf, binary.BigEndian, int32(1))
		fmt.Printf("Write int8  cap:%d,len:%d \n", writeBuf.Cap(), writeBuf.Len())

		binary.Write(&writeBuf, binary.BigEndian, int8(100))
		fmt.Printf("Write int8  cap:%d,len:%d \n", writeBuf.Cap(), writeBuf.Len())

		var int8value int8
		var int32value int32
		binary.Read(&writeBuf, binary.BigEndian, &int32value)
		fmt.Println("Read int32:", int32value)

		binary.Read(&writeBuf, binary.BigEndian, &int8value)
		fmt.Println("Read int8:", int8value)
	}

	{
		var writeBuf bytes.Buffer
		//var readBuf bytes.Buffer

		binary.Write(&writeBuf, binary.BigEndian, int32(len([]byte("hello"))))
		fmt.Printf("Write int8  cap:%d,len:%d \n", writeBuf.Cap(), writeBuf.Len())

		writeBuf.Write([]byte("hello"))
		fmt.Printf("Write int8  cap:%d,len:%d \n", writeBuf.Cap(), writeBuf.Len())

		var int32value int32
		binary.Read(&writeBuf, binary.BigEndian, &int32value)
		fmt.Println("Read int32:", int32value)

		v := make([]byte, int32value)
		i, _ := writeBuf.Read(v)
		fmt.Println("Read string:", i, string(v[:]))
	}

}
