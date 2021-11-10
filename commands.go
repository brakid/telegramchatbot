package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleCommand(message *tgbotapi.Message) (string, error) {
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
	case "list":
		return listSpendings(argString)
	case "file":
		return "Download file from ....", nil
	case "edit":
		return editSpending(argString)
	case "delete":
		return deleteSpending(argString)
	default:
		return "Invalid command", nil
	}
}

func listSpendings(s string) (string, error) {
	year, month, err := getYearAndMonth(s)

	if err != nil {
		return "", err
	}

	spendingValues, err := spendings.AllSpendings(year, month)

	if err != nil {
		return "", err
	}

	return FormatSpendings(year, month, spendingValues), nil
}

func getYearAndMonth(s string) (int, time.Month, error) {
	if s == "" {
		year, month, _ := time.Now().Date()
		return year, month, nil
	} else {
		return ExtractMonth(s)
	}
}

func editSpending(s string) (string, error) {
	id, nameAmount, err := ExtractIdAndNameAmount(s)
	if err != nil {
		return "", err
	}
	spending, err := spendings.Get(id)
	if err != nil {
		return "", err
	}

	if nameAmount.IsName {
		spending.Name = nameAmount.Name
	} else {
		spending.Amount = nameAmount.Amount
	}

	err = spendings.Update(spending)
	if err != nil {
		return "", err
	}

	return FormatUpdatedSpending(spending), nil
}

type NameAmount struct {
	Name   string
	Amount float32
	IsName bool
}

func deleteSpending(s string) (string, error) {
	id, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return "", err
	}

	err = spendings.Delete(uint(id))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Deleted Spending with ID %v", id), nil
}
