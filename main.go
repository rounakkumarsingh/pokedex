package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strings"

	"github.com/rounakkumarsingh/pokedex/internal/pokedexapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, ...string) error
}

type Config struct {
	mapOffset int
	pokemons  map[string]pokedexapi.Pokemon
}

var config Config

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "lists location",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "See the previous list",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore the pokemon",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
	}
}

func commandExit(_ *Config, _ ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return ErrExit
}

var ErrExit = errors.New("exit")

func commandHelp(_ *Config, _ ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	// stable order
	keys := make([]string, 0, len(commands))
	for k := range commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		c := commands[k]
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(config *Config, _ ...string) error {
	s := pokedexapi.GetLocations(&config.mapOffset)
	for _, v := range s {
		fmt.Println(v)
	}
	return nil
}

func commandMapb(config *Config, args ...string) error {
	config.mapOffset -= 40
	if config.mapOffset < 0 {
		config.mapOffset = 0
		fmt.Println("You are at the start")
		return nil
	}
	return commandMap(config)
}

func commandExplore(_ *Config, args ...string) error {
	fmt.Printf("Exploring %s...\n", args[1])
	pokemons, err := pokedexapi.GetPokemons(args[1])
	if err != nil {
		return err
	}
	for _, pokemon := range pokemons {
		fmt.Printf(" - %s\n", pokemon)
	}
	return nil
}

func commandCatch(_ *Config, args ...string) error {
	pokemon := args[1]
	pokemonInfo, err := pokedexapi.GetPokemonDetails(pokemon)
	if err != nil {
		if err.Error() == "Not Found!!" {
			fmt.Println("IDK what that is, but it ain't a Pokemon")
			return nil
		}
		return err
	}
	fmt.Printf("Throwing a pokeball at %s...\n", pokemon)
	prob := 1.0 / (1.0 + math.Exp(-(float64(pokemonInfo.BaseExperience))))
	u := rand.Float64()
	if u < prob {
		fmt.Printf("%s was caught!\n", pokemon)

	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}
	return nil
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		s.Scan()
		args := cleanInput(s.Text())
		switch args[0] {
		case "exit":
			if err := commands["exit"].callback(&config); err != ErrExit {
				panic(err)
			}
		case "help":
			if err := commandHelp(&config); err != nil {
				panic(err)
			}
		case "map":
			if err := commandMap(&config); err != nil {
				panic(err)
			}
		case "mapb":
			if err := commandMapb(&config); err != nil {
				panic(err)
			}
		case "explore":
			if err := commandExplore(&config, args...); err != nil {
				panic(err)
			}
		case "catch":
			if err := commandCatch(&config, args...); err != nil {
				panic(err)
			}
		}
	}
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	return strings.Fields(text)
}
