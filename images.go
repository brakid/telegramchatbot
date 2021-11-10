package main

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleImage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (string, error) {
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
