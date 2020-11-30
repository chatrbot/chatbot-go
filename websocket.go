package chatbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/tidwall/gjson"

	"github.com/gorilla/websocket"
)

var (
	handshakeTimeout = 5 * time.Second
)

//监听WebSocket的推送消息
type WsServer struct {
	lock  sync.Mutex
	host  string
	token string
	con   *websocket.Conn
	//接收到消息后的插件调用队列
	plugins []Plugin
}

//新建WebSocket连接
func newWsServer(host, token string) (*WsServer, error) {
	con, err := connect(host, token)
	if err != nil {
		return nil, err
	}
	return &WsServer{
		con:     con,
		plugins: make([]Plugin, 0, 10),
		host:    host,
		token:   token,
	}, nil
}

//connect 建立WebSocket连接
func connect(host, token string) (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "ws",
		Host:   host,
		Path:   "/ws",
	}
	query, err := url.ParseQuery("token=" + token)
	if err != nil {
		return nil, err
	}
	u.RawQuery = query.Encode()
	log.Printf("connecting to %s", u.String())

	dialer := websocket.Dialer{
		HandshakeTimeout: handshakeTimeout,
	}
	c, rsp, err := dialer.Dial(u.String(), nil)
	if err != nil {
		if rsp != nil {
			body, _ := ioutil.ReadAll(rsp.Body)
			rsp.Body.Close()
			return nil, fmt.Errorf("%s:%s",
				err,
				gjson.ParseBytes(body).Get("msg").String(),
			)
		}
		return nil, err
	}
	return c, nil
}

//Listen 开始监听服务端消息和调用插件
func (ws *WsServer) Listen() {
	for {
		_, msg, err := ws.con.ReadMessage()
		if err != nil {
			log.Println("read msg err:", err)
			log.Println("连接断开,开始重连...")
			ws.reConnectWebSocket()
			continue
		}
		log.Println("收到消息:", string(msg))
		if len(ws.plugins) > 0 {
			var rec PushMessage
			_ = json.Unmarshal(msg, &rec)

			for _, p := range ws.plugins {
				if err := p.Do(&rec); err != nil {
					log.Printf("%s handle error:%s \n", p.Name(), err)
				}
			}
		}
	}
}

//reConnectWebSocket 断线重连
func (ws *WsServer) reConnectWebSocket() {
	for {
		if ws.con != nil {
			_ = ws.con.Close()
			ws.con = nil
		}
		con, err := connect(ws.host, ws.token)
		if err == nil {
			ws.con = con
			return
		}
		log.Println("重连WebSocket失败:", err, ",5s后重试")
		time.Sleep(5 * time.Second)
	}
}

//addPlugin 添加插件
func (ws *WsServer) addPlugin(plugin ...Plugin) {
	ws.plugins = append(ws.plugins, plugin...)
}

//writeMessage gorilla的WebSocket默认发送消息会有并发问题
func (ws *WsServer) writeMessage(message string) error {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	return ws.con.WriteMessage(websocket.TextMessage, []byte(message))
}

//关闭WebSocket连接
func (ws *WsServer) Close() error {
	return ws.con.Close()
}
