package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var response string
	var err error

	if message.Command() != "" {
		response, err = HandleCommand(message)
	} else if len(message.Photo) > 0 {
		response, err = HandleImage(bot, message)
	} else {
		response, err = HandleTextMessage(message)
	}

	if err != nil {
		response = fmt.Sprintf("Error: %v", err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = tgbotapi.ModeHTML

	if _, err := bot.Send(msg); err != nil {
		log.Fatalln(err)
	}
}
