package main

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/chatrbot/chatbot-go"
)

var (
	token = flag.String("token", "", "ChatBot Token")
	host  = flag.String("host", "118.25.84.114:18881", "WebSocket Server Host")
)

// 发送小程序示例
func main() {
	flag.Parse()

	bot, err := chatbot.New(*host, *token)
	if err != nil {
		log.Fatalln("连接服务器失败:", err)
	}

	repeat := NewTranscoder(bot)
	bot.Use(repeat)

	bot.Run()
}

var _ chatbot.Plugin = new(MiniProgramDemo)

type MiniProgramDemo struct {
	bot  *chatbot.ChatBot
	name string
}

func NewTranscoder(bot *chatbot.ChatBot) *MiniProgramDemo {
	return &MiniProgramDemo{
		bot:  bot,
		name: "MiniProgramDemo",
	}
}

func (ts *MiniProgramDemo) Name() string {
	return "MiniProgramDemo"
}

// 发送"小程序"三个字给机器人会返回一个肯德基的小程序
func (ts *MiniProgramDemo) Do(msg *chatbot.PushMessage) error {
	if msg.MsgType == chatbot.CusMsgTypeUser {
		// 获取接收人等基本信息
		message := &chatbot.UserMessage{}
		_ = json.Unmarshal(msg.Data, message)

		if message.MsgType == chatbot.MsgTypeText && message.Content == "小程序" {
			// 小程序相关字段需要从收到的xml中解析
			// 可以用机器人接收一次小程序,观察下收到的xml结构
			// 其中一些封面图片等非关键字段不一定需要一一对应,可以改成自己想要的
			return ts.bot.SendMiniProgram(&chatbot.SendMiniProgramRequest{
				ToUser:            message.FromUser,
				ThumbUrl:          "http://mmbiz.qpic.cn/mmbiz_png/SE9ICmPPKWiaibdENZqwnjeIWiblOvnX4QFZMr2PJ704lOyphLBicqjwYbt9Rsiak2mYM8UBtTX91XgMg3lqs98DMMA/640?wx_fmt=png&wxfrom=200",
				Title:             "肯德基自助点餐",
				Des:               "肯德基+",
				Url:               "https://mp.weixin.qq.com/mp/waerrpage?appid=wx23dde3ba32269caa&type=upgrade&upgradetype=3#wechat_redirect",
				SourceUserName:    "gh_50338e5b8c9d",
				SourceDisplayName: "肯德基+",
				Username:          "gh_50338e5b8c9d",
				AppId:             "wx23dde3ba32269caa",
				Type:              2,
				Version:           92,
				IconUrl:           "http://mmbiz.qpic.cn/mmbiz_png/SE9ICmPPKWiaibdENZqwnjeIWiblOvnX4QFZMr2PJ704lOyphLBicqjwYbt9Rsiak2mYM8UBtTX91XgMg3lqs98DMMA/640?wx_fmt=png&wxfrom=200",
				PagePath:          "pages/home/home.html",
			})
		}
	}
	return nil
}
