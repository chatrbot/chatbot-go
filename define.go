package chatbot

import "encoding/json"

type PushMsgType int

// 自定义的消息类型,在整个消息体的最外围,即外层的msgType字段
// 用于区分是普通的聊天消息还是群内事件等
const (
	// 机器人收到的用户信息、群信息
	// 包括文本、图片、音频等多种类型
	CusMsgTypeUser PushMsgType = 10000 + iota

	// 群事件消息
	// 包括群内成员变动、机器人加入和退出群
	CusMsgTypeGroupEvent
)

// 收到的推送消息统一格式
// 根据msgType字段来区别消息类型
// data为具体的消息内容
type PushMessage struct {
	MsgType PushMsgType     `json:"msgType"`
	Data    json.RawMessage `json:"data"`
}

// 收到的转发消息具体分类
const (
	MsgTypeText   = 1  // 文本消息
	MsgTypeImg    = 3  // 图片消息
	MsgTypeVoice  = 34 // 语音消息
	MsgTypeVideo  = 43 // 视频消息
	MsgTypeEmoji  = 47 // 表情动图消息
	MsgTypeApplet = 49 // 小程序
)

const (
	RoleMember = 1 + iota
	RoleAdmin
	RoleOwner
)

// 接收到的转发消息
type UserMessage struct {
	NewMsgID       int64    `json:"newMsgId"` // 消息id,在下载语音时候会用到
	FromUser       string   `json:"fromUser"`
	AtList         []string `json:"atList"`
	CreateTime     int      `json:"createTime"`
	PushContent    string   `json:"pushContent"`
	ClientUserName string   `json:"clientUserName"`
	ToUser         string   `json:"toUser"`
	ImgBuf         string   `json:"imgBuf"`
	MsgType        int      `json:"msgType"`
	Content        string   `json:"content"`
	MsgSource      string   `json:"msgSource"`
	// 群内消息才会用到的字段
	WhoAtBot            string `json:"whoAtBot"`            // 谁@的机器人,微信昵称,方便客户端机器人反向@
	GroupMember         string `json:"groupMember"`         // 如果是群聊消息,则为分离content后的发言人微信号
	GroupMemberNickname string `json:"groupMemberNickname"` // 群内说话人的微信昵称
	GroupMemberRole     int8   `json:"groupMemberRole"`     // 用户在群内身份,1成员,2管理员,3群主
	GroupContent        string `json:"groupContent"`        // 如果是群消息,则为分离content后的群消息内容
}

func (m *UserMessage) IsAdmin() bool {
	return m.GroupMemberRole == RoleAdmin
}

func (m *UserMessage) IsGroupOwner() bool {
	return m.GroupMemberRole == RoleOwner
}

type GroupEvent int

const (
	// 机器人被邀请进群
	GroupEventInvited GroupEvent = 100000 + iota
	// 机器人被踢出群
	GroupEventKicked
	// 群内有新用户加群
	GroupEventNewMember
	// 群内有用户离开
	GroupEventMemberQuit
)

// 接收到的群内事件
type GroupBotEvent struct {
	// 事件id
	Event GroupEvent `json:"event"`
	// 事件中文提示
	EventText string `form:"eventText"`
	// 群信息
	Group GroupBase `json:"group"`
	// 变动的群成员
	Members []MemberBase `json:"members"`
}

// 群基本信息
type GroupBase struct {
	GroupUserName string `json:"groupUserName"`
	GroupNickName string `json:"groupNickName"`
	GroupHeadImg  string `json:"groupHeadImg"`
}

// 成员基本信息
type MemberBase struct {
	UserName string `json:"userName"`
	NickName string `json:"nickName"`
	HeadImg  string `json:"headImg"`
}

