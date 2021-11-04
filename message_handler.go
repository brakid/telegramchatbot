package main

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func fetchImage(imageUrl string) ([]byte, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	var response string
	var err error

	if message.Command() != "" {
		response, err = handleCommand(message)
	} else if len(message.Photo) > 0 {
		response, err = handleImage(bot, message)
	} else {
		response, err = handleTextMessage(message)
	}

	if err != nil {
		response = fmt.Sprintf("Error: %v", err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	if _, err := bot.Send(msg); err != nil {
		log.Fatalln(err)
	}
}

func handleTextMessage(message *tgbotapi.Message) (string, error) {
	name, amount, currency, err := Extract(message.Text)

	if err != nil {
		return "", err
	}

	spending := Spending{Name: name, Amount: amount, Date: time.Now(), Currency: currency}
	err = spendings.Add(&spending)

	if err != nil {
		return "", err
	}

	total, err := spendings.TotalAmount()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Spent in total: %.2f", total), nil
}

func handleImage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (string, error) {
	fmt.Println("Image received")
	photos := message.Photo

	imageUrl, err := bot.GetFileDirectURL(photos[len(photos)-1].FileID)
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(imageUrl, ".jpg") && !strings.HasSuffix(imageUrl, ".jpeg") {
		return "", fmt.Errorf("unsupported file type")
	}

	imageBytes, err := fetchImage(imageUrl)
	if err != nil {
		return "", err
	}

	fmt.Println(imageBytes[0:100])

	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return "", err
	}

	fmt.Println(img)

	return "Image was processed", nil
}

func handleCommand(message *tgbotapi.Message) (string, error) {
	command := message.Command()
	argString := message.CommandArguments()

	switch command {
	case "configure":
		args := strings.Split(argString, " ")
		if len(args) != 2 {
			return "", fmt.Errorf("invalid args length")
		}

		fmt.Printf("%s = %s", args[0], args[1])
		configuration.Set(args[0], args[1])
		err := configuration.Save()
		if err != nil {
			return "", err
		}

		return "Configuration changed", nil
	case "show":
		return fmt.Sprint(configuration), nil
	case "about":
		return "Spending tracker bot", nil
	case "start":
		return "Spending tracker bot", nil
	case "spendings":
		{
			spendingValues, err := spendings.AllSpendings()

			if err != nil {
				return "", err
			}

			return fmt.Sprintf("Spendings: %v", spendingValues), nil
		}
	default:
		return "Invalid command", nil
	}
}
