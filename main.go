package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
