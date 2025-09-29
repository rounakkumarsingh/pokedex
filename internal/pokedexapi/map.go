package pokedexapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rounakkumarsingh/pokedex/internal/pokecache"
)

var c pokecache.Cache = pokecache.NewCache(time.Hour * 1)

type LocationAreaAPIResponse struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
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
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetLocations(offset *int) []string {
	var locations []string
	currentOffset := *offset

	for len(locations) < 20 {
		var s string
		var err error

		for range 5 {
			s, err = GetLocation(currentOffset)
			if err == nil {
				break // Success, exit retry loop
			}
			// If this was the last retry, we'll move to next location
		}

		if err == nil {
			locations = append(locations, s)
		}
		// Always move to next location index, regardless of success/failure
		currentOffset++

		// Safety check to prevent infinite loop if API consistently fails
		if currentOffset > *offset+100 {
			break
		}
	}

	*offset = currentOffset
	return locations
}

func GetLocation(offSet int) (string, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%d/", offSet)
	e, ok := c.Get(url)
	if ok {
		return string(e), nil
	}
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}
	res.Body.Close()
	var v LocationAreaAPIResponse
	err = json.Unmarshal(body, &v)
	if err != nil {
		return "", nil
	}
	c.Add(url, []byte(v.Name))
	return v.Name, nil
}
