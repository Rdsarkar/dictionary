package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var dictionary map[string]bool

func loadDictionaryFromURL(url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	dictionary = make(map[string]bool)
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		word := scanner.Text()
		dictionary[word] = true
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func spellCheck(word string) string {
	if dictionary[word] {
		return word
	}

	// If the word is not found, attempt to find a similar word
	for dictWord := range dictionary {
		if strings.HasPrefix(dictWord, word[:len(word)-1]) {
			return dictWord
		}
	}

	// If no similar word is found, return an empty string
	return ""
}

func main() {
	// Load the dictionary from the raw URL
	err := loadDictionaryFromURL("https://raw.githubusercontent.com/first20hours/google-10000-english/master/google-10000-english-usa-no-swears.txt")
	if err != nil {
		fmt.Println("Error loading dictionary:", err)
		return
	}

	// Prompt the user for input until they enter "Stop" (case insensitive)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter a word (\"Stop\" to stop): ")
		word, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		// Clean up the input
		word = strings.TrimSpace(word)

		// If the input word is empty, continue the loop
		if word == "" {
			continue
		}

		// Check if the user wants to stop
		if word == "Stop" {
			break
		}

		// Check the spelling
		correctWord := spellCheck(word)
		if correctWord != "" {
			fmt.Printf("Did you mean '%s'?\n", correctWord)
		} else {
			fmt.Println("Word not found in dictionary.")
		}
	}
}
