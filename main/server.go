package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/305983806/gotunnel/util"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Tunnel struct {
	util.Config
	ControlConn *net.TCPConn
	httpConn *net.TCPConn
}

var (
	httpPort = 8000
	acceptPort = 8001
	controlPort = 8002
	tunnelPort = 8003
	tunnels = make(map[string]*Tunnel)
	httpmaps = make(map[string]string)
	currTag = make(chan string, 1)
	lock = sync.Mutex{}
)

// http 反向代理
// @param port 监听端口
func HttpProxy(port int) {
	http.HandleFunc("/", ServeHttp)
	http.ListenAndServe(":" + strconv.Itoa(port), nil)
}

func ServeHttp(w http.ResponseWriter, r *http.Request) {
	path, err := getRawPath(r)
	if err != nil {
		fmt.Println(err)
	}
	r.URL.Path = path[2]

	currTag <- path[1]

	target, _ := url.Parse("http://127.0.0.1:" + strconv.Itoa(acceptPort))
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}

func getRawPath(r *http.Request) ([]string, error) {
	flysnowRegexp := regexp.MustCompile(`^/([0-9A-Za-z]+)([0-9_/-A-Za-z]+)$`)
	paths := flysnowRegexp.FindStringSubmatch(r.URL.Path)
	if len(paths) < 3 || paths[2] == "" {
		return nil, errors.New("url path error: invalid path")
	}
	return paths, nil
}

// 接收代理转发过来的 http 请求
// @param port 监听端口
func AcceptHttp(port int) {
	addr, _ := net.ResolveTCPAddr("tcp", ":" + strconv.Itoa(port))
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("listen TCP error: " + strconv.Itoa(port))
		panic(err)
	}
	defer l.Close()
	for {
		c, err := l.AcceptTCP()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}
		select {
		case tag := <- currTag:
			addHttpconnToTunnels(c, tag)
			sendMessage("new", tag)
		case <-time.After(time.Second):
			continue
		}
	}
}

func addHttpconnToTunnels(c *net.TCPConn, tag string) {
	//TODO 需要针对，httpmaps[tag] 对象为 nil 进行处理
	tunnelName := httpmaps[tag]
	tunnels[tunnelName].httpConn = c
}

func sendMessage(message string, tag string) {
	//TODO 需要针对，httpmaps[tag] 对象为 nil 进行处理
	tunnelName := httpmaps[tag]
	if tunnels[tunnelName] != nil {
		_, e := tunnels[tunnelName].ControlConn.Write([]byte(message + "_" + tag + "\n"))
		if e != nil {
			fmt.Println("消息 New 发送异常", e)
		}
	} else {
		fmt.Println("没有客户端连接，无法发送消息")
	}
}

// 启动中央服务器
// @param port 监听端口
func StartControl(port int) {
	addr, _ := net.ResolveTCPAddr("tcp", ":" + strconv.Itoa(port))
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("listen error: ", err)
		panic(err)
	}
	defer l.Close()
	for {
		c, err := l.AcceptTCP()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}
		go controlPipe(c)
	}
}

func controlPipe(conn *net.TCPConn) {
	defer func() {
		fmt.Println("Control -> Disconnected: " + conn.RemoteAddr().String())
		conn.Close()
	}()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil || err == io.EOF{
			break
		}
		fmt.Println(message)

		//TODO 接收客户端接入，并注册 tunnel
		var tunnel Tunnel
		str := strings.Replace(message, "\n", "", -1)
		err = json.Unmarshal([]byte(str), &tunnel)
		if err != nil {
			fmt.Println(err)
		}
		tunnel.ControlConn = conn
		tunnels[tunnel.Config.Tunnel] = &tunnel
		for _, v := range tunnel.Rules  {
			httpmaps[v.Tag] = tunnel.Name
		}

		//TODO 告诉客户端接入成功
		conn.Write([]byte("success\n"))

		//TODO 发送心跳包
		go sendHeartbeat(conn, tunnel.Name)
	}
}

// 发送心跳包
func sendHeartbeat(conn *net.TCPConn, name string) {
	for {
		_, e := conn.Write([]byte("hi\n"))
		if e != nil {
			delete(tunnels, name)
		}
		time.Sleep(time.Second * 2)
	}
}

// tunnel 监听服务，用于建立 tunnel隧道
// @param port 监听端口
func ServeTunnel(port int)  {

}

func main() {
	go HttpProxy(httpPort)
	go AcceptHttp(acceptPort)
	StartControl(controlPort)
}