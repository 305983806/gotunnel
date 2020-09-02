package main

import (
	"bufio"
	"encoding/json"
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
	filePath = "./config.yaml"
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
	err := config.GetConfig(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := json.Marshal(config)
	confStr := string(b) + "\n"
	fmt.Println(confStr)
	conn.Write([]byte(confStr))
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
			// 建立 tunnel 隧道
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
	var rule util.Rule
	var local *net.TCPConn = nil
	var remote = connectRemote()
	for _, rule = range config.Rules {
		if rule.Tag != tag {
			continue
		}
		local = connectLocal(rule.Host, rule.Port)
	}
	if remote != nil && local != nil {
		// 数据交换
		joinConn(local, remote)
	} else {
		if remote != nil {
			if err := remote.Close(); err != nil {
				fmt.Println("connection close error: " + err.Error())
			}
		}
		if local != nil {
			if err := local.Close(); err != nil {
				fmt.Println("connection close error: " + err.Error())
			}
		}
	}
}

func connectRemote() *net.TCPConn {
	addr, _ := net.ResolveTCPAddr("tcp", serverHost + ":" + strconv.Itoa(tunnPort))
	c, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("connectRemote -> DialTCP error: ", err)
	}
	return c
}

func connectLocal(host string, port int) *net.TCPConn {
	addr, _ := net.ResolveTCPAddr("tcp", host + ":" + strconv.Itoa(port))
	c, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("connectLocal -> dialTCP error: ", err)
	}
	return c
}

func joinConn(c1, c2 *net.TCPConn) {
	f := func(w, r *net.TCPConn) {
		defer w.Close()
		defer r.Close()
		_, err := io.Copy(w, r)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	go f(c1, c2)
	go f(c2, c1)
}

func main() {
	Start(serverHost, serverPort)
}