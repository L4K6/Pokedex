package main

import "github.com/L4K6/Pokedex/internal/pokecache"

type config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
}

type locationAreaResp struct {
	Next     string
	Previous string
	Results  []locationArea
}

type locationArea struct {
	Name string
	URL  string
}
