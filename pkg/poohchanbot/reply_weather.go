package poohchanbot

import (
	"errors"
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/yuta1402/pooh-chan-bot/pkg/weatherhacks"
)

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
