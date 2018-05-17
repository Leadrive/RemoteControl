package main
import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)
func main() {
	command := "ping"
	params := []string{"-a", "127.0.0.1"}
	//执行cmd命令: ls -l
	execCommand(command, params)
	//    command := "ipconfig"
	//    params := []string{"/all"}
	//    //执行cmd命令: ls -l
	//    ip := getip(command, params)
	//    fmt.Println(ip)
	//    ip2 := IncIP(ip, 1)
	//    ip3 := IncIP(ip, 2)
	//    ip4 := IncIP(ip, 3)
	//    ip5 := IncIP(ip, 4)
	//    fmt.Println(ip2)
	//    fmt.Println(ip3)
	//    fmt.Println(ip4)
	//    fmt.Println(ip5)
}
func IncIP(ip string, n int) string {
	ips := strings.Split(ip, ".")
	ip3, error := strconv.Atoi(ips[3])
	if error != nil {
		fmt.Println("字符串转换成整数失败")
	}
	ip3 = ip3 + n
	ip3_d := strconv.Itoa(ip3) //数字变成字符串
	ip_new := ips[0] + "." + ips[1] + "." + ips[2] + "." + ip3_d
	//fmt.Println(ip_new)
	return ip_new
}
func execCommand(commandName string, params []string) bool {
	cmd := exec.Command(commandName, params...)
	//显示运行的命令
	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		//        enc := mahonia.NewDecoder("UTF-8")
		//        goStr := enc.ConvertString(line)
		fmt.Println(line)
	}
	cmd.Wait()
	return true
}
func getip(commandName string, params []string) string {
	cmd := exec.Command(commandName, params...)
	//显示运行的命令
	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	result := ""
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		//fmt.Println(line)
		reg := regexp.MustCompile("\\d{1,3}.\\d{1,3}.\\d{1,3}.\\d{1,3}") //六位连续的数字
		//返回str中第一个匹配reg的字符串
		data := reg.Find([]byte(line))
		if data != nil {
			result = string(data)
			//fmt.Println(result)
			break
		}
	}
	cmd.Wait()
	return result
}