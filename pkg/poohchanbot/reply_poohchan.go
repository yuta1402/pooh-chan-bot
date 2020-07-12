package poohchanbot

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

func replyPoohChan(string) *sendingMessageQueue {
	q := newSendingMessageQueue()
	q.enque(linebot.NewTextMessage("ぷぅちゃん！"))
	return q
}
