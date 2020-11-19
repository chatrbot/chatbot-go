package main

import (
	"chatbot-go"
	"flag"
	"log"
	"strings"
)

// token 修改为从机器人处获得的token
// 在这里获取token  https://github.com/chatrbot/chatbot#faq
// host websocket的服务端地址
var (
	token = flag.String("token", "", "chatbot token")
	host  = flag.String("host", "118.25.84.114:18881", "websocket server host")
)

func init() {
	flag.Parse()
}

func main() {
	bot, err := chatbot.New(*host, *token)
	if err != nil {
		log.Fatalln("连接服务器失败:", err)
	}

	pic := NewPicPlugin(bot)
	bot.Use(pic)

	bot.Run()
}

type HelloPlugin struct {
	bot  *chatbot.ChatBot
	name string
}

var _ chatbot.Plugin = new(HelloPlugin)

func NewPicPlugin(bot *chatbot.ChatBot) *HelloPlugin {
	return &HelloPlugin{name: "HelloPlugin", bot: bot}
}

func (p *HelloPlugin) Name() string {
	return p.name
}

func (p *HelloPlugin) Do(msg *chatbot.ReceiveMessage) error {
	if strings.ToLower(msg.Content) == "hello" {
		return p.bot.SendText(msg.FromUser, "world")
	}
	return nil
}
