package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/305983806/gotunnel/util"
	"io"
	"net"
)

func NewClient() {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8100")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Client connect error ! " + err.Error())
		return
	}
	fmt.Printf("已成功连接到服务器端[%s][local:%s]", conn.RemoteAddr().String(), conn.LocalAddr().String())
	reader := bufio.NewReader(conn)

	// 向服务器登记tunnel
	conf := util.Conf{}
	err = util.getYaml("./server_config.json", &conf)
	if err != nil {
		fmt.Println(err)
	}
	jsonStr, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}
	_, err = conn.Write([]byte(string(jsonStr) + "\n"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("已经向中央服务器成功登记 Tunnel[%s]", conf.Server.Tunnel)

	fmt.Printf("从服务端获得心跳包的发送时间间隔为：2秒")
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

func main() {
	NewClient()
}