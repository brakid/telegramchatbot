package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

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

func Extract(s string) (string, float32, string, error) {
	r := regexp.MustCompile(`(?P<Name>.+) +(?P<Amount>\d+((\.|,)\d+)?) *(?P<Currency>€|EUR)?`)
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
