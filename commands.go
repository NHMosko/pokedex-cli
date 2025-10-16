package main

import (
	"math/rand/v2"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(string) error
}

var allCommands = map[string]cliCommand{}


func initMap() {
	allCommands["exit"] = cliCommand{
		name:        "x - exit",
		description: "Exits the Pokedex",
		callback:    exitCommand,
	}
	allCommands["help"] = cliCommand{
		name:        "h - help",
		description: "Explains how to use the Pokedex",
		callback:    helpCommand,
	}
	allCommands["map"] = cliCommand{
		name:        "m - map",
		description: "Displays next 20 Pokemon areas",
		callback:    mapCommand,
	}
	allCommands["mapb"] = cliCommand{
		name:        "b - mapb",
		description: "Displays previous 20 Pokemon areas",
		callback:    mapbCommand,
	}
	allCommands["explore"] = cliCommand{
		name:        "e - explore",
		description: "Displays all Pokemon located at the selected area",
		callback:    exploreCommand,
	}
	allCommands["catch"] = cliCommand{
		name:        "c - catch",
		description: "Attempts to catch the selected pokemon",
		callback:    catchCommand,
	}
	allCommands["inspect"] = cliCommand{
		name:        "i - inspect",
		description: "Gives information about a caught pokemon",
		callback:    inspectCommand,
	}
	allCommands["pokedex"] = cliCommand{
		name:        "p - pokedex",
		description: "Shows all caught pokemon",
		callback:    pokedexCommand,
	}
	allCommands["again"] = cliCommand{
		name:        "a - again",
		description: "Repeats last command",
		callback:    againCommand,
	}

	caughtPokemon = make(map[string]Pokemon)
}


func helpCommand(string) error {
	if len(allCommands) == 0 {
		return fmt.Errorf("No commands found, program broken")
	}
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	fmt.Println()
	for _, v := range allCommands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func exitCommand(string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	defer os.Exit(0)
	return nil
}

type AreaResults struct {
	Name string `json:"name"`
	Url string `json:"url"`
}

type LocationResponse struct {
    Count int `json:"count"`
    Next *string `json:"next"` 
    Previous *string `json:"previous"` 
    Results []AreaResults `json:"results"` 
}

var locationResponse LocationResponse

func mapCommand(string) error {
	url := "https://pokeapi.co/api/v2/location-area/"

	if locationResponse.Next != nil {
		url = *locationResponse.Next
	}

	jsonData, err := getHTTPResponse(url)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonData, &locationResponse); err != nil {
		return err
	}
	for _, location := range locationResponse.Results {
		fmt.Println(location.Name)
	}
	
	return nil
}

func mapbCommand(string) error {
	url := ""

	if locationResponse.Previous != nil {
		url = *locationResponse.Previous
	} else {
		fmt.Println("No previous areas found")
		return fmt.Errorf("No previous areas found")
	}

	jsonData, err := getHTTPResponse(url)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonData, &locationResponse); err != nil {
		return err
	}
	for _, location := range locationResponse.Results {
		fmt.Println(location.Name)
	}
	
	return nil
}


func getHTTPResponse(url string) ([]byte, error) {
	if val, ok := globalCache.Get(url); ok {
		return val, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	jsonData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	globalCache.Add(url, jsonData)
	return jsonData, nil
}


type Stat struct {
	Name string `json:"name"` 
}

type PokemonStat struct {
	Info Stat `json:"stat"`
	Value int `json:"base_stat"`
}

type Type struct {
	Name string `json:"name"`
}

type PokemonType struct {
	Type Type `json:"type"`
}

type Pokemon struct {
	Name string `json:"name"`
	Base_experience int `json:"base_experience"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []PokemonStat `json:"stats"`
	Types []PokemonType `json:"types"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type SpecificAreaResponse struct {
	Pokemon_encounters []PokemonEncounter `json:"pokemon_encounters"`
}

var caughtPokemon map[string]Pokemon

func exploreCommand(areaName string) error {
	if areaName == "" {
		fmt.Println("Specify which area to explore")
		return fmt.Errorf("No area given")
	}

	fullUrl := "https://pokeapi.co/api/v2/location-area/" + areaName

	jsonData, err := getHTTPResponse(fullUrl)
	if err != nil {
		return err
	}

	var areaResponse SpecificAreaResponse
	if err := json.Unmarshal(jsonData, &areaResponse); err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range areaResponse.Pokemon_encounters {
		fmt.Println("-", encounter.Pokemon.Name)
	}
	return nil
}

func catchCommand(pokemonName string) error {
	if pokemonName == "" {
		fmt.Println("Specify which pokemon to catch")
		return fmt.Errorf("No pokemon given")
	}

	fullUrl := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	jsonData, err := getHTTPResponse(fullUrl)
	if err != nil {
		return err
	}

	var pokemon Pokemon
	if err := json.Unmarshal(jsonData, &pokemon); err != nil {
		fmt.Printf("%v might not be a pokemon\n")
		return err
	}	

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if roll := rand.IntN(pokemon.Base_experience + 100); roll > pokemon.Base_experience {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		caughtPokemon[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	return nil
}

func inspectCommand(pokemonName string) error {
	if pokemonName == "" {
		fmt.Println("Specify which pokemon to inspect")
		return fmt.Errorf("No pokemon given")
	}

	pokemon, ok := caughtPokemon[pokemonName]
	if !ok {
		fmt.Println("You have not caught that pokemon")
		return nil
	}

	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("-%s: %v\n", stat.Info.Name, stat.Value)
	}
	fmt.Println("Types:")
	for _, pokeType := range pokemon.Types {
		fmt.Printf("- %s\n", pokeType.Type.Name)
	}

	return nil
}

func pokedexCommand(string) error {
	if len(caughtPokemon) == 0 {
		fmt.Println("Empty... use 'catch' to get new pokemon")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for k, _ := range caughtPokemon {
		fmt.Println("-", k)
	}
	return nil
}

func againCommand(string) error {
	err := lastCommand(lastArgument)
	if err != nil {
		fmt.Println("Failed to repeat command")
		return err
	}
	return nil
}
