package main

import (
	"fmt"
	"go-image-generator/pkg/utils"
)

func main() {
	testDates := []string{
		"2020-09-29",
		"2021-01-01",
		"2024-05-23",
		"2024-12-03",
		"2024-11-21",
		"2024-02-02",
	}

	for _, date := range testDates {
		formatted, err := utils.ParseEventDate(date)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", date, err)
		} else {
			fmt.Printf("%s -> %s\n", date, formatted)
		}
	}
}
