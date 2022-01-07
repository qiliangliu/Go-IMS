package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的 channel
	Message chan string
}

//创建一个server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

//广播消息方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

//处理当前链接的业务
func (this *Server) Handler(conn net.Conn) {
	//...当前链接的业务
	// fmt.Println("链接建立成功...")

	//创建一个用户对象
	user := NewUser(conn)

	//User上线，将用户加到 OnlineMap 中去
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	//广播当前用户上线
	this.BroadCast(user, "已上线")

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf) //n 为读入数据的长度
			if n == 0 {              //当客户端退出的时候回返回一个0值
				this.BroadCast(user, "下线")
				this.mapLock.Lock()
				delete(this.OnlineMap, user.Name) //用户下线删除在OnlineMap中的内容
				this.mapLock.Unlock()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			//提取用户消息（去除'\n'）
			msg := string(buf[:n-1])
			//讲用户的消息进行广播
			this.BroadCast(user, msg)
		}
	}()

	//阻塞 handler 结束
	// select {}
}

//监听Message管道中的一个gorountine, 一旦这个管道中有消息，我们就把这个消息广播到OnlineMap中所用用户
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//将所有的消息推送给在线用户
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

//启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	//close listen socket
	defer listener.Close()

	//启动监听进程
	go this.ListenMessager()

	//accept
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}

}
