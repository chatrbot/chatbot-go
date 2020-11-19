package chatbot

//ReceiveMessage 接收到的推送消息
type ReceiveMessage struct {
	MsgSeq         int      `json:"msgSeq"`
	ClientID       string   `json:"clientId"`
	ReportMsgType  int      `json:"reportMsgType"`
	FromUser       string   `json:"fromUser"`
	AtList         []string `json:"atList"`
	CreateTime     int      `json:"createTime"`
	NewMsgID       int64    `json:"newMsgId"`
	PushContent    string   `json:"pushContent"`
	ClientUserName string   `json:"clientUserName"`
	ToUser         string   `json:"toUser"`
	MsgID          int      `json:"msgId"`
	ImgBuf         string   `json:"imgBuf"`
	MsgType        int      `json:"msgType"`
	Content        string   `json:"content"`
	MsgSource      string   `json:"msgSource"`
}

type (
	//发送文本
	SendTextRequest struct {
		ToUser  string   `json:"toUser"`  //发给谁
		AtList  []string `json:"atList"`  //群内at谁
		Content string   `json:"content"` //发送内容 存在at时候必须有@xxx标识
	}
	SendTextResponse struct {
		CreateTime  int64 `json:"createTime"`  //客户端时间
		ClientMsgId int64 `json:"clientMsgId"` //客户端消息ID
		ServerTime  int64 `json:"serverTime"`  //服务端时间
		MsgId       int64 `json:"msgId"`       //服务端消息ID
		NewMsgId    int64 `json:"newMsgId"`    //服务端消息ID
	}
	//发送图片
	SendPicRequest struct {
		ToUser string `json:"toUser"` //发给谁
		ImgUrl string `json:"imgUrl"` //图片地址
	}
	SendPicResponse struct {
		ClientMsgId string `json:"clientMsgId"` //客户端消息ID
		MsgId       int64  `json:"msgId"`       //服务端消息ID
		NewMsgId    int64  `json:"newMsgId"`    //服务端消息ID
	}
	//发送表情
	SendEmojiRequest struct {
		ToUser        string `json:"toUser"`        //发给谁
		EmojiMd5      string `json:"emojiMd5"`      //表情md5值
		GifUrl        string `json:"gifUrl"`        //动图地址,和md5和len互斥,不为空时候上传动图
		EmojiTotalLen int64  `json:"emojiTotalLen"` //表情长度
	}
	SendEmojiResponse struct {
		MsgId    int64  `json:"msgId"`    //服务端消息ID
		NewMsgId int64  `json:"newMsgId"` //服务端消息ID
		Md5      string `json:"md5"`      //表情md5值
		TotalLen int64  `json:"totalLen"` //表情长度
	}
)
