package internal

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type AppResponse struct {
	Applist struct {
		Apps []struct {
			Appid int    `json:"appid"`
			Name  string `json:"name"`
		} `json:"apps"`
	} `json:"applist"`
}

// Query the Steam API to get a JSON response object.
func querySteamAPI() (AppResponse, error) {
	var steamResponse AppResponse
	client := http.Client{
		Timeout: time.Second * 30, // Timeout after 30 seconds
	}

	req, err := http.NewRequest(http.MethodGet, SteamApiUrl, nil)
	if err != nil {
		return steamResponse, err
	}

	req.Header.Set("User-Agent", "cryoutilities")

	res, err := client.Do(req)
	if err != nil {
		return steamResponse, err
	}

	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(res.Body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return steamResponse, err
	}

	err = json.Unmarshal(body, &steamResponse)
	if err != nil {
		return steamResponse, err
	}

	return steamResponse, nil
}

// Parse a SteamAPI response object and create a map of all games.
func generateGameMap(response AppResponse) map[int]string {
	gameMap := map[int]string{}

	for i := range response.Applist.Apps {
		gameMap[response.Applist.Apps[i].Appid] = response.Applist.Apps[i].Name
	}

	return gameMap
}
