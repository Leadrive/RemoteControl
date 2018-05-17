package main

import (
	"fmt"
	"net"

	"log"
	"os/exec"

	"regexp"

	"bufio"
	"io"
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

		analyseCommand(data, conn)
	}

}

func analyseCommand(data string, c net.Conn) {
	re := regexp.MustCompile(`"type":[^"]+"([^"]+)".+"command":[^"]+"([^"]+)"`)
	match := re.FindStringSubmatch(data)
	fmt.Println("match:", match[1:])
	cmdType, cmdString := match[1], match[2]

	if cmdType == "commandInConfig" {
		cmdAndarg := strings.Split(cmdString, " ")
		if len(cmdAndarg) > 1 {
			Cmd, arg := cmdAndarg[0], cmdAndarg[1]
			msg := execCommand(Cmd, arg)
			SendtoMaster(c, msg)
		} else {
			msg := execCommand(cmdAndarg[0], "")
			SendtoMaster(c, msg)
		}




	} else if cmdType == "commandInWrite" {
		writeCommand(cmdString)
	} else {
		log.Println("command type error! please input right command")
	}
}

//for example: python test.py

func execCommand(Cmd string, arg string) string{

	fmt.Printf("cmd: %s , arg: %s\n", Cmd, arg)

	cmd := exec.Command(Cmd, arg)
	fmt.Println("going to exe command...")

	//显示运行的命令
	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	msg := ""
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		msg += line
	}
	fmt.Println(msg)
	cmd.Wait()
	return msg
}

func writeCommand(cmdString string) {
	fmt.Println("going to write code...")

	f, err := os.Create("python.py")
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	//cmdString = strings.Replace(cmdString, `\n`, `\r\n`, -1)
	f.WriteString(cmdString)
	fmt.Println("writing code:", cmdString)
	fmt.Println("writen!")

}

func SendtoMaster(c net.Conn, msg string) {
	message := `{"to": "` + " " + `", "from": "` + "slave" + `", "command": "` + msg  + `", "type": "` + " " + `"}#finished#`
	fmt.Println("message to Master:", message)
	c.Write([]byte(message))
}