package poohchanbot

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

type WordCommand struct {
	AllHookWords []string
	MakeReply    func(string) linebot.SendingMessage
}

func (wc WordCommand) canHook(text string) bool {
	for _, word := range wc.AllHookWords {
		if !strings.Contains(text, word) {
			return false
		}
	}

	return true
}

func Handler() http.Handler {
	handler, err := httphandler.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)

	if err != nil {
		log.Fatal(err)
	}

	commands := []WordCommand{
		{[]string{"ぷぅちゃん", "天気"}, replyWeather},
		{[]string{"ぷーちゃん", "天気"}, replyWeather},
		{[]string{"ぷぅちゃん"}, replyRandomMessage},
		{[]string{"ぷーちゃん"}, replyRandomMessage},
	}

	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		bot, err := handler.NewClient()
		if err != nil {
			log.Print(err)
			return
		}

		for _, event := range events {
			if event.Type != linebot.EventTypeMessage {
				return
			}

			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if strings.HasPrefix(message.Text, "/print-id") {
					switch event.Source.Type {
					case linebot.EventSourceTypeUser:
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(event.Source.UserID)).Do(); err != nil {
							log.Print(err)
						}

					case linebot.EventSourceTypeGroup:
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(event.Source.GroupID)).Do(); err != nil {
							log.Print(err)
						}

					case linebot.EventSourceTypeRoom:
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(event.Source.RoomID)).Do(); err != nil {
							log.Print(err)
						}
					}
				}

				// 退出処理
				if strings.Contains(message.Text, "ぷぅちゃん") && strings.Contains(message.Text, "ハウス") {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ばいばい...")).Do(); err != nil {
						log.Print(err)
					}

					switch event.Source.Type {
					case linebot.EventSourceTypeGroup:
						if _, err := bot.LeaveGroup(event.Source.GroupID).Do(); err != nil {
							log.Print(err)
						}
					case linebot.EventSourceTypeRoom:
						if _, err := bot.LeaveRoom(event.Source.RoomID).Do(); err != nil {
							log.Print(err)
						}
					}

					break
				}

				// 返信処理
				for _, command := range commands {
					if command.canHook(message.Text) {
						replyMessage := command.MakeReply(message.Text)
						if _, err = bot.ReplyMessage(event.ReplyToken, replyMessage).Do(); err != nil {
							log.Print(err)
						}
						break
					}
				}
			}
		}
	})

	return handler
}
