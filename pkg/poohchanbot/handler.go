package poohchanbot

import (
	"errors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

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

func getForecastMessage(cityID string) (string, error) {
	res, err := weatherhacks.GetForecast(cityID)
	if err != nil {
		return "", err
	}

	if len(res.Forecasts) <= 0 {
		return "", errors.New("forecast data is not available")
	}

	city := res.Location.City
	telop := res.Forecasts[0].Telop

	text := city + ":\n" +
		"    " + telop + "\n"

	tempmin := res.Forecasts[0].Temperature.Min.Celsius
	tempmax := res.Forecasts[0].Temperature.Max.Celsius

	if tempmin != "" && tempmax != "" {
		text += "    " + "最低 " + tempmin + "℃" + " / " + "最高 " + tempmax + "℃"
	}

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

	text := "今日の天気だよ♪\n" + "\n" + text0 + "\n" + text1
	return text, nil
}

func replyWeather(string) linebot.SendingMessage {
	text, err := forecastMessage()
	if err != nil {
		log.Print(err)
	}

	return linebot.NewTextMessage(text)
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
