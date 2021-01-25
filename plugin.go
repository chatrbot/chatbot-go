package chatbot

// Plugin 机器人插件接口
type Plugin interface {
	Name() string
	Do(msg *PushMessage) error
}
