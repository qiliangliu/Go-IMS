package main

import (
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

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>>>>>>>>链接建立失败<<<<<<<<<<<<")
	}
	fmt.Println(">>>>>>>>>>>链接建立成功<<<<<<<<<<<<")

	//启动客户端业务

	select {}
}
