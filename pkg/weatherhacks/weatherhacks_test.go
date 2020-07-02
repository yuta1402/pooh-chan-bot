package weatherhacks

import (
	"log"
	"testing"
)

func TestGetForecast(t *testing.T) {
	res, err := GetForecast("130010")
	if err != nil {
		t.Fatal(err)
	}

	if res.Location.City != "東京" {
		t.Fatalf("city: %s, expected city: %s", res.Location.City, "東京")
	}

	if res.Location.Area != "関東" {
		t.Fatalf("area: %s, expected area: %s", res.Location.Area, "関東")
	}

	if res.Location.Prefecture != "東京都" {
		t.Fatalf("prefecture: %s, expected prefecture: %s", res.Location.Prefecture, "東京都")
	}

	// dateLabel := res.Forecasts[0].DateLabel
	telop := res.Forecasts[0].Telop
	city := res.Location.City

	text := city + ": " + telop
	// text := dateLabel + "の" + city + "の天気は「" + telop + "」だよ♪"
	log.Print(text)
}
