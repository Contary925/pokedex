package main

import (

)

type config struct {
	commands	map[string]cliCommand
	next		*string
	previous 	*string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}
	conf := config{commands: cmds}
	REPL(&conf)
}

