package main

import (
	"fmt"
	"net"

	"regexp"

	"log"

	"sync"

	"strings"
)

const (
	bufsize = 4096 * 100
	MASTER  = "master"
	SLAVE   = "slave"
	SERVER  = "server"
)

var (
	slaveConnPool []net.Conn
	masterConn    net.Conn
	mu            sync.Mutex
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:5000")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("========================\na new conn: ", conn.RemoteAddr().String())
		AddslaveConnPool(conn)

		go handleConn(conn)
	}
	defer listen.Close()

}

func handleConn(c net.Conn) {
	defer c.Close()
	for {

		buf := make([]byte, bufsize)
		nbyte, err := c.Read(buf)
		if err != nil {
			log.Printf("the socker has been closed by the %s", c.RemoteAddr().String())
			break
		}
		fmt.Println("recv data:", string(buf[:nbyte]))

		to, from, cmdString, cmdType := analyseCommand(string(buf[:nbyte]))
		fmt.Printf("from: %s, to: %s, cmd: %s, type: %s\n", from, to, cmdString, cmdType)

		//add to slavepool
		if from == MASTER {
			masterConn = c
			//delete the master conn from pool
			DelslaveConnPool(c)

			dispatch(to, from, cmdString, cmdType)

		} else {
			log.Println("don't know who send this message")
		}
	}
}

func AddslaveConnPool(c net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	slaveConnPool = append(slaveConnPool, c)
	fmt.Println("conn pool :", slaveConnPool)
}

func DelslaveConnPool(c net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	//fmt.Println("pool before delete", slaveConnPool)
	var pool []net.Conn
	for _, v := range slaveConnPool {
		if v == c {

		} else {
			pool = append(pool, v)
		}
	}
	slaveConnPool = pool
	//fmt.Println("pool after delete", slaveConnPool)
}

//return to, from, cmdString and cmdType
func analyseCommand(data string) (_, _, _, _ string) {
	//re := regexp.MustCompile(`"to"[^"]+"([^"]+)"[^"]+"from"[^"]+"([^"]+)"[^"]+"command"[^"]+"([^"+])"[^"]+"type"[^"]+"([^"]+)".+`)

	re := regexp.MustCompile(`"to"[^"]+"([^"]*)"[^"]+"from"[^"]+"([^"]*)"[^"]+"command"[^"]+"([^"]*)"[^"]+"type"[^"]+"([^"]*)"`)
	match := re.FindStringSubmatch(data)

	if len(match) < 5 {
		log.Println("not enough arguments! please input again...")
		return
	} else {
		return match[1], match[2], match[3], match[4]
	}
}

func dispatch(to, from, cmdString, cmdType string) {
	if to == SERVER || to == "" {
		cmdtoServer(cmdString, cmdType)
	} else {
		cmdtoSlave(to, from, cmdString, cmdType)
	}
}

func cmdtoServer(cmdString, cmdType string) {

	if cmdString == "listSlave" {
		slavelist := listSlave()
		fmt.Println("slavelist:", slavelist)
		SendtoMaster(slavelist)
	}
}

func cmdtoSlave(to, from, cmdString, cmdType string) {
	mu.Lock()
	defer mu.Unlock()

	for _, c := range slaveConnPool {
		if c.RemoteAddr().String() == strings.Replace(to, "-", ":", -1) {
			SendtoSlave(c, cmdString, cmdType)
		}
	}
}

func listSlave() []string {
	mu.Lock()
	defer mu.Unlock()
	var list []string
	for _, c := range slaveConnPool {
		list = append(list, c.RemoteAddr().String())
	}
	return list
}

func SendtoMaster(list []string) {
	//{"slaveList": ["127.0.0.1-7345", "127.0.0.1"]}#finished#
	message := `{"slaveList": ["`
	for k, v := range list {
		v = strings.Replace(v, ":", "-", -1)
		message += v + `"`
		if k < len(list)-1 {
			message += `, "`
		}
	}
	message += `]}#finished#`

	fmt.Println("message to master:", message)
	masterConn.Write([]byte(message))
}

func SendtoSlave(c net.Conn, cmdString, cmdType string) {
	message := `{"type": "` + cmdType +`", "command": "` + cmdString + `"}#finished#`
	fmt.Println("message to Slave:", message)
	c.Write([]byte(message))
}
