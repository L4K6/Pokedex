package main

type config struct {
	Next     string
	Previous string
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
