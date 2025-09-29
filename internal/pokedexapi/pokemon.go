package pokedexapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rounakkumarsingh/pokedex/internal/pokecache"
)

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Types          []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

var pokemonCache pokecache.Cache = pokecache.NewCache(24 * time.Hour)

func GetPokemonDetails(pokemonName string) (Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)
	e, ok := pokemonCache.Get(url)
	if ok {
		var v Pokemon
		if err := json.Unmarshal(e, &v); err != nil {
			return Pokemon{}, err
		}
		return v, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	if res.StatusCode == 404 {
		return Pokemon{}, errors.New("Not Found!!")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}
	res.Body.Close()
	var v Pokemon
	err = json.Unmarshal(body, &v)
	if err != nil {
		return Pokemon{}, err
	}
	e, err = json.Marshal(v)
	if err != nil {
		return Pokemon{}, err
	}
	pokemonCache.Add(url, e)
	return v, nil
}
