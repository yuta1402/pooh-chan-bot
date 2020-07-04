package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"github.com/yuta1402/pooh-chan-bot/pkg/weatherhacks"
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

func replyRandomMessage(string) linebot.SendingMessage {
	messages := []string{
		"ぷぅちゃん！",
		"ぷぅ？",
		"ぷっぷぷ〜ぷっぷぷ〜ぷっぷぷっぷぷぅ〜♪",
		"シャーッ！！！",
		"ギャーギャーギャーギャー...",
	}

	r := rand.Intn(len(messages))
	text := messages[r]

	return linebot.NewTextMessage(text)
}

func getForecastMessage(cityID string) (string, error) {
	res, err := weatherhacks.GetForecast(cityID)
	if err != nil {
		return "", err
	}

	if len(res.Forecasts) <= 0 {
		return "", errors.New("forecast data is not available")
	}

	telop := res.Forecasts[0].Telop
	city := res.Location.City

	text := city + ": " + telop

	return text, nil
}

func forecastMessage() (string, error) {
	text0, err := getForecastMessage("130010")
	if err != nil {
		return "", err
	}

	text1, err := getForecastMessage("230010")
	if err != nil {
		return "", err
	}

	text := "今日の天気は、\n" + "  " + text0 + "\n" + "  " + text1 + "\n" + "だよ♪"
	return text, nil
}

func replyWeather(string) linebot.SendingMessage {
	text, err := forecastMessage()
	if err != nil {
		log.Print(err)
	}

	return linebot.NewTextMessage(text)
}

func main() {
	rand.Seed(time.Now().Unix())

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

	http.Handle("/api/v1/webhook", handler)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "pong")
	})

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
