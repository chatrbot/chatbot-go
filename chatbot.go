package chatbot

import (
	"errors"
	"strconv"

	"github.com/beevik/etree"
)

//机器人实例
type ChatBot struct {
	//WebSocket和Http接口调用的token
	token string
	//服务端地址
	host string
	//机器人Http接口的包装
	bot *BotServer
	//ws连接实例
	ws *WsServer
}

//New 新建一个ChatBot实例
//@host WebSocket的服务端地址
//@token激活后机器人给的token
func New(host, token string) (*ChatBot, error) {
	w, err := newWsServer(host, token)
	if err != nil {
		return nil, err
	}
	return &ChatBot{
		token: token,
		host:  host,
		ws:    w,
		bot:   newBotServer(host, token),
	}, nil
}

//Run 连接WebSocket服务并且开始监听
func (cb *ChatBot) Run() {
	cb.ws.Listen()
}

//Use 添加处理消息的插件
func (cb *ChatBot) Use(plugin ...Plugin) {
	cb.ws.addPlugin(plugin...)
}

//SendText 发送文本形式的消息
//@toUser 接收人微信号,一般为机器人推送过来的消息发送人,即你自己
//@content 文本内容,如果有人被@需要填写对方昵称
//@atList 被@人列表,这里填写的是对方微信号
func (cb *ChatBot) SendText(toUser, content string, atList []string) error {
	//个人消息不存在@
	if !IsGroupMessage(toUser) {
		atList = nil
	}
	_, err := cb.bot.sendTextMessage(&SendTextRequest{
		ToUser:  toUser,
		AtList:  atList,
		Content: content,
	})
	return err
}

//SendPic 发送图片消息
//@toUser 接收人微信号
//@imgUrl 图片的网络地址
func (cb *ChatBot) SendPic(toUser, imgUrl string) error {
	_, err := cb.bot.sendPicMessage(&SendPicRequest{
		ToUser: toUser,
		ImgUrl: imgUrl,
	})
	return err
}

//SendVoice 发送语音
//@toUser 接收人微信号
//@url 音频文件网络地址
func (cb *ChatBot) SendVoice(toUser, url string) error {
	_, err := cb.bot.sendVoiceMessage(&SendVoiceRequest{
		ToUser:  toUser,
		SilkUrl: url,
	})
	return err
}

//SendVideo 发送视频
//@toUser 接收人微信号
//@videoUrl 视频网络地址
//@thumbUrl 封面缩略图地址
//视频发送必须有封面图,如果需要根据视频内容截取封面
//可以自行搜索ffmpeg相关的资料
func (cb *ChatBot) SendVideo(toUser, videoUrl, thumbUrl string) error {
	if thumbUrl == "" {
		return errors.New("thumbUrl is empty")
	}
	_, err := cb.bot.sendVideoMessage(&SendVideoRequest{
		ToUser:        toUser,
		VideoUrl:      videoUrl,
		VideoThumbUrl: thumbUrl,
	})
	return err
}

//SendEmoji 发送表情动图
//toUser 接收人微信号
//emojiMd5 从收到的xml中可以解析md5字段
//emojiLen 从收到的xml可以解析len字段
func (cb *ChatBot) SendEmoji(toUser, emojiMd5, emojiLen string) error {
	l, err := strconv.ParseInt(emojiLen, 10, 0)
	if err != nil {
		return err
	}
	_, err = cb.bot.sendEmojiMessage(&SendEmojiRequest{
		ToUser:        toUser,
		EmojiTotalLen: l,
		EmojiMd5:      emojiMd5,
	})
	return err
}

//DownloadPic 下载图片
func (cb *ChatBot) DownloadPic(xml string) (*DownloadImageResponse, error) {
	return cb.bot.downloadPic(&DownloadImageRequest{XML: xml})
}

//DownloadVideo 下载视频
func (cb *ChatBot) DownloadVideo(xml string) (*DownloadVideoResponse, error) {
	return cb.bot.downloadVideo(&DownloadVideoRequest{XML: xml})
}

//DownloadVoice 下载音频
func (cb *ChatBot) DownloadVoice(msgID int64, xml string) (*DownloadVoiceResponse, error) {
	return cb.bot.downloadVoice(&DownloadVoiceRequest{NewMsgId: msgID, XML: xml})
}

//DownloadEmoji 下载表情或者动态图片
//注意这个方法只能提取表情的下载地址不能用于直接发送
//发送(转发)表情需要用拿到的xml中的md5和len字段发送,可以使用ParseEmojiXML方法来获取
func (cb *ChatBot) DownloadEmoji(xml string) (string, error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return "", err
	}
	url := doc.FindElement("//emoji").SelectAttr("cdnurl").Value
	return url, nil
}

//ParseEmojiXML 解析emoji表情中的md5和len字段
func (cb *ChatBot) ParseEmojiXML(xml string) (md5, length string, err error) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		return "", "", err
	}
	e := doc.FindElement("//emoji")
	md5 = e.SelectAttr("md5").Value
	length = e.SelectAttr("len").Value
	return
}

//DelGroupMembers 删除群成员
func (cb *ChatBot) DelGroupMembers(group string, members []string) ([]string, error) {
	rsp, err := cb.bot.delGroupMembers(&DelGroupRequest{
		Group:      group,
		MemberList: members,
	})
	if err != nil {
		return nil, err
	}
	return rsp.DelMemberList, nil
}
