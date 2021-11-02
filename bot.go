package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var configuration *Configuration
var spendings *Spendings

func main() {
	godotenv.Load()

	botApiKey := os.Getenv("BOT_API_KEY")
	chatId, err := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if err != nil {
		log.Fatalln(err)
	}

	bot, err := tgbotapi.NewBotAPI(botApiKey)
	if err != nil {
		log.Fatalln(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	configuration, err = LoadConfiguration()
	if err != nil {
		log.Fatalln(err)
	}

	defer configuration.Save()

	spendings = New()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Chat.ID != chatId {
			continue
		}

		HandleMessage(bot, update.Message)
	}
}
