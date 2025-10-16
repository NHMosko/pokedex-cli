package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nhmosko/pokedexcli/internal/pokecache"
)

var globalCache *pokecache.Cache
var lastCommand func(string) error
var lastArgument string

func main() {
	initMap()
	scanner := bufio.NewScanner(os.Stdin)
	globalCache = pokecache.NewCache(5 * time.Second)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		inputWords := cleanInput(scanner.Text())

		if len(inputWords) < 1 {
			continue
		}

		inputCommand := inputWords[0]
		switch inputCommand {
		case "x":
			inputCommand = "exit"
		case "h":
			inputCommand = "help"
		case "m":
			inputCommand = "map"
		case "b":
			inputCommand = "mapb"
		case "e":
			inputCommand = "explore"
		case "c":
			inputCommand = "catch"
		case "p":
			inputCommand = "pokedex"
		case "a":
			inputCommand = "again"
		case "i":
			inputCommand = "inspect"
		}

		argument := ""
		if len(inputWords) > 1 {
			argument = inputWords[1]
		}

		cmd, ok := allCommands[inputCommand]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if inputCommand != "again" {
			lastCommand = cmd.callback
			lastArgument = argument
		}
		
		cmd.callback(argument)
	}
}

func cleanInput(text string) []string {
	var out []string
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, ",", " ")
	text = strings.Trim(text, " ")
	splitText := strings.Split(text, " ")
	for _, word := range splitText {
		//fmt.Println(word)
		word = strings.Trim(word, " ")
		if word != "" {
			out = append(out, word)
		}
	}
	//fmt.Println(out)
	return out
}
