package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var allCommands = map[string]cliCommand{}

func helpCommand() error {
	if len(allCommands) == 0 {
		return fmt.Errorf("No commands found, program broken")
	}
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, v := range allCommands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func exitCommand() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	defer os.Exit(0)
	return nil
}

type Results struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type LocationArea struct {
    Count int `json:"count"`
    Next *string `json:"next"` 
    Previous *string `json:"previous"` 
    Results []Results `json:"results"` 
}

var locationAreas LocationArea

func mapCommand() error {
	url := "https://pokeapi.co/api/v2/location-area/"

	if locationAreas.Next != nil {
		url = *locationAreas.Next
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonData, &locationAreas); err != nil {
		return err
	}
	for _, location := range locationAreas.Results {
		fmt.Println(location.Name)
	}
	
	return nil
}

func mapbCommand() error {
	url := ""

	if locationAreas.Previous != nil {
		url = *locationAreas.Previous
	} else {
		fmt.Println("you're on the first page")
		return fmt.Errorf("No previous areas found")
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonData, &locationAreas); err != nil {
		return err
	}
	for _, location := range locationAreas.Results {
		fmt.Println(location.Name)
	}
	
	return nil
}


func initMap() {
	allCommands["exit"] = cliCommand{
		name:        "exit",
		description: "Exits the Pokedex",
		callback:    exitCommand,
	}
	allCommands["help"] = cliCommand{
		name:        "help",
		description: "Explains how to use the Pokedex",
		callback:    helpCommand,
	}
	allCommands["map"] = cliCommand{
		name:        "map",
		description: "Displays next 20 Pokemon areas",
		callback:    mapCommand,
	}
	allCommands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Displays previous 20 Pokemon areas",
		callback:    mapbCommand,
	}
}

func main() {
	initMap()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		inputWords := cleanInput(scanner.Text())

		if len(inputWords) < 1 {
			continue
		}

		inputCommand := inputWords[0]

		cmd, ok := allCommands[inputCommand]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		cmd.callback()
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
