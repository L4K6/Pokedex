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
		var stringInp string
		if command.name == "explore" || command.name == "catch" || command.name == "inspect" {
			if len(cleanedInput) < 2 {
				stringInp = ""
			} else {
				stringInp = cleanedInput[1]
			}
		}
		err := command.callback(cfg, stringInp)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func cleanInput(text string) []string {
	lowerCaseText := strings.ToLower(text)
	textSlice := strings.Fields(lowerCaseText)
	return textSlice
}
