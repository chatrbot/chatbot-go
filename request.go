package chatbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultTimeOut = time.Second * 30         //默认超时
	urlSendText    = "/api/v1/chat/sendText"  //发送文本
	urlSendPic     = "/api/v1/chat/sendPic"   //发送图片
	urlSendEmoji   = "/api/v1/chat/sendEmoji" //发送表情
)

//BotServer 调用机器人http接口的服务
//主要用于基本的消息发送
type BotServer struct {
	host  string
	token string
}

func newBotServer(host, token string) *BotServer {
	return &BotServer{
		host:  host,
		token: token,
	}
}

//baseRequest 拼接请求
//接口都为post,token需要在url中携带
func (bs *BotServer) baseRequest(addr string, body []byte, duration time.Duration, APIRsp interface{}) error {
	u := url.URL{
		Scheme: "http",
		Host:   bs.host,
		Path:   addr,
	}
	query, _ := url.ParseQuery("token=" + bs.token)
	u.RawQuery = query.Encode()
	log.Println("request:", u.String())
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := http.Client{Timeout: duration}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("response status code err:%d", rsp.StatusCode))
	}

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	rsp.Body.Close()
	if rspBody == nil {
		return errors.New("body is nill")
	}

	return json.Unmarshal(rspBody, APIRsp)
}

func (bs *BotServer) toJson(req interface{}) []byte {
	j, _ := json.Marshal(req)
	return j
}

//sendTextMessage 发送文本消息
func (bs *BotServer) sendTextMessage(req *SendTextRequest) (*SendTextResponse, error) {
	rsp := &SendTextResponse{}
	err := bs.baseRequest(urlSendText, bs.toJson(req), defaultTimeOut, rsp)
	return rsp, err
}

//sendEmojiMessage 发送图片
func (bs *BotServer) sendPicMessage(req *SendPicRequest) (*SendPicResponse, error) {
	rsp := &SendPicResponse{}
	err := bs.baseRequest(urlSendPic, bs.toJson(req), defaultTimeOut, rsp)
	return rsp, err
}

//sendEmojiMessage 发送表情
func (bs *BotServer) sendEmojiMessage(req *SendEmojiRequest) (*SendEmojiResponse, error) {
	rsp := &SendEmojiResponse{}
	err := bs.baseRequest(urlSendEmoji, bs.toJson(req), defaultTimeOut, rsp)
	return rsp, err
}
