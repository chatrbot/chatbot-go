package chatbot

import (
	"encoding/json"
	"errors"
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
	// 握手建立超时时间
	handshakeTimeout = 5 * time.Second
)

// 监听WebSocket的推送消息
type WsServer struct {
	mu      sync.Mutex
	host    string
	token   string
	con     *websocket.Conn
	plugins []Plugin

	pingTimer *time.Timer
}

// 新建WebSocket连接
func newWSClient(host, token string) (*WsServer, error) {
	con, err := connect(host, token)
	if err != nil {
		return nil, err
	}
	log.Println("connect server success")
	server := &WsServer{
		con:     con,
		plugins: make([]Plugin, 0, 10),
		host:    host,
		token:   token,
	}
	server.startHeartBeat()
	return server, nil
}

// connect 建立WebSocket连接
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
	log.Println("connecting to", u.String())

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

// ReceiveCallbackMessage 开始监听服务端消息和调用插件
func (ws *WsServer) ReceiveCallbackMessage() {
	for {
		msgType, msg, err := ws.con.ReadMessage()
		if err != nil {
			log.Println("连接断开,开始重连...")
			ws.Close()
			ws.reconnect()
			ws.startHeartBeat()
			log.Println("重连成功")
			continue
		}
		if string(msg) == "pong" {
			continue
		}
		if msgType == websocket.TextMessage {
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
}

// reconnect 断线重连
func (ws *WsServer) reconnect() {
	for {
		con, err := connect(ws.host, ws.token)
		if err == nil {
			ws.mu.Lock()
			ws.con = con
			ws.mu.Unlock()
			return
		}
		log.Println("连接WebSocket失败:", err, ",5s后重试")
		time.Sleep(5 * time.Second)
	}
}

// startHeartBeat 心跳包
func (ws *WsServer) startHeartBeat() {
	log.Println("开始发送心跳包")
	ws.pingTimer = time.AfterFunc(10*time.Second, func() {
		if err := ws.ping(); err != nil {
			ws.Close()
		}
		ws.mu.Lock()
		if ws.pingTimer != nil {
			ws.pingTimer.Reset(10 * time.Second)
		}
		ws.mu.Unlock()
	})
}

// addPlugin 添加插件
func (ws *WsServer) addPlugin(plugin ...Plugin) {
	ws.plugins = append(ws.plugins, plugin...)
}

func (ws *WsServer) ping() error {
	return ws.writeMessage("ping")
}

func (ws *WsServer) writeMessage(message string) error {
	return ws.write(websocket.TextMessage, []byte(message))
}

func (ws *WsServer) write(messageType int, message []byte) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	if ws.con == nil {
		return errors.New("WebSocket con is nil")
	}
	return ws.con.WriteMessage(messageType, message)
}

// 关闭WebSocket连接
func (ws *WsServer) Close() {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	if ws.con != nil {
		ws.con.Close()
		ws.con = nil
	}
	if ws.pingTimer != nil {
		ws.pingTimer.Stop()
		ws.pingTimer = nil
	}
}
