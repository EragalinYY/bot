package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type WeatherResponse struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Name string `json:"name"`
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func getWeather(city string) (string, error) {
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	url := "http://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=" + apiKey + "&units=metric"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var weatherResponse WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return "", err
	}

	weatherDescription := weatherResponse.Weather[0].Description
	temperature := weatherResponse.Main.Temp
	cityName := weatherResponse.Name

	return "Weather in " + cityName + ": " + weatherDescription + ", Temperature: " + fmt.Sprintf("%.2f", temperature) + "°C", nil
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я твой бот. Чем я могу помочь вам сегодня?")
				bot.Send(msg)
			case "help":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я могу помочь вам со следующими командами:\n/start - Запустите бота\n/weather - Запросить прогноз погоды\n/help - Отобразить справку")
				bot.Send(msg)
			case "weather":
				city := update.Message.CommandArguments()
				if city == "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, укажите название населенного пункта")
					bot.Send(msg)
					continue
				}
				weatherInfo, err := getWeather(city)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Не удалось получить информацию о погоде")
					bot.Send(msg)
					continue
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, weatherInfo)
				bot.Send(msg)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не знаю этой команды")
				bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			bot.Send(msg)
		}
	}
}
