package pkg

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
)

var cache []Emoji = nil

func GetEmoji() ([]Emoji, error) {

	if cache != nil {
		return cache, nil
	}

	var res []Emoji

	response, err := http.Get("https://api.github.com/repositories/516731265/contents/Emojis/Objects")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &res)

	if err != nil {
		return nil, err
	}

	cache = res

	return res, nil

}

func GetRandomEmoji() (*Emoji, error) {

	data, err := GetEmoji()

	if err != nil {
		return nil, err
	}

	index := rand.Intn(len(data))

	return &data[index], nil
}
