package message

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"io"
)

type Message struct {
	From string
	To   string
	Type string
	Cmd  string
	Args string
}

func SendMsg(c net.Conn, m Message) {
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}
	c.Write(data)
}

func RecvMsg(c net.Conn, buf []byte) Message {
	nbyte, err := c.Read(buf)
	if err != nil {
		log.Fatal("recving msg error: ", err)
	}
	fmt.Println("recv data :", string(buf[:nbyte]))
	return AnalyseMsg(buf[:nbyte])
}

func AnalyseMsg(data []byte) Message {
	var m Message
	err := json.Unmarshal(data, &m)
	if err != nil {
		log.Fatalf("JSON unmarshaling failed: %s", err)
	}
	return m
}


func RecvFile(c net.Conn, fn string){
	f,err := os.Create(fn)
	if err!=nil{
		log.Fatal("open file error:", err)
	}
	finfo, _ := os.Lstat(fn)
	io.CopyN(f, c, finfo.Size())
}
func SendFile(c net.Conn, fn string){
	f,err := os.Open(fn)
	if err!=nil{
		log.Fatal("open file error:", err)
	}
	finfo, _ := os.Lstat(fn)
	io.CopyN(c, f, finfo.Size())
}

func TransFile(src, dst net.Conn){
	io.CopyN(src, dst, 4)
}