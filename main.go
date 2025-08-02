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
	"math/rand"
)

type config struct {
	Next string
	Previous string
	Location string
	Encounter string // Name of the pokemon
	Caught map[string]PokeData
	Cache    *pokecache.Cache
	Inspect string
}

type cliCommand struct {
	name string
	description string
	callback func(*config) error
}

type PokeLocation struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type PokeData struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        any `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []struct {
		Abilities []struct {
			Ability  any  `json:"ability"`
			IsHidden bool `json:"is_hidden"`
			Slot     int  `json:"slot"`
		} `json:"abilities"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"past_abilities"`
	PastTypes []any `json:"past_types"`
	Species   struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       string `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  string `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      string `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale string `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       string `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      string `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale string `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       string `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  string `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      string `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale string `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       string `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  string `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      string `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

type SavedPokemon struct {
	Name string
	Height int
	Weight int
	Stats  map[string]int
    Types  []string
}

func savePokemon(pokemonData PokeData) SavedPokemon {

	stats := make(map[string]int)
    for _, s := range pokemonData.Stats {
        stats[s.Stat.Name] = s.BaseStat
    }

    types := []string{}
    for _, t := range pokemonData.Types {
        types = append(types, t.Type.Name)
    }

    return SavedPokemon{
        Name:   pokemonData.Name,
        Height: pokemonData.Height,
        Weight: pokemonData.Weight,
		Stats:  stats,
        Types:  types,
    }
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

func commandExplore(param *config) error {
	
	location := param.Location

	fmt.Printf("Exploring %v...\n",location)

	pokeUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", location)
	
	resp, err := http.Get(pokeUrl)
	
	if err != nil {
		fmt.Println("Error calling the PokeURl")
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading the response body")
		return err
	}

	var rawOutput PokeLocation

	if err := json.Unmarshal(body, &rawOutput); err != nil {
		fmt.Println("Error unmarshaling the data")
		return err
	}
	
	fmt.Println("Found Pokemon:")

	for _,check := range rawOutput.PokemonEncounters {
		fmt.Printf("- %v \n",check.Pokemon.Name)
	}

	return nil
}

func cleanInput(text string) []string {
	fmt.Println(text)
	conversion := strings.Fields(strings.ToLower(text))
	return conversion
}

func catchProbability(experienceLevel int) bool {

	chance := rand.Intn(100)
	difficulty := experienceLevel / 5

	if difficulty >= 90 {
		difficulty = 90
	}

	return chance > difficulty
}

func commandCatch(param *config) error {
	
	pokemonName := strings.ToLower(param.Encounter)

	pokeUrl := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", pokemonName)
	
	fmt.Printf("Throwing a Pokeball at %v...\n",pokemonName)

	resp, err := http.Get(pokeUrl)
	
	if err != nil {
		fmt.Println("Error calling GET req on pokeAPI")
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Err reading response body")
		return err
	}

	var pokemonData PokeData

	err = json.Unmarshal(body, &pokemonData)

	if err != nil {
		fmt.Println("Error finding the pokemon, Make sure it exists")
		return err
	}


	if pokemonData.BaseExperience != 0 {

		if catchProbability(pokemonData.BaseExperience) {
			fmt.Printf("%v was caught!\n",pokemonName)
			
			if param.Caught == nil {
                param.Caught = make(map[string]PokeData)
            }

			param.Caught[strings.ToLower(pokemonData.Name)] = pokemonData

		} else {
			fmt.Printf("%v escaped!\n",pokemonName)
		}

	}
	return nil
}

func commandPokedex(param *config) error {
	fmt.Println("Your Pokedex:")

	if len(param.Caught) == 0 {
        fmt.Println(" - (no Pokémon caught yet)")
        return nil
    }

	for name := range param.Caught {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandInspect(param *config) error {
    pokemonName := strings.ToLower(param.Inspect)

    pokemonData, exists := param.Caught[pokemonName]
    if !exists {
        fmt.Println("Pokémon not found in your caught list.")
        return nil
    }

    saved := savePokemon(pokemonData)

    fmt.Printf("Name: %s\n", saved.Name)
    fmt.Printf("Height: %d\n", saved.Height)
    fmt.Printf("Weight: %d\n", saved.Weight)
	fmt.Println("Stats:")
    for stat, value := range saved.Stats {
        fmt.Printf("  - %s: %d\n", stat, value)
    }
    fmt.Printf("Types: %v\n", saved.Types)


    return nil
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
		"explore": {
			name: "explore",
			description: "takes in a map and then we explore",
			callback: commandExplore,
		},
		"catch": {
			name: "catch",
			description: "command to catch the pokemon",
			callback: commandCatch,
		},
		"inspect": {
			name: "inspect",
			description: "inspecting",
			callback: commandInspect,
		},
		"pokedex": {
			name: "pokedex",
			description: "pokedex",
			callback: commandPokedex,
		},
	}

    configPagination := &config{
		Cache: cache,
		Location: "",
		Encounter: "",
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
			
			if command.name == "explore" {

				if len(cleanInput) > 1 {
					configPagination.Location = cleanInput[1]
					command.callback(configPagination)
				} else {
					fmt.Println("User forgot to mention region")
				}

			} else if command.name == "catch" {

				if len(cleanInput) > 1 {
					configPagination.Encounter = cleanInput[1]
					command.callback(configPagination)
				} else {
					fmt.Println("No Pokemon encountered")
				}
			} else if command.name == "inspect" {

				if len(cleanInput) > 1 {
					configPagination.Inspect = cleanInput[1]
					command.callback(configPagination)
				} else {
					fmt.Println("No Pokemon mentioned")
				}
			} else {
				command.callback(configPagination)
			}

		} else {
			fmt.Println("Unknown Command")
		}
		
	}
}

