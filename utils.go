package chatbot

import (
	"strings"
)

//IsGroupMessage 是否为群聊消息
func IsGroupMessage(userName string) bool {
	return strings.HasSuffix(userName, "@chatroom")
}

//IsBotBeenAt 机器人是否被@了
func IsBotBeenAt(msg *UserMessage) bool {
	for _, u := range msg.AtList {
		if msg.ClientUserName == u {
			return true
		}
	}
	return false
}

//SplitAtContent 分割包含@的消息内容
//例如 @小明 你吃饭了么
//其中消息体中间会有个"空格"，这个可能是个特殊字符，也可能是个真的空格
func SplitAtContent(msgContent string) string {
	content := msgContent

	contents := strings.SplitN(msgContent, "\u2005", 2)
	if len(contents) != 2 {
		contents = strings.SplitN(msgContent, " ", 2)
	}
	if len(contents) == 2 {
		content = strings.TrimSpace(contents[1])
	}
	return content
}
