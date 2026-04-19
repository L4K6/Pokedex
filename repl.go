package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	var conf config
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
		err := command.callback(&conf)
		if err != nil {
			fmt.Println(err)
		}

	}
}

func cleanInput(text string) []string {
	lowerCaseText := strings.ToLower(text)
	textSlice := strings.Fields(lowerCaseText)
	return textSlice
}
