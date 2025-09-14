package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pbenteu/pokedexcli/internal/pokecache"
)

type getPokemonRes struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effot"`
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
}

type Pokemon struct {
	Id             int
	Name           string
	BaseExperience int
	Height         int
	Weight         int
	Stats          struct {
		Hp             int
		Attack         int
		Defense        int
		SpecialAttack  int
		SpecialDefense int
		Speed          int
	}
	Types []string
}

func GetPokemon(name string) (Pokemon, error) {
	fullURL := baseURL + "/pokemon/" + name

	var jsonBytes []byte

	value, exists := pokecache.Cache.Get(fullURL)
	if exists {
		jsonBytes = value
	} else {
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			return Pokemon{}, err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return Pokemon{}, err
		}

		defer res.Body.Close()

		jsonBytes, err = io.ReadAll(res.Body)
		if err != nil {
			return Pokemon{}, err
		}

		pokecache.Cache.Add(fullURL, jsonBytes)
	}

	var response getPokemonRes
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return Pokemon{}, err
	}

	pokemon := Pokemon{
		Id:             response.Id,
		Name:           response.Name,
		BaseExperience: response.BaseExperience,
		Height:         response.Height,
		Weight:         response.Weight,
	}

	for _, item := range response.Stats {
		switch item.Stat.Name {
		case "hp":
			pokemon.Stats.Hp = item.BaseStat
		case "attack":
			pokemon.Stats.Attack = item.BaseStat
		case "defense":
			pokemon.Stats.Defense = item.BaseStat
		case "special-attack":
			pokemon.Stats.SpecialAttack = item.BaseStat
		case "special-defense":
			pokemon.Stats.SpecialDefense = item.BaseStat
		case "speed":
			pokemon.Stats.Speed = item.BaseStat
		}
	}

	for _, item := range response.Types {
		pokemon.Types = append(pokemon.Types, item.Type.Name)
	}

	return pokemon, nil
}
