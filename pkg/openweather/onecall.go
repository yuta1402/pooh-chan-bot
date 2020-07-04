package openweather

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// とりあえず必要な値だけ定義
type OnecallResponse struct {
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Timezone string  `json:"timezone"`

	Daily []OnecallDailyResponse `json:"daily"`
}

// とりあえず必要な値だけ定義
type OnecallDailyResponse struct {
	Dt   int64 `json:"dt"`
	Temp struct {
		Day   float64 `json:"day"`
		Min   float64 `json:"min"`
		Max   float64 `json:"max"`
		Night float64 `json:"night"`
		Eve   float64 `json:"eve"`
		Morn  float64 `json:"morn"`
	} `json:"temp"`
}

func GetOnecall(lat float64, lon float64, appid string) (OnecallResponse, error) {
	values := url.Values{}
	values.Add("lat", strconv.FormatFloat(lat, 'f', 2, 64))
	values.Add("lon", strconv.FormatFloat(lon, 'f', 2, 64))
	values.Add("appid", appid)

	// 単位を metric で固定
	values.Add("units", "metric")

	resp, err := http.Get(Endpoint + "/onecall" + "?" + values.Encode())
	if err != nil {
		return OnecallResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return OnecallResponse{}, err
	}

	res := OnecallResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return OnecallResponse{}, err
	}

	return res, err
}
