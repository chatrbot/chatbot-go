package chatbot

import "strings"

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
