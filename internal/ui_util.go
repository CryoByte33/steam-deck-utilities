package internal

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type GameStatus struct {
	GameName    string
	IsInstalled bool
}

// Show an error message over the main window.
func presentErrorInUI(err error, win fyne.Window) {
	CryoUtils.ErrorLog.Println(err)
	dialog.ShowError(err, CryoUtils.MainWindow)
}

// Create a CheckGroup of game data to allow for selection.
func createGameDataList() (*widget.CheckGroup, error) {
	cleanupList := widget.NewCheckGroup([]string{}, func(strings []string) {})
	cleanupList.Enable()
	cleanupList.Refresh()

	// Use the cached API Response if already present
	if CryoUtils.SteamAPIResponse == nil {
		// Make a map of all games stored in the steam API
		steamResponse, err := querySteamAPI()
		if err != nil {
			CryoUtils.ErrorLog.Println(err)
		}
		CryoUtils.SteamAPIResponse = generateGameMap(steamResponse)
	}

	// Get a list of games that Steam classifies as installed
	libraries, err := findDataFolders()
	if err != nil {
		return nil, err
	}
	// Get a list of all the folders in each location
	var storage []string
	for i := range libraries {
		// Get the paths to crawl
		var compat, shader string
		if libraries[i].Path == SteamDataRoot {
			compat = SteamCompatRoot
			shader = SteamShaderRoot
		} else {
			compat = filepath.Join(libraries[i].Path, ExternalCompatRoot)
			shader = filepath.Join(libraries[i].Path, ExternalShaderRoot)
		}

		// Crawl the compat path and append the folders
		// Append if no error, to prevent crashing for users that haven't synced data first.
		dir, _ := getDirectoryList(compat, true)
		if err == nil {
			storage = append(storage, dir...)
		}

		// Crawl the shader path and append the folders
		// Append if no error, to prevent crashing for users that haven't synced data first.
		dir, _ = getDirectoryList(shader, true)
		if err == nil {
			storage = append(storage, dir...)
		}
	}

	// Store all games present on the filesystem, in all compat and shader directories
	localGames := make(map[int]GameStatus)

	for i := range storage {
		intGame, _ := strconv.Atoi(storage[i])
		localGames[intGame] = GameStatus{
			GameName:    CryoUtils.SteamAPIResponse[intGame],
			IsInstalled: false,
		}
	}

	// Loop through each library location's installed games
	for i := range libraries {
		for j := range libraries[i].InstalledGames {
			// If the game has an entry in the list of localGames, then update the isInstalled flag to true for that game.
			val, keyExists := localGames[libraries[i].InstalledGames[j]]
			if keyExists {
				val.IsInstalled = true
				localGames[libraries[i].InstalledGames[j]] = val
			}
		}
	}

	var sortedMap []int

	for i := range localGames {
		// Add each key to the sortedMap slice, so we can sort it afterwards.
		sortedMap = append(sortedMap, i)
	}
	// Sort the slice
	sort.Ints(sortedMap)

	// For each entry in the completed list, add an entry to the check group to return
	for key := range sortedMap {
		// Skips non-game prefixes
		if sortedMap[key] == 0 || sortedMap[key] >= SteamGameMaxInteger {
			continue
		}

		var optionStr string
		var gameStr string

		// If the game name is known, use that, otherwise ???.
		if localGames[sortedMap[key]].GameName != "" {
			gameStr = localGames[sortedMap[key]].GameName
		} else {
			gameStr = "???"
		}

		if localGames[sortedMap[key]].IsInstalled {
			optionStr = fmt.Sprintf("%d - %s - Installed", sortedMap[key], gameStr)
		} else {
			optionStr = fmt.Sprintf("%d - %s - Not Installed", sortedMap[key], gameStr)
		}
		cleanupList.Append(optionStr)
	}

	return cleanupList, nil
}

// Get data to move values as canvas elements.
func getDataToMoveUI(data DataToMove) (*widget.List, *widget.List, error) {
	var leftList, rightList *widget.List

	// Use the cached API Response if already present
	if CryoUtils.SteamAPIResponse == nil {
		// Make a map of all games stored in the steam API
		steamResponse, err := querySteamAPI()
		if err != nil {
			CryoUtils.ErrorLog.Println(err)
		}
		CryoUtils.SteamAPIResponse = generateGameMap(steamResponse)
	}

	// Get lists of data to move
	leftList = widget.NewList(
		func() int {
			return len(data.right)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Left Side")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			gameInt, _ := strconv.Atoi(data.right[i])
			gameName := CryoUtils.SteamAPIResponse[gameInt]

			o.(*widget.Label).SetText(data.right[i] + " - " + gameName)
		})

	rightList = widget.NewList(
		func() int {
			return len(data.left)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Right Side")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			gameInt, _ := strconv.Atoi(data.left[i])
			gameName := CryoUtils.SteamAPIResponse[gameInt]

			o.(*widget.Label).SetText(data.left[i] + " - " + gameName)
		})

	return leftList, rightList, nil
}

