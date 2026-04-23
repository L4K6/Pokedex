package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, inputstr string) error
}

func commandExit(c *config, inputstr string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, inputstr string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, c := range getCommands() {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(c *config, inputstr string) error {
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

func commandBmap(c *config, inputstr string) error {
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

func commandExplore(c *config, location string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + location
	var pokemon LocationAreaStruct
	if len(location) == 0 {
		return fmt.Errorf("No location input")
	}
	if bytes, ok := c.Cache.Get(url); ok {
		data := bytes
		if err := json.Unmarshal(data, &pokemon); err != nil {
			return fmt.Errorf("Error Unmarshalling the response body: %w", err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error making a request: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode > 299 {
			return fmt.Errorf("Unknown Location")
		}

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading out the data: %w", err)
		}
		if err = json.Unmarshal(data, &pokemon); err != nil {
			return fmt.Errorf("Error Unmarshalling the response body: %w", err)
		}
		c.Cache.Add(url, data)
	}
	fmt.Println("Exploring " + location + "...")
	fmt.Println("Found Pokemon:")
	for _, name := range pokemon.PokemonEncounters {
		fmt.Println(" - " + name.Pokemon.Name)
	}
	return nil
}

func commandCatch(c *config, name string) error {
	fmt.Println("Throwing a Pokeball at " + name + "...")
	url := "https://pokeapi.co/api/v2/pokemon/" + name
	if len(name) == 0 {
		return fmt.Errorf("No Pokemon input")
	}
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error making a request: %w", err)
	}
	if res.StatusCode > 299 {
		return fmt.Errorf("Unknown Pokemon")
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error reading out the data: %w", err)
	}
	var catch Pokemon
	if err = json.Unmarshal(data, &catch); err != nil {
		return fmt.Errorf("Error Unmarshalling the response body: %w", err)
	}
	randonNumber := rand.Intn(catch.BaseExperience)
	fmt.Println(randonNumber)
	if 42 > randonNumber {
		fmt.Printf("%s was caught!\n", name)
		c.CaughtPokemon[name] = catch
	} else {
		fmt.Printf("%s esceaped!\n", name)
	}
	return nil
}

func commandInspect(c *config, name string) error {
	if len(name) == 0 {
		return fmt.Errorf("no pokemon name")
	}
	if pokemon, ok := c.CaughtPokemon[name]; !ok {
		return fmt.Errorf("Unknown Pokemon name")
	} else {
		fmt.Println("Name: " + pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("-%s: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, typeInfo := range pokemon.Types {
			fmt.Printf(" -%s\n", typeInfo.Type.Name)
		}
		return nil
	}
}

func commandPokedex(c *config, name string) error {
	fmt.Println("Your Pokedex:")
	for _, pok := range c.CaughtPokemon {
		fmt.Printf(" - %s\n", pok.Name)
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
		"explore": {
			name:        "explore",
			description: "Explores location and shows the pokemon that can be found there",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Tries to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a pokemon if it is within the pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows pokemon caught in the pokedex",
			callback:    commandPokedex,
		},
	}
}
