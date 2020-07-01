package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

type WordCommand struct {
	AllHookWords []string
	MakeReply    func(string) linebot.SendingMessage
}

type WeatherHack struct {
	Title       string   `xml:"channel>title"`
	Description []string `xml:"channel>item>description"`
}

type ForecastResponse struct {
	Forecasts []struct {
		DataLabel string `json:"data_label"`
		Telop     string `json:"telop"`
		Date      string `json:"date"`
	} `json:"forecasts"`
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

func replyWeather(string) linebot.SendingMessage {
	values := url.Values{}
	values.Add("city", "400040")
	resp, err := http.Get("http://weather.livedoor.com/forecast/webservice/json/v1" + "?" + values.Encode())
	if err != nil {
		log.Print(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return nil
	}

	// fmt.Println(string(body))

	res := ForecastResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		log.Print(err)
	}

	log.Print(res)
	return nil
	// return linebot.NewTextMessage(string(body))
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
		{[]string{"ぷぅちゃん", "天気"}, replyWeather},
		{[]string{"ぷーちゃん", "天気"}, replyWeather},
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
