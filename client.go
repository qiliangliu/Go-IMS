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
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       -1,
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

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("4.查询在线用户")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 4 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>请输入合法范围内的数字<<<")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		//根据不同的模式处理不同的业务
		switch client.flag {
		case 1: //公聊模式
			fmt.Println("公聊模式")
			break
		case 2: //私聊模式
			fmt.Println("私聊模式")
			break
		case 3: //更改用户名
			fmt.Println("更改用户名")
			break
		case 4: //查询在线用户
			fmt.Println("查询在线用户")
			break
		}
	}
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
	client.Run()
}
