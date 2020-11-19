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
	lock    sync.Mutex
	con     *websocket.Conn
	plugins []Plugin
}

//新建监听器
func newWsServer(host, token string) (*WsServer, error) {
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
	return &WsServer{con: c, plugins: make([]Plugin, 0, 10)}, nil
}

func (ws *WsServer) Listen() error {
	for {
		_, msg, err := ws.con.ReadMessage()
		if err != nil {
			log.Println("read msg err:", err)
			return err
		}
		log.Println("收到消息", string(msg))

		if len(ws.plugins) > 0 {
			var rec ReceiveMessage
			_ = json.Unmarshal(msg, &rec)

			for _, p := range ws.plugins {
				if err := p.Do(&rec); err != nil {
					log.Printf("%s handle error:%s \n", p.Name(), err)
				}
			}
		}
	}
}

func (ws *WsServer) addPlugin(plugin ...Plugin) {
	ws.plugins = append(ws.plugins, plugin...)
}

//writeMessage gorilla的websocket库默认发送消息会有并发问题
func (ws *WsServer) writeMessage(message string) error {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	return ws.con.WriteMessage(websocket.TextMessage, []byte(message))
}

//关闭WebSocket连接
func (ws *WsServer) Close() error {
	return ws.con.Close()
}
