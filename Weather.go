package main

import (
	"fmt"
	"log"
	"time"

	owm "github.com/briandowns/openweathermap"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {
	// Получите токен вашего бота от BotFather в Telegram
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Fatal(err)
	}
	// Получите токен на сайте openweather
	apiKey := ""

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = "Привет! Я бот, который предоставляет информацию о погоде. Просто отправь мне название города, и я скажу тебе текущую погоду."
			case "help":
				msg.Text = "Просто отправь мне название города, и я скажу тебе текущую погоду."
			default:
				msg.Text = "Неизвестная команда. Попробуй /start или /help."
			}

			bot.Send(msg)
		} else {
			// Обработка текстовых сообщений
			city := update.Message.Text
			weatherText, err := getWeather(apiKey, city)
			if err != nil {
				weatherText = "Ошибка при получении погоды: " + err.Error()
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, weatherText)
			bot.Send(msg)
		}
	}
}

func getWeather(apiKey, city string) (string, error) {
	w, err := owm.NewCurrent("C", "RU", apiKey)
	if err != nil {
		return "", err
	}

	err = w.CurrentByName(city)
	if err != nil {
		return "", err
	}


	var loc *time.Location
	if w.Sys.Country != "" {
		timezoneOffset := int(w.Timezone)
		loc = time.FixedZone("", timezoneOffset)
	} else {
		// Use a default time zone if the country code is not available
		log.Println("Код страны недоступен.")
		loc = time.UTC
	}

	currentTime := time.Now().In(loc)


	temperature := int(w.Main.Temp)

	weatherText := fmt.Sprintf("Погода в %s:\nТемпература: %d°C\nОписание: %s\nТекущее время: %s",
		w.Name, temperature, w.Weather[0].Description, currentTime.Format("2006-01-02 15:04:05"))

	return weatherText, nil
}
