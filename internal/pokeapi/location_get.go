package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationDetails struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string
			URL  string
		}
	} `json:"pokemon_encounters"`
}

func GetLocation(name string) (LocationDetails, error) {
	fullURL := baseURL + "/location-area/" + name

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return LocationDetails{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return LocationDetails{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("error getting location: %s", res.Status)
		return LocationDetails{}, err
	}

	jsonBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationDetails{}, err
	}

	var response LocationDetails

	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return LocationDetails{}, err
	}

	return response, nil
}
