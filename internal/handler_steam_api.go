// CryoUtilities
// Copyright (C) 2023 CryoByte33 and contributors to the CryoUtilities project

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
			Name  string `json:"name"`
			Appid int    `json:"appid"`
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
		defer res.Body.Close()
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
