package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server //当前用户所对应的服务器
}

//创建一个用户类
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动一个监听user channel 消息的goroutine
	go user.ListenMessage()

	return user
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

//用户上线功能
func (this *User) Online() {
	//User上线，将用户加到 OnlineMap 中去
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播当前用户上线
	this.server.BroadCast(this, "已上线")
}

//用户下线功能
func (this *User) Offline() {
	//User下线，将用户从 OnlineMap 中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//广播当前用户上线
	this.server.BroadCast(this, "下线")
}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

//用户消息处理功能
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + "在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else {
		//将用户消息进行广播
		this.server.BroadCast(this, msg)
	}
}
