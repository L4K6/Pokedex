package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, c := range getCommands() {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(c *config) error {
	url := c.Next
	var location locationAreaResp
	if c.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	if bytes, ok := c.Cache.Get(url); ok {
		data := bytes
		if err := json.Unmarshal(data, &location); err != nil {
			return fmt.Errorf("Error Unmarshalling the response body: %w", err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error making a request: %w", err)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading out the data: %w", err)
		}
		if err = json.Unmarshal(data, &location); err != nil {
			return fmt.Errorf("Error Unmarshalling the response body: %w", err)
		}
		c.Cache.Add(url, data)
	}

	c.Previous = location.Previous
	c.Next = location.Next
	for _, name := range location.Results {
		fmt.Println(name.Name)
	}
	return nil
}

func commandBmap(c *config) error {
	url := c.Previous
	var location locationAreaResp
	if c.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	if bytes, ok := c.Cache.Get(url); ok {
		data := bytes
		if err := json.Unmarshal(data, &location); err != nil {
			return fmt.Errorf("Error Unmarshalling the response body: %w", err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error making a request: %w", err)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading out the data: %w", err)
		}
		if err = json.Unmarshal(data, &location); err != nil {
			return fmt.Errorf("Error Unmarshalling the response body: %w", err)
		}
		c.Cache.Add(url, data)
	}

	c.Previous = location.Previous
	c.Next = location.Next
	for _, name := range location.Results {
		fmt.Println(name.Name)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays 20 locations",
			callback:    commandMap,
		},
		"bmap": {
			name:        "bmap",
			description: "Displays 20 previous locations",
			callback:    commandBmap,
		},
	}
}
