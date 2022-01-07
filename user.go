package main

import (
	"net"
	"strings"
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
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//重命名消息格式：rename|张三
		newName := msg[7:]
		//判断newName是否已经存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名已经被占用!\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.Name = newName
			this.server.OnlineMap[this.Name] = this
			this.server.mapLock.Unlock()
			this.SendMsg("您已更新用户名为：" + this.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//消息格式：to|用户名|消息内容
		//1.获取对方的用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("消息格式不正确，请使用：to|用户名|消息内容\n")
			return
		}
		//2.根据用户名得到用户对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("当前用户不存在\n")
			return
		}
		//3.获取消息内容，并通过User对象发送过去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("无消息发送，请重发\n")
			return
		}
		remoteUser.SendMsg(this.Name + "对您说：" + content + "\n")
	} else {
		//将用户消息进行广播
		this.server.BroadCast(this, msg)
	}
}
