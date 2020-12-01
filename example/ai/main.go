/**
这是一个接入思知API让机器人成为AI机器人的例子
里面需要的思知token需要自行注册获取 https://www.ownthink.com/
*/
package main

import (
	chatbot "chatbot-go"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

// token 修改为从机器人处获得的token
// 在这里获取token  https://github.com/chatrbot/chatbot#faq
// host WebSocket的服务端地址
var (
	AIToken = flag.String("ai", "", "OwnThink Token")
	token   = flag.String("token", "", "ChatBot Token")
	host    = flag.String("host", "118.25.84.114:18881", "WebSocket Server Host")
)

func init() {
	flag.Parse()
}

func main() {
	bot, err := chatbot.New(*host, *token)
	if err != nil {
		log.Fatalln("连接服务器失败:", err)
	}

	repeat := NewRepeatPlugin(bot)
	bot.Use(repeat)

	bot.Run()
}

//AI插件,接入AI API,可以和用户做智能对话
type AIPlugin struct {
	bot  *chatbot.ChatBot
	name string
}

var _ chatbot.Plugin = new(AIPlugin)

func NewRepeatPlugin(bot *chatbot.ChatBot) *AIPlugin {
	return &AIPlugin{name: "AIPlugin", bot: bot}
}

func (p *AIPlugin) Name() string {
	return p.name
}

func (p *AIPlugin) Do(msg *chatbot.PushMessage) error {
	if msg.MsgType == chatbot.CusMsgTypeUser {
		message := &chatbot.UserMessage{}
		if err := json.Unmarshal(msg.Data, message); err != nil {
			return err
		}
		//如果是群内消息
		if chatbot.IsGroupMessage(message.FromUser) {
			//如果是机器人被@了
			if chatbot.IsBotBeenAt(message) {
				keyword := chatbot.SplitAtContent(message.GroupContent)
				reply, err := OwnThinkAPI(keyword)
				if err != nil {
					return err
				}
				if err := p.bot.SendText(
					message.FromUser,
					fmt.Sprintf("@%s %s", message.WhoAtBot, reply),
					[]string{message.GroupMember},
				); err != nil {
					return fmt.Errorf("发送群内@回复消息失败:%w", err)
				}
			}
		} else {
			reply, err := OwnThinkAPI(message.Content)
			if err != nil {
				return err
			}
			if err := p.bot.SendText(
				message.FromUser,
				reply,
				nil,
			); err != nil {
				return fmt.Errorf("发送回复消息失败:%w", err)
			}
		}
	}
	return nil
}

//OwnThinkBot 思知AI接口
func OwnThinkAPI(content string) (string, error) {
	log.Println("Receive Content", content)
	c := url.QueryEscape(content)
	rsp, err := http.DefaultClient.Get(fmt.Sprintf("https://api.ownthink.com/bot?appid=%s&userid=user&spoken=%s", *AIToken, c))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", fmt.Errorf("read body error:%e", err)
	}
	_ = rsp.Body.Close()

	reply := gjson.ParseBytes(body).Get("data.info.text").String()
	log.Println("收到回复:", reply)
	return reply, nil
}
