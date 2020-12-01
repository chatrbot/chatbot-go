package chatbot

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestOwnThink(t *testing.T) {
	appID := ""
	ques := "姚明"
	q := url.QueryEscape(ques)
	rsp, err := http.DefaultClient.Get(fmt.Sprintf("https://api.ownthink.com/bot?appid=%s&spoken=%s", appID, q))
	if err != nil {
		t.Error(err)
		return
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Error("read body err", err)
		return
	}
	rsp.Body.Close()

	t.Log(string(body))
}

func TestSplitAtContent(t *testing.T) {
	text := "@test  讲个笑话"
	content := SplitAtContent(text)
	t.Log(content)
}
