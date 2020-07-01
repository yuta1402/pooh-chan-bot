package weatherhacks

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	Endpoint = "http://weather.livedoor.com/forecast/webservice/json/v1"
)

type ForecastResponse struct {
	PinpointLocations []struct {
		Link string `json:"link"`
		Name string `json:"name"`
	} `json:"pinpointLocations"`

	Link string `json:"link"`

	Forecasts []struct {
		DateLabel string `json:"dateLabel"`
		Telop     string `json:"telop"`
		Date      string `json:"date"`
	} `json:"forecasts"`

	Location struct {
		City       string `json:"city"`
		Area       string `json:"area"`
		Prefecture string `json:"prefecture"`
	} `json:"location"`

	PublicTime string `json:"publicTime"`

	Title string `json:"title"`

	Description struct {
		Text       string `json:"text"`
		PublicTime string `json:"publicTime"`
	} `json:"description"`
}

func GetForecast(cityID string) (ForecastResponse, error) {
	values := url.Values{}
	values.Add("city", cityID)
	resp, err := http.Get(Endpoint + "?" + values.Encode())
	if err != nil {
		return ForecastResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ForecastResponse{}, err
	}

	res := ForecastResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return ForecastResponse{}, err
	}

	return res, nil
}
