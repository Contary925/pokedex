package main

import (
	"github.com/Contary925/pokedex/internal/pokecache"
	"time"
)

type config struct {
	commands	map[string]cliCommand
	next		*string
	previous 	*string
	cache		*pokecache.Cache
	pokedex 	map[string]pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}


func main() {
	cmds := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:		 "help",
			description: "View help message",
			callback: 	 commandHelp,
		},
		"map": {
			name:		 "map",
			description: "View next map",
			callback: 	 commandMap,
		},
		"mapb": {
			name: 		 "mapb",
			description: "View previous map",
			callback:	 commandMapb,
		},
		"explore": {
			name:		 "explore",
			description: "View pokemons on the map",
			callback:	 commandExplore,
		},
		"catch": {
			name: 		 "catch",
			description: "Attempt to catch a pokemon",
			callback: 	 commandCatch,
		},
		"pokedex": {
			name: 		 "pokedex",
			description: "View your pokedex",
			callback: 	 commandPokedex,
		},
		"inspect": {
			name: 	 	 "inspect",
			description: "Inspect a pokemon in your pokedex",
			callback: 	 commandInspect,
		},
	}
	conf := config{
		commands: cmds,
		cache: pokecache.NewCache(5 * time.Second),
		pokedex: make(map[string]pokemon),
	}
	REPL(&conf)
}

