package main

import (
	"fmt"
	"os"
	"encoding/json"
	"net/http"
	"math/rand"
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

type statInfo struct {
	Name string `json:"name"`
}

type stat struct {
	Value int `json:"base_stat"`
	Info statInfo `json:"stat"`
}

type typeInfo struct {
	Name string `json:"name"`
}

type pokemonType struct {
	TypeInfo typeInfo `json:"type"`
}

type pokemon struct {
	Name string `json:"name"`
	BaseExperience int `json:"base_experience"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []stat `json:"stats"`
	Types []pokemonType `json:"types"`
}

type pokemonEncounter struct {
	Pokemon pokemon `json:"pokemon"`
}

type areaData struct {
	PokemonEncounters []pokemonEncounter `json:"pokemon_encounters"`
}

func commandExit(conf *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, args ...string) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: View the next map
mapb: View the previous map`)
	return nil
}

func commandMap(conf *config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if conf.next != nil {
		url = *conf.next
	}
	var data mapData
	var err error
	if cached_value, ok := conf.cache.Get(url); ok {
		err = json.Unmarshal(cached_value, &data)
		if err != nil {
			return err
		}
	} else {
		data, err = fetchMapFromUrl(url)
		if err != nil {
			return err
		}
		byte_data, err := json.Marshal(data)
		if err != nil {
			return err
		}
		conf.cache.Add(url, byte_data)
	}
	openMap(conf, data)
	return nil
}

func commandMapb(conf *config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if conf.previous != nil {
		url = *conf.previous
	} else {
		fmt.Println("You're on a first page!")
		return nil
	}
	var data mapData
	var err error
	if cached_value, ok := conf.cache.Get(url); ok {
		err = json.Unmarshal(cached_value, &data)
		if err != nil {
			return err
		}
	} else {
		data, err = fetchMapFromUrl(url)
		if err != nil {
			return err
		}
		byte_data, err := json.Marshal(data)
		if err != nil {
			return err
		}
		conf.cache.Add(url, byte_data)
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

func commandExplore(conf *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("No arguments!")
	}
	url := "https://pokeapi.co/api/v2/location-area/"+args[0]+"/"
	var data areaData
	var err error
	if cached_value, ok := conf.cache.Get(url); ok {
		err = json.Unmarshal(cached_value, &data)
		if err != nil {
			return err
		}
	} else {
		data, err = fetchAreaDataFromUrl(url)
		if err != nil {
			return err
		}
		byte_data, err := json.Marshal(data)
		if err != nil {
			return err
		}
		conf.cache.Add(url, byte_data)
	}
	if err != nil {
		return err
	}
	if len(data.PokemonEncounters) == 0 {
		fmt.Printf("No pokemons in this area!\n")
		return nil
	}
	fmt.Printf("Found pokemons:\n")
	for _, encounter := range data.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func fetchAreaDataFromUrl(url string) (areaData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return areaData{}, err
	}
	defer resp.Body.Close()
	var data areaData
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return areaData{}, err
	}
	return data, nil
}

func fetchPokemonDataFromUrl(url string) (pokemon, error) {
	resp, err := http.Get(url)
	if err != nil {
		return pokemon{}, err
	}
	defer resp.Body.Close()
	var data pokemon
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return pokemon{}, err
	}
	return data, nil
}

func commandCatch(conf *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("No arguments!")
	}
	url := "https://pokeapi.co/api/v2/pokemon/"+args[0]+"/"
	var pokemon pokemon
	var err error
	if cached_value, ok := conf.cache.Get(url); ok {
		err = json.Unmarshal(cached_value, &pokemon)
		if err != nil {
			return err
		}
	} else {
		pokemon, err = fetchPokemonDataFromUrl(url)
		if err != nil {
			return err
		}
		byte_data, err := json.Marshal(pokemon)
		if err != nil {
			return err
		}
		conf.cache.Add(url, byte_data)
	}
	if err != nil {
		return err
	}
	if pokemon.Name == "" {
		fmt.Printf("Unknown pokemon!\n")
		return nil
	}
	if _, ok := conf.pokedex[pokemon.Name]; ok {
		fmt.Printf("You have already caught this pokemon!\n")
		return nil
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon.Name)
	chance := 2500 / pokemon.BaseExperience
	//fmt.Printf("Success chance: %v\n", chance)
	roll := rand.Intn(100) //0-99
	success := roll < chance
	if success {
		fmt.Printf("Successfully caught %v!\n", pokemon.Name)
		conf.pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("Failed to catch %v...\n", pokemon.Name)
	}
	return nil
}

func commandPokedex(conf *config, args ...string) error {
	if len(conf.pokedex) == 0 {
		fmt.Printf("You have not caught any pokemons yet!\n")
		return nil
	}
	fmt.Printf("Pokemons in your pokedex:\n")
	for name, _ := range conf.pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandInspect(conf *config, args ...string) error {
	var pokemon pokemon
	if value, ok := conf.pokedex[args[0]]; !ok {
		fmt.Printf("You have not caught this pokemon yet!\n")
		return nil
	} else {
		pokemon = value
	}
	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%v: %v\n", stat.Info.Name, stat.Value)
	}
	fmt.Printf("Types:\n")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("  - %v\n", pokemonType.TypeInfo.Name)
	}
	return nil
}