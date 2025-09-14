package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/pbenteu/pokedexcli/internal/pokeapi"
)

func handleExitCommand(gameState *GameState, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func handleHelpCommand(gameState *GameState, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	for command, config := range commandConfigs {
		fmt.Printf("%s: %s\n", command, config.description)
	}

	return nil
}

func handleMapCommand(gameState *GameState, args ...string) error {
	fmt.Println("Loading next locations...")

	gameState.PlayerLocationPage += 1

	locations, err := pokeapi.ListLocations(gameState.PlayerLocationPage)
	if err != nil {
		return err
	}

	for _, location := range locations {
		fmt.Println(location.Name)
	}

	return nil
}

func handleMapbCommand(gameState *GameState, args ...string) error {
	if gameState.PlayerLocationPage <= 1 {
		fmt.Println("You are already in the first location")
		return nil
	}

	fmt.Println("Loading previous locations...")

	gameState.PlayerLocationPage -= 1

	locations, err := pokeapi.ListLocations(gameState.PlayerLocationPage)
	if err != nil {
		return err
	}

	for _, location := range locations {
		fmt.Println(location.Name)
	}

	return nil
}

func handleExploreCommand(gameState *GameState, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: explore <location>")
	}

	response, err := pokeapi.GetLocation(args[0])
	if err != nil {
		return err
	}

	for _, value := range response.PokemonEncounters {
		fmt.Println(value.Pokemon.Name)
	}

	return nil
}

func handleCatchCommand(gameState *GameState, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: catch <pokemon>")
	}

	pokemonToCatch := args[0]

	response, err := pokeapi.GetPokemon(pokemonToCatch)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonToCatch)

	odds := 1 - ((math.Log10(float64(response.BaseExperience)) / 10) * 2)
	roll := rand.Float64()

	// fmt.Println(response)
	// fmt.Printf("odds: %.2f, roll: %.2f\n", odds, roll)

	if roll <= odds {
		fmt.Printf("You catched a %s!!!\n", pokemonToCatch)
		gameState.Pokedex[pokemonToCatch] = response
	} else {
		fmt.Printf("%s escaped...\n", pokemonToCatch)
	}

	return nil
}

func handlePokedexCommand(gameState *GameState, args ...string) error {
	for _, item := range gameState.Pokedex {
		fmt.Printf("- %s\n", item.Name)
	}

	return nil
}

func handleInspectCommand(gameState *GameState, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: inspect <pokemon>")
	}

	pokemon, exists := gameState.Pokedex[args[0]]
	if !exists {
		fmt.Println("Pokemon not found in pokedex:", args[0])
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	fmt.Printf("   -hp: %d\n", pokemon.Stats.Hp)
	fmt.Printf("   -attack: %d\n", pokemon.Stats.Hp)
	fmt.Printf("   -defense: %d\n", pokemon.Stats.Defense)
	fmt.Printf("   -special-attack: %d\n", pokemon.Stats.SpecialAttack)
	fmt.Printf("   -special-defense: %d\n", pokemon.Stats.SpecialDefense)
	fmt.Printf("   -speed: %d\n", pokemon.Stats.Speed)
	fmt.Println("Types:")
	for _, item := range pokemon.Types {
		fmt.Printf("   -%s\n", item)
	}

	return nil
}
