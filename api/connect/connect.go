package connect

import (
	"github.com/gorilla/websocket"
	"go-websocket/clientvar"
	"go-websocket/servers"
	"go-websocket/tools/util"
	"log"
	"net/http"
)

const (
	// 最大的消息大小
	maxMessageSize = 8192
)

type Controller struct {
	Upgrader *websocket.Upgrader
}

type toClient struct {
	ClientId string `json:"clientId"`
}

var defaultUpgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *Controller) Run(w http.ResponseWriter, r *http.Request) {
	c.Upgrader = defaultUpgrader
	conn, err := c.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	clientId := util.GenClientId()

	//给客户端绑定ID
	clientvar.AddClient(clientId, conn)

	//返回给客户端
	if err = conn.WriteJSON(toClient{ClientId: clientId}); err != nil {
		_ = conn.Close()
	}

	log.Printf("客户端已连接:%s 总连接数：%d", clientId, clientvar.ClientNumber())

	//设置读取消息大小上线
	conn.SetReadLimit(maxMessageSize)

	//发送心跳
	servers.SendJump(conn)

	//读取消息并发送 在这不提供
	//wh.readMessage(conn, clientId)

	//阻塞main线程
	select {}
}