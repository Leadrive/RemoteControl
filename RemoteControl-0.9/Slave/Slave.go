package main

import (
	"fmt"
	"net"

	"log"
	"os/exec"

	"regexp"

	"os"
	"strings"
)

var (
	DATA = make([]string, 4096)
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:5000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 4096)
	for {

		length, _ := conn.Read(buf)
		//delete "#finished#"
		data := string(buf[0 : length-10])
		fmt.Println("recv data :", data)

		analyseCommand(data)
	}

}

func analyseCommand(data string) {
	re := regexp.MustCompile(`"type":[^"]+"([^"]+)".+"command":[^"]+"([^"]+)"`)
	match := re.FindStringSubmatch(data)
	fmt.Println("match:", match[1:])
	cmdType, cmdString := match[1], match[2]

	if cmdType == "commandInConfig" {
		cmdAndarg := strings.Split(cmdString, " ")
		if len(cmdAndarg) > 1{
			Cmd, arg := cmdAndarg[0], cmdAndarg[1]
			exeCommand(Cmd, arg)
		}else {
			exeCommand(cmdAndarg[0], "")
		}

	} else if cmdType == "commandInWrite" {
		writeCommand(cmdString)
	} else {
		log.Println("command type error! please input right command")
	}
}

//for example: python test.py
func exeCommand(Cmd string, arg string) {
	fmt.Printf("cmd: %s , arg: %s\n", Cmd, arg)

	cmd := exec.Command(Cmd, arg)
	fmt.Println("going to exe command...")
	err := cmd.Start()
	if err != nil {
		log.Println(err)
	}
}

func writeCommand(cmdString string) {
	fmt.Println("going to write code...")

	f, err := os.Create("python.py")
	if err != nil{
		log.Println(err)
	}
	defer f.Close()


	//cmdString = strings.Replace(cmdString, `\n`, `\r\n`, -1)
	f.WriteString(cmdString)
	fmt.Println("writing code:", cmdString)
	fmt.Println("writen!")


}
