package chatbot

//Plugin 机器人插件
type Plugin interface {
	Name() string
	Do(msg *ReceiveMessage) error
}
