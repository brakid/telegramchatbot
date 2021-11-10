package main

import (
	"fmt"
	"strings"
	"time"
)

func FormatSpendings(year int, month time.Month, spendingValues *[]Spending) string {
	var stringBuilder strings.Builder

	stringBuilder.WriteString(fmt.Sprintf("<b>Spendings %v %v:</b>\n", month.String(), year))

	total := float32(0.0)

	for _, spending := range *spendingValues {
		total += spending.Amount
		stringBuilder.WriteString(fmt.Sprintf("%v <i>%v</i> <b>%v</b>: %v%v\n", spending.Date.Format("02.01.2006"), spending.ID, spending.Name, spending.Amount, spending.Currency))
	}

	stringBuilder.WriteString("_________________\n")
	stringBuilder.WriteString(fmt.Sprintf("<b>Total</b>: %.2fEUR", total))

	return stringBuilder.String()
}

func FormatSpending(spending *Spending) string {
	return fmt.Sprintf("Added spending <i>%v: %.2f%v</i> with ID: <b>%v</b>", spending.Name, spending.Amount, spending.Currency, spending.ID)
}

func FormatUpdatedSpending(spending *Spending) string {
	return fmt.Sprintf("Updated spending <i>%v: %.2f%v</i> with ID: <b>%v</b>", spending.Name, spending.Amount, spending.Currency, spending.ID)
}
