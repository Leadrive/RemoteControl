package main

import (
	"fmt"
	"net"
	//"regexp"
	//"os"
)

var (
	VoidCmd   []string //不用等回复的
	UnVoidCmd []string //要等回复的
	MasterCmd string
	to        string
	DATA      = make([]string, 4096)
)

func init() {
	UnVoidCmd = append(UnVoidCmd, "listSlave")
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:5000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 4096)
	for {
		//输入命令
		fmt.Println("please select your command: ")
		fmt.Scanln(&MasterCmd)
		fmt.Println("please select your slave (or server): ")
		fmt.Scanln(&to)

		for _, v := range UnVoidCmd {
			if MasterCmd == v {
				//sending
				SendMsg(conn, to, MasterCmd)

				fmt.Println("starting to recv msg...")
				length, _ := conn.Read(buf)
				//delete "#finished#"
				data := string(buf[0 : length-10])
				fmt.Println("recv data :", data)

				continue

			} else {
				SendMsg(conn, to, MasterCmd)
				continue
			}
		}

		//analyseCommand(data)

	}

}

func SendMsg(c net.Conn, to, cmdString string) {
	//{"to": "", "from": "master", "command": "listSlave", "type": ""}
	if to == "server" {
		to = ""
	}
	message := `{"to": "` + to + `", "from": "` + "master" + `", "command": "` + cmdString + `", "type": "` + "commandInConfig" + `"}#finished#`
	fmt.Println("message to Slave:", message)
	c.Write([]byte(message))
}
