package main

import (
	"chatbot-go"
	"encoding/json"
	"flag"
	"fmt"
	"log"
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

	repeat := NewRepeatPlugin(bot)
	bot.Use(repeat)

	bot.Run()
}

//复读机插件
//会重复群内用户的发送内容
//用于展示不同消息的收发
type RepeatPlugin struct {
	bot  *chatbot.ChatBot
	name string
}

var _ chatbot.Plugin = new(RepeatPlugin)

func NewRepeatPlugin(bot *chatbot.ChatBot) *RepeatPlugin {
	return &RepeatPlugin{name: "RepeatPlugin", bot: bot}
}

func (p *RepeatPlugin) Name() string {
	return p.name
}

func (p *RepeatPlugin) Do(msg *chatbot.PushMessage) error {
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

//handleMessage 处理机器人收到的聊天消息
//其中包含了私聊消息和群消息 需要自己判断
func (p *RepeatPlugin) handleMessage(msg *chatbot.UserMessage) error {
	if chatbot.IsGroupMessage(msg.FromUser) {
		if chatbot.IsBotBeenAt(msg) {
			if err := p.bot.SendText(msg.FromUser, fmt.Sprintf("@%s 叫我干嘛", msg.WhoAtBot), []string{msg.GroupMember}); err != nil {
				log.Println("发送@回复失败", err)
			}
		} else {
			switch msg.MsgType {
			case chatbot.MsgTypeText:
				return p.bot.SendText(msg.FromUser, msg.PushContent, nil)
			case chatbot.MsgTypeImg:
				if rsp, err := p.bot.DownloadPic(msg.GroupContent); err != nil {
					return fmt.Errorf("下载图片失败:%w", err)
				} else {
					log.Println("图片地址", rsp.ImgUrl)
					if err := p.bot.SendPic(msg.FromUser, rsp.ImgUrl); err != nil {
						return fmt.Errorf("发送图片消息失败:%w", err)
					}
				}
			case chatbot.MsgTypeVoice:
				if rsp, err := p.bot.DownloadVoice(msg.NewMsgID, msg.GroupContent); err != nil {
					return fmt.Errorf("下载语音失败:%w", err)
				} else {
					log.Println("语音地址", rsp.VoiceUrl)
					if err := p.bot.SendVoice(msg.FromUser, rsp.VoiceUrl); err != nil {
						return fmt.Errorf("发送图片消息失败:%w", err)
					}
				}
			case chatbot.MsgTypeVideo:
				if rsp, err := p.bot.DownloadVideo(msg.GroupContent); err != nil {
					return fmt.Errorf("下载视频失败:%w", err)
				} else {
					log.Println("视频地址", rsp.VideoUrl)
					if err := p.bot.SendVideo(msg.FromUser, rsp.VideoUrl, "https://ss1.bdstatic.com/70cFvXSh_Q1YnxGkpoWK1HF6hhy/it/u=1472019148,226459533&fm=26&gp=0.jpg"); err != nil {
						return fmt.Errorf("发送图片消息失败:%w", err)
					}
				}
			case chatbot.MsgTypeEmoji:
				if md5, l, err := p.bot.ParseEmojiXML(msg.GroupContent); err != nil {
					return fmt.Errorf("解析表情失败:%w", err)
				} else {
					if err := p.bot.SendEmoji(msg.FromUser, md5, l); err != nil {
						return fmt.Errorf("发送Emoji图片消息失败:%w", err)
					}
				}
			default:
				log.Println("未知消息类型:", msg.MsgType)
			}
		}
	}
	return nil
}

//handleGroupEvent 处理群内事件
func (p *RepeatPlugin) handleGroupEvent(msg *chatbot.GroupBotEvent) error {
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
