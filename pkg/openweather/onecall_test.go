package openweather

import (
	"os"
	"testing"
)

func TestGetOnecall(t *testing.T) {
	lat := 35.69
	lon := 139.69
	timezone := "Asia/Tokyo"

	res, err := GetOnecall(lat, lon, os.Getenv("OPENWEATHER_APIKEY"))
	if err != nil {
		t.Fatal(err)
	}

	if res.Lat != lat {
		t.Fatalf("res.Lat: %f, expected res.Lat: %f", res.Lat, lat)
	}

	if res.Lon != lon {
		t.Fatalf("res.Lon: %f, expected res.Lon: %f", res.Lon, lon)
	}

	if res.Timezone != timezone {
		t.Fatalf("res.Timezone: %s, expected: %s", res.Timezone, timezone)
	}
}
