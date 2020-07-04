package poohchanbot

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

func replyPoohChan(string) linebot.SendingMessage {
	return linebot.NewTextMessage("ぷぅちゃん！")
}
