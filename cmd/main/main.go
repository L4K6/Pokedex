package main

import (
	"time"

	"github.com/L4K6/Pokedex/internal/pokecache"
)

func main() {
	cfg := &config{
		Cache: pokecache.NewCache(5 * time.Second),
	}
	startRepl(cfg)
}
