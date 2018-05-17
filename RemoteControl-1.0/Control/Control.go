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
	cmdtype 	string
	cmdArgs string
	DATA      = make([]string, 4096)
)

func init() {
	UnVoidCmd = append(UnVoidCmd, "listSlave")
	UnVoidCmd = append(UnVoidCmd, "ping")
	UnVoidCmd = append(UnVoidCmd, "ipconfig")
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:5000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 4096)
	for {
		to =""
		cmdtype =""
		MasterCmd=""
		cmdArgs= ""

		//输入命令
		fmt.Println("please select your slave (or server): ")
		fmt.Scanln(&to)
		fmt.Println("please select your command type (commandInConfig or commandInWrite) : ")
		fmt.Scanln(&cmdtype)
		fmt.Println("please input your command or python code: ")
		fmt.Scanln(&MasterCmd)
		fmt.Println("please input your command args: ")
		fmt.Scanln(&cmdArgs)

		flag := false
		for _, v := range UnVoidCmd {
			if MasterCmd == v {
				//sending
				SendMsg(conn, to, cmdtype, MasterCmd, cmdArgs)

				fmt.Println("starting to recv msg...")
				length, _ := conn.Read(buf)
				//delete "#finished#"
				data := string(buf[0 : length-10])
				fmt.Println("recv data :", data)
				flag = true
				break
			}
		}
		if !flag {
			SendMsg(conn, to, cmdtype, MasterCmd, cmdArgs)
		}

		//analyseCommand(data)

	}

}

func SendMsg(c net.Conn, to, cmdType, cmdString , cmdArgs string) {
	//{"to": "", "from": "master", "command": "listSlave", "type": ""}
	if cmdArgs == "" {
		message := `{"to": "` + to + `", "from": "` + "master" + `", "command": "` + cmdString + `", "type": "` + cmdType + `"}#finished#`
		fmt.Println("message to Slave:", message)
		c.Write([]byte(message))
	}else {
		message := `{"to": "` + to + `", "from": "` + "master" + `", "command": "` + cmdString +" "+ cmdArgs + `", "type": "` + cmdType + `"}#finished#`
		fmt.Println("message to Slave:", message)
		c.Write([]byte(message))
	}


}
