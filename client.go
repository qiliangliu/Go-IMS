package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	//链接sever服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil { //如果建立不成功，返回空nil
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn
	return client
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888		//通过这一个行在执行client可执行程序的时候加上 -ip 和 port 两个参数给对应的：serverIp和serverPort进行赋值
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器的Ip地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器的端口Port（默认是8888）")
	//命令行解析
	flag.Parse()
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>>>>>>>>链接建立失败<<<<<<<<<<<<")
	}
	fmt.Println(">>>>>>>>>>>链接建立成功<<<<<<<<<<<<")

	//启动客户端业务

	select {}
}
