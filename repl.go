package main

import "strings"

func cleanInput(text string) []string {
	lowerCaseText := strings.ToLower(text)
	textSlice := strings.Fields(lowerCaseText)
	return textSlice
}
