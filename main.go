package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
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
