package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"net/http"
	"io"
	"encoding/json"
	"time"
	"github.com/girik21/pokedexcli/internal/pokecache"
)

type config struct {
	Next string
	Previous string
	Cache    *pokecache.Cache
}

type cliCommand struct {
	name string
	description string
	callback func(*config) error
}

func commandExit(param *config) error {
	fmt.Printf("Closing the Pokedex... Goodbye! \n")
	os.Exit(0)
	return nil
}

func commandHelp(param *config) error {
	fmt.Printf("Welcome to the Pokedex! \n")
	fmt.Print("Usage: \n")

	fmt.Printf("\n")
	fmt.Printf("\n")
	
	fmt.Printf("map: shows the map of pokemon\n")
	fmt.Printf("help: Displays a help message \n")
	fmt.Printf("exit: Exit the Pokedex \n")
	return nil
}

func commandBack(param *config) error {
	pokeUrl := "https://pokeapi.co/api/v2/location-area/?limit=20"
	if param.Previous != "" {
		pokeUrl = param.Previous
	}

	var body []byte
	var err error

	// ✅ Check cache first
	if data, ok := param.Cache.Get(pokeUrl); ok {
		body = data
	} else {
		resp, err := http.Get(pokeUrl)
		if err != nil {
			fmt.Println("Error calling the PokeURL")
			return err
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		param.Cache.Add(pokeUrl, body)
	}

	var rawOutput map[string]interface{}
	if err = json.Unmarshal(body, &rawOutput); err != nil {
		return err
	}

	if output, ok := rawOutput["results"].([]interface{}); ok {
		for _, item := range output {
			if location, ok := item.(map[string]interface{}); ok {
				fmt.Println(location["name"])
			}
		}
	}

	if next, ok := rawOutput["next"].(string); ok {
		param.Next = next
	}
	if previous, ok := rawOutput["previous"].(string); ok {
		param.Previous = previous
	}
	return nil
}

func commandMap(param *config) error {
	pokeUrl := "https://pokeapi.co/api/v2/location-area/?limit=20"
	if param.Next != "" {
		pokeUrl = param.Next
	}

	var body []byte
	var err error

	// ✅ Check cache first
	if data, ok := param.Cache.Get(pokeUrl); ok {
		body = data
	} else {
		resp, err := http.Get(pokeUrl)
		if err != nil {
			fmt.Println("Error calling the API:", err)
			return err
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading the body")
			return err
		}

		param.Cache.Add(pokeUrl, body)
	}

	var rawOutput map[string]interface{}
	if err = json.Unmarshal(body, &rawOutput); err != nil {
		fmt.Println("Error unmarshaling the body")
		return err
	}

	if result, ok := rawOutput["results"].([]interface{}); ok {
		for _, item := range result {
			if location, ok := item.(map[string]interface{}); ok {
				fmt.Println(location["name"])
			}
		}
	}

	if next, ok := rawOutput["next"].(string); ok {
		param.Next = next
	}
	if previous, ok := rawOutput["previous"].(string); ok {
		param.Previous = previous
	}
	return nil
}

func cleanInput(text string) []string {
	fmt.Println(text)
	conversion := strings.Fields(strings.ToLower(text))
	return conversion
}


func main() {
	userInput := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(5 * time.Second)

	commands := map[string]cliCommand{
		"exit": {
			name: "exit",
			description: "Exit the pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "User is asking for help",
			callback: commandHelp,
		},
		"map": {
			name: "map",
			description: "User wants to see the map",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "User wants to go back to the prev page",
			callback: commandBack,
		},
	}

    configPagination := &config{
		Cache: cache,
	}


	for {
		fmt.Print("Pokedex > ") // Printing the REPL to show that the pokdex started

		userInput.Scan() // Scans the user input

		cleanInput := strings.Fields(strings.ToLower(userInput.Text())) // Cleaning the text Feilds removes the space between the content and ToLower converts it to lowercase

		if len(cleanInput) == 0 { // Handles the panic gracefully if the user pressed enter without doing something
			continue
		}

		userCommand := cleanInput[0]
		
		if command, exists := commands[userCommand]; exists {
			command.callback(configPagination)
		} else {
			fmt.Println("Unknown Command")
		}
		
	}
}

