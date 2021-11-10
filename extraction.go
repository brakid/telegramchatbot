package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const NAME_REGEX = `(?P<Name>.+)`
const AMOUNT_REGEX = `(?P<Amount>\d+((\.|,)\d+)?)`
const CURRENCY_REGEX = `(?P<Currency>€|EUR)`
const MONTH_REGEX = `(?P<Month>\d{1,2})`
const YEAR_REGEX = `(?P<Year>\d{4})`
const ID_REGEX = `(?P<ID>\d+)`

var newSendingRegex = fmt.Sprintf("%v +%v *%v?", NAME_REGEX, AMOUNT_REGEX, CURRENCY_REGEX)
var monthAndYearRegex = fmt.Sprintf(`%v\.%v`, MONTH_REGEX, YEAR_REGEX)
var idAndNameAmount = fmt.Sprintf(`%v +((%v *%v)|%v)`, ID_REGEX, AMOUNT_REGEX, CURRENCY_REGEX, NAME_REGEX)

func ExtractSpending(s string) (string, float32, string, error) {
	r := regexp.MustCompile(newSendingRegex)
	matches := getMatches(r.FindStringSubmatch(s), r.SubexpNames())

	fmt.Println(matches)

	if len(*matches) < 2 {
		return "", 0.0, "", fmt.Errorf("no data extracted")
	}

	currency, ok := (*matches)["Currency"]

	if !ok {
		currency, _ = configuration.Get("defaultCurrency")
	} else {
		if currency == "€" {
			currency = "EUR"
		}
	}

	amount, err := parseAmount((*matches)["Amount"])
	if err != nil {
		return "", 0.0, "", fmt.Errorf("amount parsing failed")
	}

	return strings.TrimSpace((*matches)["Name"]), amount, currency, nil
}

func ExtractMonth(s string) (int, time.Month, error) {
	r := regexp.MustCompile(monthAndYearRegex)
	matches := getMatches(r.FindStringSubmatch(s), r.SubexpNames())

	if len(*matches) != 2 {
		return 0, time.January, fmt.Errorf("invalid date format, expecting MM.YYYY")
	}

	monthString := (*matches)["Month"]
	monthInt, err := strconv.Atoi(monthString)
	if err != nil {
		return 0, time.January, err
	}
	month := time.Month(monthInt)

	yearString := (*matches)["Year"]
	year, err := strconv.Atoi(yearString)
	if err != nil {
		return 0, time.January, err
	}

	return year, month, nil
}

func ExtractIdAndNameAmount(s string) (uint, *NameAmount, error) {
	r := regexp.MustCompile(idAndNameAmount)
	matches := getMatches(r.FindStringSubmatch(s), r.SubexpNames())

	if len(*matches) < 2 {
		return 0, nil, fmt.Errorf("invalid format, expecting <ID> <name>|<amount with currency>")
	}

	idString, ok := (*matches)["ID"]
	if !ok {
		return 0, nil, fmt.Errorf("no id found")
	}
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		return 0, nil, fmt.Errorf("id is not a number")
	}

	var nameAmount NameAmount
	name, ok := (*matches)["Name"]
	if ok {
		nameAmount.Name = strings.TrimSpace(name)
		nameAmount.IsName = true
		return uint(id), &nameAmount, nil
	}

	amountString, ok := (*matches)["Amount"]
	if !ok {
		return 0, nil, fmt.Errorf("no amount is provided")
	}
	amount, err := parseAmount(amountString)
	if err != nil {
		return 0, nil, fmt.Errorf("amount is not a number")
	}
	nameAmount.Amount = amount
	return uint(id), &nameAmount, nil
}

func getMatches(submatches []string, names []string) *map[string]string {
	result := make(map[string]string)

	for index, _ := range submatches {
		submatch := submatches[index]
		name := names[index]

		if name != "" && submatch != "" {
			result[name] = submatch
		}
	}

	return &result
}

func parseAmount(value string) (float32, error) {
	amount, err := strconv.ParseFloat(strings.ReplaceAll(value, ",", "."), 32)
	if err != nil {
		return 0.0, fmt.Errorf("amount parsing failed")
	}

	return float32(amount), nil
}