type (
	// 发送文本
	SendTextRequest struct {
		ToUser  string   `json:"toUser"`  // 发送对象
		AtList  []string `json:"atList"`  // 群内at的人微信号
		Content string   `json:"content"` // 发送内容 存在at时候必须有@xxx标识,xxx为对方昵称
	}
	SendTextResponse struct {
		CreateTime  int64 `json:"createTime"`  // 客户端时间
		ClientMsgId int64 `json:"clientMsgId"` // 客户端消息ID
		ServerTime  int64 `json:"serverTime"`  // 服务端时间
		MsgId       int64 `json:"msgId"`       // 服务端消息ID
		NewMsgId    int64 `json:"newMsgId"`    // 服务端消息ID
	}
	// 发送图片
	SendPicRequest struct {
		ToUser string `json:"toUser"` // 发送对象
		ImgUrl string `json:"imgUrl"` // 图片地址
	}
	SendPicResponse struct {
		ClientMsgId string `json:"clientMsgId"` // 客户端消息ID
		MsgId       int64  `json:"msgId"`       // 服务端消息ID
		NewMsgId    int64  `json:"newMsgId"`    // 服务端消息ID
	}
	// 发送表情
	SendEmojiRequest struct {
		ToUser        string `json:"toUser"`        // 发送对象
		EmojiMd5      string `json:"emojiMd5"`      // 表情md5值
		GifUrl        string `json:"gifUrl"`        // 动图地址,和md5和len互斥,不为空时候上传动图
		EmojiTotalLen int64  `json:"emojiTotalLen"` // 表情长度
	}
	SendEmojiResponse struct {
		MsgId    int64  `json:"msgId"`    // 服务端消息ID
		NewMsgId int64  `json:"newMsgId"` // 服务端消息ID
		Md5      string `json:"md5"`      // 表情md5值
		TotalLen int64  `json:"totalLen"` // 表情长度
	}
	// 发送语音
	SendVoiceRequest struct {
		ToUser  string `json:"toUser"`  // 发送对象
		SilkUrl string `json:"silkUrl"` // 语音链接
	}
	SendVoiceResponse struct {
		ClientMsgId string `json:"clientMsgId"` // 客户端消息ID
		MsgId       int64  `json:"msgId"`       // 服务端消息ID
		NewMsgId    int64  `json:"newMsgId"`    // 服务端消息ID
	}
	// 发送视频
	SendVideoRequest struct {
		ToUser        string `json:"toUser"`        // 发送对象
		VideoUrl      string `json:"videoUrl"`      // 视频地址
		VideoThumbUrl string `json:"videoThumbUrl"` // 视频缩略图地址
	}
	SendVideoResponse struct {
		ClientMsgId string `json:"clientMsgId"` // 客户端消息ID
		MsgId       int64  `json:"msgId"`       // 服务端消息ID
		NewMsgId    int64  `json:"newMsgId"`    // 服务端消息ID
	}
	// 发送小程序
	SendMiniProgramRequest struct {
		ToUser            string `json:"toUser"`            // 接收人微信号/ID
		ThumbUrl          string `json:"thumbUrl"`          //	缩略图地址
		Title             string `json:"title"`             // 标题
		Des               string `json:"des"`               // 描述
		Url               string `json:"url"`               //	地址
		SourceUserName    string `json:"sourceUserName"`    // 来源用户名
		SourceDisplayName string `json:"sourceDisplayName"` // 来源显示名
		Username          string `json:"username"`          // 用户名
		AppId             string `json:"appid"`             // 小程序 APPID
		Type              int    `json:"type"`              // 类型
		Version           int    `json:"version"`           // 版本
		IconUrl           string `json:"iconUrl"`           //  图标地址
		PagePath          string `json:"pagePath"`          // 启动页
	}
	SendMiniProgramResponse struct {
		ClientMsgId string `json:"clientMsgId"` // 客户端消息ID
		MsgId       int64  `json:"msgId"`       // 服务端消息ID
		NewMsgId    int64  `json:"newMsgId"`    // 服务端消息ID
	}
	// 下载图片
	DownloadImageRequest struct {
		XML string `json:"xml"`
	}
	DownloadImageResponse struct {
		Content []byte `json:"content"` // 下载失败后的提示
		ImgUrl  string `json:"imgUrl"`  // 图片地址
	}
	// 下载视频
	DownloadVideoRequest struct {
		XML string `json:"xml"`
	}
	DownloadVideoResponse struct {
		Content  []byte `json:"content"`  // 下载失败后的提示
		VideoUrl string `json:"videoUrl"` // 视频地址
	}
	// 下载语音
	DownloadVoiceRequest struct {
		NewMsgId int64  `json:"newMsgId"` // 服务端ID
		XML      string `json:"xml"`      // 内容xml
	}
	DownloadVoiceResponse struct {
		Content     []byte `json:"content"`     // 下载失败后的提示
		VoiceLength int64  `json:"voiceLength"` // 语音长度
		VoiceUrl    string `json:"voiceUrl"`    // 语音地址
	}
	// 删除群成员
	DelGroupRequest struct {
		Group      string   `json:"chatroom"`   // 群号
		MemberList []string `json:"memberList"` // 删除人员
	}
	DelGroupResponse struct {
		DelMemberList []string `json:"delMemberList"`
	}
)