func (app *Config) refreshSwapContent() {
	app.InfoLog.Println("Refreshing Swap data...")
	swap, err := getSwapFileSize()
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		swapStr := "Current Swap Size: Unknown"
		app.SwapText.Text = swapStr
		app.SwapText.Color = Gray
	} else {
		humanSwapSize := swap / int64(GigabyteMultiplier)
		swapStr := fmt.Sprintf("Current Swap Size: %dGB", humanSwapSize)
		app.SwapText.Text = swapStr
		if swap == RecommendedSwapSizeBytes {
			app.SwapText.Color = Green
		} else {
			app.SwapText.Color = Red
		}
	}

	app.SwapText.Refresh()
}

func (app *Config) refreshSwappinessContent() {
	app.InfoLog.Println("Refreshing Swappiness data...")
	swappiness, err := getSwappinessValue()
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		swappinessStr := "Current Swappiness: Unknown"
		app.SwappinessText.Text = swappinessStr
		app.SwappinessText.Color = Gray
	} else {
		swappinessStr := fmt.Sprintf("Current Swappiness: %d", swappiness)
		app.SwappinessText.Text = swappinessStr
		if strconv.Itoa(swappiness) == RecommendedSwappiness {
			app.SwappinessText.Color = Green
		} else {
			app.SwappinessText.Color = Red
		}
	}

	app.SwappinessText.Refresh()
}

func (app *Config) refreshHugePagesContent() {
	app.InfoLog.Println("Refreshing HugePages data...")
	if getHugePagesStatus() {
		app.HugePagesButton.Text = "Disable HugePages"
		app.HugePagesText.Color = Green
	} else {
		app.HugePagesButton.Text = "Enable HugePages"
		app.HugePagesText.Color = Red
	}
	app.HugePagesButton.Refresh()
	app.HugePagesText.Refresh()
}

func (app *Config) refreshShMemContent() {
	app.InfoLog.Println("Refreshing shmem data...")
	if getShMemStatus() {
		app.ShMemButton.Text = "Disable Shared Memory in THP"
		app.ShMemText.Color = Green
	} else {
		app.ShMemButton.Text = "Enable Shared Memory in THP"
		app.ShMemText.Color = Red
	}
	app.ShMemButton.Refresh()
	app.ShMemText.Refresh()

}

func (app *Config) refreshCompactionProactivenessContent() {
	app.InfoLog.Println("Refreshing compaction proactiveness data...")
	if getCompactionProactivenessStatus() {
		app.CompactionProactivenessButton.Text = "Revert Compaction Proactiveness"
		app.CompactionProactivenessText.Color = Green
	} else {
		app.CompactionProactivenessButton.Text = "Set Compaction Proactiveness"
		app.CompactionProactivenessText.Color = Red
	}
	app.CompactionProactivenessButton.Refresh()
	app.CompactionProactivenessText.Refresh()
}

func (app *Config) refreshDefragContent() {
	app.InfoLog.Println("Refreshing defragmentation data...")
	if getDefragStatus() {
		app.DefragButton.Text = "Enable Huge Page Defragmentation"
		app.DefragText.Color = Green
	} else {
		app.DefragButton.Text = "Disable Huge Page Defragmentation"
		app.DefragText.Color = Red
	}
	app.DefragButton.Refresh()
	app.DefragText.Refresh()
}

func (app *Config) refreshPageLockUnfairnessContent() {
	app.InfoLog.Println("Refreshing page lock unfairness data...")
	if getPageLockUnfairnessStatus() {
		app.PageLockUnfairnessButton.Text = "Revert Page Lock Unfairness"
		app.PageLockUnfairnessText.Color = Green
	} else {
		app.PageLockUnfairnessButton.Text = "Set Page Lock Unfairness"
		app.PageLockUnfairnessText.Color = Red
	}
	app.PageLockUnfairnessButton.Refresh()
	app.PageLockUnfairnessText.Refresh()
}

func (app *Config) refreshPageClusterContent() {
	app.InfoLog.Println("Refreshing page cluster data...")
	if getPageClusterStatus() {
		app.PageClusterButton.Text = "Revert Page Cluster"
		app.PageClusterText.Color = Green
	} else {
		app.PageClusterButton.Text = "Set Page Cluster"
		app.PageClusterText.Color = Red
	}
	app.PageClusterButton.Refresh()
	app.PageClusterText.Refresh()
}

func (app *Config) refreshAllContent() {
	app.refreshSwapContent()
	app.refreshSwappinessContent()
	app.refreshHugePagesContent()
	app.refreshCompactionProactivenessContent()
	app.refreshShMemContent()
	app.refreshDefragContent()
	app.refreshPageLockUnfairnessContent()
}
