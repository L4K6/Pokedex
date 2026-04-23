package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		cleanedInput := cleanInput(userInput)
		if len(cleanedInput) == 0 {
			continue
		}
		command, exists := getCommands()[cleanedInput[0]]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}
		if command.name == "explore" && len(cleanedInput) >= 2 {
			err := command.callback(cfg, cleanedInput[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		if command.name == "catch" && len(cleanedInput) >= 2 {
			err := command.callback(cfg, cleanedInput[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			err := command.callback(cfg, "")
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func cleanInput(text string) []string {
	lowerCaseText := strings.ToLower(text)
	textSlice := strings.Fields(lowerCaseText)
	return textSlice
}
