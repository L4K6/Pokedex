package main

import (
	"time"

	"github.com/L4K6/Pokedex/internal/pokecache"
)

func main() {
	cfg := &config{
		Cache:         pokecache.NewCache(50 * time.Second),
		CaughtPokemon: map[string]Pokemon{},
	}
	startRepl(cfg)
}
