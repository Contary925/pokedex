package main

import (
	"fmt"
	"os"
	"encoding/json"
	"net/http"
)

type area struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type mapData struct {
	Count    int     `json:"count"`
	Next     *string  `json:"next"`
	Previous *string  `json:"previous"`
	Results  []area  `json:"results"`
}

func commandExit(conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: View the next map
mapb: View the previous map`)
	return nil
}

func commandMap(conf *config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if conf.next != nil {
		url = *conf.next
	}
	data, err := fetchMapFromUrl(url)
	if err != nil {
		return err
	}
	openMap(conf, data)
	return nil
}

func commandMapb(conf *config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if conf.previous != nil {
		url = *conf.previous
	} else {
		fmt.Println("You're on a first page!")
		return nil
	}
	data, err := fetchMapFromUrl(url)
	if err != nil {
		return err
	}
	openMap(conf, data)
	return nil
}

func fetchMapFromUrl(url string) (mapData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return mapData{}, err
	}
	defer resp.Body.Close()
	var data mapData
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return mapData{}, err
	}
	return data, nil
}

func openMap(conf *config, data mapData) {
	conf.previous = data.Previous
	conf.next = data.Next
	for _, area := range data.Results {
		fmt.Println(area.Name)
	}
	return
}