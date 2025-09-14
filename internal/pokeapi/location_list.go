package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pbenteu/pokedexcli/internal/pokecache"
)

type Location struct {
	Name string `json:"name"`
	Url  string `json:"Url"`
}

type ListLocationsRes struct {
	Results []Location `json:"results"`
}

func ListLocations(page int) ([]Location, error) {
	var response ListLocationsRes

	fullURL := baseURL + "/location-area"

	offset := (page - 1) * 20

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return []Location{}, err
	}

	query := url.Values{}
	query.Add("limit", "20")
	query.Add("offset", strconv.Itoa(offset))
	req.URL.RawQuery = query.Encode()

	cacheKey := req.URL.String()

	var byteData []byte

	cacheEntry, exists := pokecache.Cache.Get(cacheKey)
	if exists {

		byteData = cacheEntry
	} else {
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return []Location{}, err
		}

		defer res.Body.Close()

		byteData, err = io.ReadAll(res.Body)
		if err != nil {
			return []Location{}, err
		}

		pokecache.Cache.Add(cacheKey, byteData)
	}

	err = json.Unmarshal(byteData, &response)
	if err != nil {
		return []Location{}, err
	}

	return response.Results, nil
}
