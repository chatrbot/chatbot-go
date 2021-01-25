package chatbot

import (
	"testing"
)

var (
	bot      *ChatBot
	testUser = ""
)

func TestMain(t *testing.M) {
	host := "127.0.0.1:18083"
	// 测试服的token,你拿去用是无效的
	token := "63c5a2edf6ff4418b59419f09cba35a4"
	bot = &ChatBot{
		token: token,
		host:  host,
		bot:   newBotServer(host, token),
	}
	t.Run()
}

func TestChatBot_SendText(t *testing.T) {
	err := bot.SendText(testUser, "test", nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("ok")
}
