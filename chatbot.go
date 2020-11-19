package chatbot

import "log"

type ChatBot struct {
	token string
	host  string
	bot   *BotServer
	ws    *WsServer
}

//New 新建一个chatbot实例
//host websocket的服务端地址
//token激活后机器人给的token
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

//Run 连接websocket服务并且开始监听
func (cb *ChatBot) Run() {
	log.Fatalln(cb.ws.Listen())
}

//Use 添加处理消息的插件
func (cb *ChatBot) Use(plugin ...Plugin) {
	cb.ws.addPlugin(plugin...)
}

//SendText 发送文本形式的消息
//toUser 接收人微信号,一般为机器人推送过来的消息发送人,即你自己
//content 文本内容
func (cb *ChatBot) SendText(toUser, content string) error {
	_, err := cb.bot.sendTextMessage(&SendTextRequest{
		ToUser:  toUser,
		AtList:  nil,
		Content: content,
	})
	return err
}

//SendPic 发送图片消息
//toUser 接收人微信号
//imgUrl 图片的网络地址
func (cb *ChatBot) SendPic(toUser, imgUrl string) error {
	_, err := cb.bot.sendPicMessage(&SendPicRequest{
		ToUser: toUser,
		ImgUrl: imgUrl,
	})
	return err
}

//SendEmoji 发送表情动图
//toUser 接收人微信号
//imgUrl 表情地址
//emojiMd5 如果是接收到的表情会有这个参数,没有可不传
//emojiLen 如果是接收到的表情会有这个参数,没有可不传
func (cb *ChatBot) SendEmoji(toUser, imgUrl, emojiMd5 string, emojiLen int64) error {
	_, err := cb.bot.sendEmojiMessage(&SendEmojiRequest{
		ToUser:        toUser,
		EmojiTotalLen: emojiLen,
		EmojiMd5:      emojiMd5,
		GifUrl:        imgUrl,
	})
	return err
}
