package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

func NewServer()  {
	// 监听客户端链接8100端口
	startTCP8100()
	// 监听公网链接80端口
	startTCP80()
	// 如果接收到来自公网的请求，建立tunnel
	// 销毁tunnel
	// 通过tunnel转发tcp 2531.93
}

var (
	cache *net.TCPConn
	allConns = make(map[string]*net.TCPConn)
)

func startTCP8100() {
	const port = "8100"
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", ":" + port)
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("服务端启动失败：" + err.Error())
		panic(err)
	}
	fmt.Printf("Server started on port: %s\n", port)

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		fmt.Printf("客户端：%s 正在接入...\n", tcpConn.RemoteAddr().String())

		// 接收tunnel注册信息
		signinTunnel(tcpConn)

		_remoteAddr := tcpConn.RemoteAddr().String()
		if _, ok := allConns[_remoteAddr]; !ok {
			allConns[_remoteAddr] = tcpConn
		}
		go checkHeartbeat(tcpConn)
	}
}

// 登记tunnel
func signinTunnel(conn *net.TCPConn) {
	// 创建消息缓冲区
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}
	// 转化为字符串输出
	msg := string(buffer[0:n])
	fmt.Println("收到消息：", conn.RemoteAddr(), msg)

	conf := Conf{}
	err = json.Unmarshal([]byte(msg), &conf)
	if err != nil {
		panic(err)
	}
	// 向http注册一个目录
}

// 心跳检测
func checkHeartbeat(conn *net.TCPConn) {
	for {
		_, e := conn.Write([]byte("hi\n"))
		if e != nil {
			cache = nil
		}
		time.Sleep(time.Second * 2)
	}
}

// 打对外端口80
func startTCP80() {
	const port = "80"
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", ":" + port)
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	defer tcpListener.Close()
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("收到来自 %s:%s 的外部请求 ", tcpConn.RemoteAddr().String(), port)
	}
}

func main() {
	NewServer()
}
