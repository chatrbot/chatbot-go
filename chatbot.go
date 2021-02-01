package chatbot

import (
	"errors"
	"strconv"

	"github.com/beevik/etree"
)

// 机器人实例
type ChatBot struct {
	// WebSocket和Http接口调用的token
	token string
	// 服务端地址
	host string
	// 机器人Http接口的包装
	bot *BotServer
	// ws连接实例
	ws *WsServer
}

// New 新建一个ChatBot实例
// @host WebSocket的服务端地址
// @token激活后机器人给的token
func New(host, token string) (*ChatBot, error) {
	ws, err := newWSClient(host, token)
	if err != nil {
		return nil, err
	}
	return &ChatBot{
		token: token,
		host:  host,
		ws:    ws,
		bot:   newBotServer(host, token),
	}, nil
}

func (bot *ChatBot) Close() {
	bot.ws.Close()
}

// Run 连接WebSocket服务并且开始监听
func (bot *ChatBot) Run() {
	bot.ws.ReceiveCallbackMessage()
}

// Use 添加处理消息的插件
func (bot *ChatBot) Use(plugin ...Plugin) {
	bot.ws.addPlugin(plugin...)
}

// SendText 发送文本形式的消息
// @toUser 接收人微信号,一般为机器人推送过来的消息发送人,即你自己
// @content 文本内容,如果有人被@需要填写对方昵称
// @atList 被@人列表,这里填写的是对方微信号
func (bot *ChatBot) SendText(toUser, content string, atList []string) error {
	// 个人消息不存在@
	if !IsGroupMessage(toUser) {
		atList = nil
	}
	_, err := bot.bot.sendTextMessage(&SendTextRequest{
		ToUser:  toUser,
		AtList:  atList,
		Content: content,
	})
	return err
}

// SendPic 发送图片消息
// @toUser 接收人微信号
// @imgUrl 图片的网络地址
func (bot *ChatBot) SendPic(toUser, imgUrl string) error {
	_, err := bot.bot.sendPicMessage(&SendPicRequest{
		ToUser: toUser,
		ImgUrl: imgUrl,
	})
	return err
}

// SendVoice 发送语音
// @toUser 接收人微信号
// @url 音频文件网络地址
func (bot *ChatBot) SendVoice(toUser, url string) error {
	_, err := bot.bot.sendVoiceMessage(&SendVoiceRequest{
		ToUser:  toUser,
		SilkUrl: url,
	})
	return err
}

// SendVideo 发送视频
// @toUser 接收人微信号
// @videoUrl 视频网络地址
// @thumbUrl 封面缩略图地址
// 视频发送必须有封面图,如果需要根据视频内容截取封面
// 可以自行搜索ffmpeg相关的资料
func (bot *ChatBot) SendVideo(toUser, videoUrl, thumbUrl string) error {
	if thumbUrl == "" {
		return errors.New("thumbUrl is empty")
	}
	_, err := bot.bot.sendVideoMessage(&SendVideoRequest{
		ToUser:        toUser,
		VideoUrl:      videoUrl,
		VideoThumbUrl: thumbUrl,
	})
	return err
}

// SendEmoji 发送表情动图
// toUser 接收人微信号
// emojiMd5 从收到的xml中可以解析md5字段
// emojiLen 从收到的xml可以解析len字段
func (bot *ChatBot) SendEmoji(toUser, emojiMd5, emojiLen string) error {
	l, err := strconv.ParseInt(emojiLen, 10, 0)
	if err != nil {
		return err
	}
	_, err = bot.bot.sendEmojiMessage(&SendEmojiRequest{
		ToUser:        toUser,
		EmojiTotalLen: l,
		EmojiMd5:      emojiMd5,
	})
	return err
}

// SendMiniProgram 发送小程序
// toUser	接收人微信号/ID
// thumbUrl	缩略图地址
// title 标题
// des 描述
// url 地址
// sourceUserName 来源用户名
// sourceDisplayName 来源显示名
// username	用户名
// appId 小程序AppId
// type 类型
// version 版本
// iconUrl 图标地址
// pagePath 启动页
func (bot *ChatBot) SendMiniProgram(req *SendMiniProgramRequest) error {
	_, err := bot.bot.sendMiniProgramMessage(req)
	return err
}

// DownloadPic 下载图片
func (bot *ChatBot) DownloadPic(xml string) (*DownloadImageResponse, error) {
	return bot.bot.downloadPic(&DownloadImageRequest{XML: xml})
}

// DownloadVideo 下载视频
func (bot *ChatBot) DownloadVideo(xml string) (*DownloadVideoResponse, error) {
	return bot.bot.downloadVideo(&DownloadVideoRequest{XML: xml})
}

// DownloadVoice 下载音频
func (bot *ChatBot) DownloadVoice(msgID int64, xml string) (*DownloadVoiceResponse, error) {
	return bot.bot.downloadVoice(&DownloadVoiceRequest{NewMsgId: msgID, XML: xml})
}

// DownloadEmoji 下载表情或者动态图片
// 注意这个方法只能提取表情的下载地址不能用于直接发送
// 发送(转发)表情需要用拿到的xml中的md5和len字段发送,可以使用ParseEmojiXML方法来获取
func (bot *ChatBot) DownloadEmoji(xml string) (string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return "", err
	}
	url := doc.FindElement("//emoji").SelectAttr("cdnurl").Value
	return url, nil
}

// ParseEmojiXML 解析emoji表情中的md5和len字段
func (bot *ChatBot) ParseEmojiXML(xml string) (md5, length string, err error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return "", "", err
	}
	e := doc.FindElement("//emoji")
	md5 = e.SelectAttr("md5").Value
	length = e.SelectAttr("len").Value
	return
}

// DelGroupMembers 删除群成员
func (bot *ChatBot) DelGroupMembers(group string, members []string) ([]string, error) {
	rsp, err := bot.bot.delGroupMembers(&DelGroupRequest{
		Group:      group,
		MemberList: members,
	})
	if err != nil {
		return nil, err
	}
	return rsp.DelMemberList, nil
}
