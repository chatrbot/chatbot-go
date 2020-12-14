/**
群管理员机器人
*/

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/chatrbot/chatbot-go"
)

// token 修改为从机器人处获得的token
// 在这里获取token  https://github.com/chatrbot/chatbot#faq
// host WebSocket的服务端地址
var (
	token = flag.String("token", "", "ChatBot Token")
	host  = flag.String("host", "118.25.84.114:18881", "WebSocket Server Host")
)

func init() {
	flag.Parse()
}

func main() {
	bot, err := chatbot.New(*host, *token)
	if err != nil {
		log.Fatalln("连接服务器失败:", err)
	}

	manager := NewGroupManagerPlugin(bot)
	bot.Use(manager)

	bot.Run()
}

type GroupManagerPlugin struct {
	bot  *chatbot.ChatBot
	name string
}

func NewGroupManagerPlugin(bot *chatbot.ChatBot) *GroupManagerPlugin {
	return &GroupManagerPlugin{name: "groupManager", bot: bot}
}

func (p *GroupManagerPlugin) Name() string {
	return p.name
}

func (p *GroupManagerPlugin) Do(msg *chatbot.PushMessage) error {
	switch msg.MsgType {
	case chatbot.CusMsgTypeUser:
		m := &chatbot.UserMessage{}
		if err := json.Unmarshal(msg.Data, m); err != nil {
			return err
		}
		return p.handleMessage(m)
	case chatbot.CusMsgTypeGroupEvent:
		e := &chatbot.GroupBotEvent{}
		if err := json.Unmarshal(msg.Data, e); err != nil {
			return err
		}
		return p.handleGroupEvent(e)
	default:
		log.Println("消息类型错误")
	}
	return nil
}

// handleMessage 处理机器人收到的聊天消息
// 使用@的方法可以做到快速踢人:@somebody 踢
// 命令机器人踢出群成员,注意机器人必须为群管理员身份
func (p *GroupManagerPlugin) handleMessage(msg *chatbot.UserMessage) error {
	kickKeyword := "踢123"
	if chatbot.IsGroupMessage(msg.FromUser) &&
		msg.MsgType == chatbot.MsgTypeText &&
		len(msg.AtList) > 0 {
		//判断身份这条消息发送人的身份
		if !msg.IsAdmin() && !msg.IsGroupOwner() {
			if err := p.bot.SendText(msg.FromUser, "你不是管理员不能命令我", []string{msg.GroupMember}); err != nil {
				log.Println("发送消息失败", err)
				return err
			}
			return errors.New("不是管理员身份,不能进行操作")
		}

		content := chatbot.SplitAtContent(msg.GroupContent)
		if content == kickKeyword {
			_, err := p.bot.DelGroupMembers(msg.FromUser, msg.AtList)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("未处理的消息类型")
}

//handleGroupEvent 处理群内事件
func (p *GroupManagerPlugin) handleGroupEvent(msg *chatbot.GroupBotEvent) error {
	switch msg.Event {
	case chatbot.GroupEventInvited:
		return p.bot.SendText(msg.Group.GroupUserName, "大家好我是机器人", nil)
	case chatbot.GroupEventKicked:
		log.Println("机器人被踢出群了!", msg.Group.GroupNickName)
	case chatbot.GroupEventNewMember:
		for _, m := range msg.Members {
			return p.bot.SendText(msg.Group.GroupUserName, fmt.Sprintf("欢迎新成员:%s", m.NickName), nil)
		}
	case chatbot.GroupEventMemberQuit:
		for _, m := range msg.Members {
			return p.bot.SendText(msg.Group.GroupUserName, fmt.Sprintf("有人离开了:%s", m.NickName), nil)
		}
	default:
		log.Println("未知事件")
	}
	return nil
}
