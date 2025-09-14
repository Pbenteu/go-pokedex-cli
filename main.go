package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pbenteu/pokedexcli/internal/pokeapi"
)

type GameState struct {
	PlayerLocationPage int
	Pokedex            map[string]pokeapi.Pokemon
}

type CommandHandler struct {
	description string
	callback    func(gameState *GameState, args ...string) error
}

var gameState = GameState{
	PlayerLocationPage: 0,
	Pokedex:            map[string]pokeapi.Pokemon{},
}

var commandConfigs = map[string]CommandHandler{}

func init() {
	commandConfigs = map[string]CommandHandler{
		"exit": {
			description: "Exit the Pokedex",
			callback:    handleExitCommand,
		},
		"help": {
			description: "Displays a help message",
			callback:    handleHelpCommand,
		},
		"map": {
			description: "Load the next avaiable locations",
			callback:    handleMapCommand,
		},
		"mapb": {
			description: "Load the previous avaiable locations",
			callback:    handleMapbCommand,
		},
		"explore": {
			description: "Explore a given location. Usage: 'explore eterna-city-area'",
			callback:    handleExploreCommand,
		},
		"catch": {
			description: "Try to catch a pokemon",
			callback:    handleCatchCommand,
		},
		"pokedex": {
			description: "List the pokemons on your pokedex",
			callback:    handlePokedexCommand,
		},
		"inspect": {
			description: "Inspect a pokemon you have caught",
			callback:    handleInspectCommand,
		},
	}
}

func CleanInput(text string) []string {
	text = strings.ToLower(text)
	return strings.Fields(text)
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("Pokedex > ")
		scanner.Scan()

		words := CleanInput(scanner.Text())

		command := words[0]

		commandConfig, ok := commandConfigs[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := commandConfig.callback(&gameState, words[1:]...)
		if err != nil {
			fmt.Printf("Error executing command:\n   %v\n", err)
		}
	}
}
