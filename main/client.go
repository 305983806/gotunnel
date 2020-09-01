package main

import (
	"bufio"
	"fmt"
	"github.com/305983806/gotunnel/util"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	serverHost = "127.0.0.1"
	serverPort = 8002
	tunnPort = 8003
	filePath = "../config.yaml"
)

var (
	config util.Config
)

func Start(ip string, port int) {
	addr, _ := net.ResolveTCPAddr("tcp", ip + ":" + strconv.Itoa(port))
	c, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("dial error: ", err)
		panic(err)
	}
	defer c.Close()
	sendConfig(c)
	receiveMsg(c)
}

// 发送配置信息，注册tunnel
func sendConfig(conn *net.TCPConn) {
	//TODO 读取配置
	config.GetConfig(filePath)
	b := []byte("{\"name\":\"cp\",\"rules\":[{\"tag\":\"neo\",\"host\":\"127.0.0.1\",\"port\":8080}]}\n")
	conn.Write(b)
}

func receiveMsg(conn *net.TCPConn)  {
	// 循环读取连接数据：
	// a) 如果是 hi，表示心跳包，直接忽略
	// b) 如果是 new，则需要主动建立起本地与服务器的tunnel连接
	// c) 如果是 success，表示成功向中央服务器登记Tunnel
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		//fmt.Println("接收到消息：", msg)

		if msg == "success\n" {
			fmt.Println("已成功向中央服务器登记Tunnel")
		} else if strings.HasPrefix(msg, "success") {
			//TODO 建立 tunnel 隧道
			fmt.Println("准备建立隧道...")
			tag := strings.Replace(msg, "\n", "", -1)
			tag = strings.Replace(tag, "new_", "", -1)
			go combine(tag)
		} else if msg == "hi\n" {
			// 心跳检测，忽略
		}
	}
}

func combine(tag string) {
	remote := connectRemote()
	if remote != nil {

	}
}

func connectRemote() *net.TCPConn {
	addr, _ := net.ResolveTCPAddr("tcp", serverHost + ":" + strconv.Itoa(tunnPort))
	c, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("DialTCP error: ", err)
	}
	return c
}

func connectLocal() {

}

func main() {
	Start(serverHost, serverPort)
}