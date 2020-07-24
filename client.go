package gotunnel

import (
	"bufio"
	"fmt"
	"io"
	"net"
)
func Initialization() {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8100")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Client connect error ! " + err.Error())
		return
	}
	fmt.Println("成功接入中央服务器。")
	reader := bufio.NewReader(conn)
	for {
		s, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		} else {
			if s == "new\n" {
				// 创建新的TCP链路
				go combine()
			} else {
				// 其他一切消息忽略，如：hi心跳包
			}
		}
	}
}

func combine() {
	local := localConn()
	tunnel := tunnelConn()
	defer close(local, "local")
	defer close(tunnel, "tunnel")
	if local != nil && tunnel != nil {
		joinConn(local, tunnel)
	}
}

func joinConn(local, tunnel *net.TCPConn) {
	f := func(dst *net.TCPConn, src *net.TCPConn) {
		_, err := io.Copy(dst, src)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("end")
	}
	go f(local, tunnel)
	go f(tunnel, local)
}

func localConn() *net.TCPConn {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("本地服务：8080 连接失败。")
		return nil
	}
	fmt.Println("本地服务：8080 连接成功")
	return conn
}

func tunnelConn() *net.TCPConn {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:8101")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("tunnel：8101 连接失败")
		return nil
	}
	fmt.Println("tunnel: 8101 连接成功")
	return conn
}

func close(conn *net.TCPConn, name string) {
	err := conn.Close()
	if err != nil {
		fmt.Println(name + " close: " + err.Error())
	}
}