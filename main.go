package main

import (
	"fmt"
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

func replyPoohChan(string) linebot.SendingMessage {
	return linebot.NewTextMessage("ぷぅちゃん！")
}

func main() {
	handler, err := httphandler.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)

	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	commands := []WordCommand{
		{[]string{"ぷぅちゃん"}, replyPoohChan},
		{[]string{"ぷーちゃん"}, replyPoohChan},
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
				for _, command := range commands {
					if command.canHook(message.Text) {
						replyMessage := command.MakeReply(message.Text)
						if _, err = bot.ReplyMessage(event.ReplyToken, replyMessage).Do(); err != nil {
							log.Print(err)
						}
					}
				}
			}
		}
	})

	http.Handle("/callback", handler)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")
	})

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
