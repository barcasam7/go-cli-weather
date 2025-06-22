package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

type weatherResponse struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Date string `json:"date"`
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	err := godotenv.Load() // load env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	q := "Bridlington"

	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	fmt.Print("How many days of weather forecast do you want? (default is 1 day):")
	var days int = 1
	_, _ = fmt.Scanln(&days)

	API_KEY := os.Getenv("API_KEY")
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=" + API_KEY + "&q=" + q + "&days=" + strconv.Itoa(days) + "&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic("Failed to fetch weather data")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather weatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		panic(err)
	}

	location, current := weather.Location, weather.Current

	fmt.Printf(
		"%s, %s: %.0fC, %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	for _, day := range weather.Forecast.Forecastday {
		fmt.Printf("\nWeather forecast for %s:\n", day.Date)
		for _, hour := range day.Hour {
			date := time.Unix(hour.TimeEpoch, 0)

			if date.Before(time.Now()) {
				continue
			}

			message := fmt.Sprintf(
				"%s - %.0fC, %.0f%%, %s\n",
				date.Format("15:04"),
				hour.TempC,
				hour.ChanceOfRain,
				hour.Condition.Text,
			)

			if hour.ChanceOfRain < 40 {
				fmt.Print(message)
			} else {
				color.Red(message)
			}
		}
	}
}
