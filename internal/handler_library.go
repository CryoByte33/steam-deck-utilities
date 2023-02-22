package internal

import (
	"os"
	"strconv"
	"strings"

	"github.com/andygrunwald/vdf"
)

type Library struct {
	Path           string
	InstalledGames []int
}

func (lib *Library) listGames() {
	for _, game := range lib.InstalledGames {
		CryoUtils.InfoLog.Println(game)
	}
}

// Get a list of installed games (parseVDF), then ensure their folders in compatdata and
// shaderdata get moved to the appropriate location. Then, create a symlink on the SSD.
func findDataFolders() ([]Library, error) {
	libraries, err := parseVDF(LibraryVDFLocation)
	if err != nil {
		CryoUtils.ErrorLog.Println("Error parsing VDF at", LibraryVDFLocation)
		return nil, err
	}

	CryoUtils.InfoLog.Println("Loading libraries saved in VDF...")
	for x := range libraries {
		// If the library was added manually, remove the end of the path
		if strings.HasSuffix(libraries[x].Path, "SteamLibrary") {
			CryoUtils.InfoLog.Println("Found manually added library at", libraries[x].Path)
			libraries[x].Path = strings.ReplaceAll(libraries[x].Path, "/SteamLibrary", "")
		} else {
			CryoUtils.InfoLog.Println("Found library at", libraries[x].Path)
		}
	}

	return libraries, nil
}

// Parse the specified VDF file and return a slice of libraries.
func parseVDF(file string) ([]Library, error) {
	var libraries []Library

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	p := vdf.NewParser(f)
	m, err := p.Parse()
	if err != nil {
		return nil, err
	}

	for _, library := range m["libraryfolders"].(map[string]interface{}) {
		var installedGames []int

		for game := range library.(map[string]interface{})["apps"].(map[string]interface{}) {
			intGame, _ := strconv.Atoi(game)
			installedGames = append(installedGames, intGame)
		}
		newLib := Library{
			Path:           library.(map[string]interface{})["path"].(string),
			InstalledGames: installedGames,
		}
		libraries = append(libraries, newLib)
	}
	return libraries, nil
}
