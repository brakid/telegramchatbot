package main

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleTextMessage(message *tgbotapi.Message) (string, error) {
	name, amount, currency, err := ExtractSpending(message.Text)

	if err != nil {
		return "", err
	}

	spending := Spending{Name: name, Amount: amount, Date: time.Now(), Currency: currency}
	err = spendings.Add(&spending)

	if err != nil {
		return "", err
	}

	return FormatSpending(&spending), nil
}
