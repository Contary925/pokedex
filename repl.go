package main

import (
	"strings"
	"bufio"
	"fmt"
	"os"
)

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func REPL(conf *config) {
	scanner := bufio.NewScanner(os.Stdin)
	cmds := conf.commands
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		response := scanner.Text()
		response_split := cleanInput(response)
		if len(response_split) == 0 {
			continue
		}
		command := response_split[0]
		res, ok := cmds[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := res.callback(conf, response_split[1:]...)
		if err != nil {
			fmt.Printf("Command '%s' failed with error: %s\n", command, err)
		}
	}
}