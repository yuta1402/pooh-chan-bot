package poohchanbot

import (
	"math/rand"

	"github.com/line/line-bot-sdk-go/linebot"
)

func replyRandomMessage(string) linebot.SendingMessage {
	messages := []string{
		"ぷぅちゃん！",
		"ぷぅちゃん♪",
		"ぷぅちゃん♡",
		"♡",
		"♪",
		"ぷぅ？",
		"ぷっぷぷ〜ぷっぷぷ〜ぷっぷぷっぷぷぅ〜♪",
		"シャーッ！！！",
		"ギャーギャー！ギャーギャー！",
	}

	r := rand.Intn(len(messages))
	text := messages[r]

	return linebot.NewTextMessage(text)
}
