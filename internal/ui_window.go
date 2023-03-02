package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func syncGameDataWindow() {
	var selectionContainer *fyne.Container
	// Create a new window
	w := CryoUtils.App.NewWindow("Sync Game Data")

	driveList, err := getListOfAttachedDrives()
	if err != nil {
		presentErrorInUI(err, w)
	}

	if len(driveList) > 1 {
		// Place a prompt near the top of the window
		prompt := canvas.NewText("Please select the devices you'd like to sync data between.", nil)
		prompt.TextSize, prompt.TextStyle = 18, fyne.TextStyle{Bold: true}

		// Make the list widgets with initial contents, excluding what the other has pre-selected
		// This is simply to make it more "one click" for users.
		leftList := widget.NewSelect(removeElementFromStringSlice(driveList[1], driveList), func(s string) {})
		rightList := widget.NewSelect(removeElementFromStringSlice(driveList[0], driveList), func(s string) {})

		// Pre-define each direction with the default values
		leftSelected := removeElementFromStringSlice(driveList[1], driveList)[0]
		leftList.Selected = leftSelected
		rightSelected := removeElementFromStringSlice(driveList[0], driveList)[0]
		rightList.Selected = rightSelected
		// Define the OnChanged functions after, so both are aware of the other's existence.
		leftList.OnChanged = func(s string) {
			leftSelected = s
			// Remove the selected option from the other list to prevent both sides being the same.
			rightList.Options = removeElementFromStringSlice(s, driveList)
		}
		rightList.OnChanged = func(s string) {
			rightSelected = s
			// Remove the selected option from the other list to prevent both sides being the same.
			leftList.Options = removeElementFromStringSlice(s, driveList)
		}

		cancelButton := widget.NewButton("Cancel", func() {
			w.Close()
		})
		submitButton := widget.NewButton("Submit", func() {
			selectionContainer.Hide()
			w.CenterOnScreen()
			populateGameDataWindow(w, leftSelected, rightSelected)
		})
		buttonBar := container.NewHSplit(cancelButton, submitButton)
		selectionContainer = container.NewVBox(prompt, leftList, rightList, buttonBar)
	} else {
		// Place a prompt near the top of the window
		prompt := canvas.NewText("Not enough drives attached to sync data", Red)
		prompt.TextSize, prompt.TextStyle = 18, fyne.TextStyle{Bold: true}

		cancelButton := widget.NewButton("Cancel", func() {
			w.Close()
		})

		selectionContainer = container.NewVBox(prompt, cancelButton)
	}

	w.SetContent(selectionContainer)
	w.CenterOnScreen()
	w.RequestFocus()
	w.Show()
}

func populateGameDataWindow(w fyne.Window, left string, right string) {
	var leftCard, rightCard *widget.Card
	var syncDataButton *widget.Button
	var data DataToMove

	p := widget.NewProgressBarInfinite()
	d := dialog.NewCustom("Finding data to move...", "Dismiss", p, w)
	d.Show()

	// Get a list of data to move
	err := data.getDataToMove(left, right)
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		d.Hide()
		presentErrorInUI(err, w)
	}

	// Get the storage totals necessary for each side
	data.getSpaceNeeded(left, right)

	leftSpaceAvailable, err := getFreeSpace(left)
	if err != nil {
		presentErrorInUI(err, w)
	}
	rightSpaceAvailable, err := getFreeSpace(right)
	if err != nil {
		presentErrorInUI(err, w)
	}

	// Place a prompt near the top of the window
	prompt := canvas.NewText("Please confirm that it is okay to move this data:", nil)
	prompt.TextSize, prompt.TextStyle = 18, fyne.TextStyle{Bold: true}

	// User-presentable strings
	leftDataStr := fmt.Sprintf("Data to be moved to %s", right)
	leftSizeStr := fmt.Sprintf("Total Size: %.2fGB", float64(data.leftSize)/float64(GigabyteMultiplier))
	rightDataStr := fmt.Sprintf("Data to be moved to %s", left)
	rightSizeStr := fmt.Sprintf("Total Size: %.2fGB", float64(data.rightSize)/float64(GigabyteMultiplier))

	// Deal with lack of space on left
	if leftSpaceAvailable < data.rightSize {
		leftCard = widget.NewCard(leftDataStr, "",
			canvas.NewText("Error: Not enough space available on destination drive.", Red))
		// Provide a button to close the window
		syncDataButton = widget.NewButton("Close", func() {
			w.Close()
		})
	}

	// Deal with lack of space on right
	if rightSpaceAvailable < data.leftSize {
		rightCard = widget.NewCard(rightDataStr, "",
			canvas.NewText("Error: Not enough space available on destination drive.", Red))
		// Provide a button to close the window
		syncDataButton = widget.NewButton("Close", func() {
			w.Close()
		})
	}

	leftList, rightList, err := getDataToMoveUI(data)
	// Deal with error
	if err != nil {
		// Create an error in each card if directories can't be listed
		leftCard = widget.NewCard(leftDataStr, "",
			canvas.NewText("Error: Failed to get list of directories to be moved.", Red))
		rightCard = widget.NewCard(rightDataStr, "",
			canvas.NewText("Error: Failed to get list of directories to be moved.", Red))
		// Provide a button to close the window
		syncDataButton = widget.NewButton("Close", func() {
			w.Close()
		})
	}

	// If there's anything to move left
	if len(data.right) != 0 {
		leftCard = widget.NewCard(rightDataStr, rightSizeStr, leftList)
	} else {
		leftCard = widget.NewCard(rightDataStr, "",
			canvas.NewText("None! Everything is synced in this direction.", Green))
	}

	// If there's anything to move right
	if len(data.left) != 0 {
		rightCard = widget.NewCard(leftDataStr, leftSizeStr, rightList)
	} else {
		rightCard = widget.NewCard(leftDataStr, "",
			canvas.NewText("None! Everything is synced in this direction.", Green))
	}

	// Create button if something can sync
	if len(data.right) != 0 || len(data.left) != 0 {
		syncDataButton = widget.NewButton("Confirm", func() {
			// Do the actual sync
			CryoUtils.InfoLog.Println("Sync data confirmed")
			progress := widget.NewProgressBar()
			CryoUtils.MoveDataProgressBar = progress
			progress.Resize(fyne.NewSize(500, 50))
			tempVBox := container.NewVBox(canvas.NewText("Moving items, please wait...", nil), progress)
			widget.ShowModalPopUp(tempVBox, w.Canvas())
			err = moveGameData(data, left, right)
			if err != nil {
				presentErrorInUI(err, w)
			} else {
				_, err := data.confirmDirectoryStatus(left, right)
				if err != nil {
					presentErrorInUI(err, w)
				} else {
					CryoUtils.InfoLog.Println("All data moved properly, printing success!")
					dialog.ShowInformation(
						"Success!",
						"Data move completed, all game data is synced to the appropriate device.",
						CryoUtils.MainWindow,
					)
					w.Close()
				}
			}
		})
	} else {
		// Otherwise, provide a button to close the window
		syncDataButton = widget.NewButton("Close", func() {
			w.Close()
		})
	}
	cancelButton := widget.NewButton("Cancel", func() {
		w.Close()
	})

	d.Hide()

	// Format the window
	syncMain := container.NewGridWithColumns(1, leftCard, rightCard)
	syncButtonBorder := container.NewGridWithColumns(2, cancelButton, syncDataButton)
	syncLayout := container.NewBorder(nil, syncButtonBorder, nil, nil, syncMain)
	w.SetContent(syncLayout)
	w.Resize(fyne.NewSize(300, 450))
	w.CenterOnScreen()
	w.RequestFocus()
	w.Show()
}

func cleanupDataWindow() {
	var cleanupCard *widget.Card
	var cleanupButton, cancelButton *widget.Button

	// Create a new window
	w := CryoUtils.App.NewWindow("Clean Game Data")

	var removeList []string
	cleanupList, err := createGameDataList()
	if err != nil {
		presentErrorInUI(err, CryoUtils.MainWindow)
	}
	cleanupList.OnChanged = func(s []string) {
		var tempList []string

		for i := range s {
			// Get only the game ID for the selected games
			tempList = append(tempList, strings.Split(s[i], " ")[0])
		}
		removeList = tempList
	}
	cleanupScroll := container.NewVScroll(cleanupList)

	// Create an error in each card if directories can't be listed
	cleanupCard = widget.NewCard("Clean Stale Game Data",
		"Choose which game's prefixes and shadercache you would like to remove.",
		cleanupScroll)
	cancelButton = widget.NewButton("Cancel", func() {
		w.Close()
	})
	cleanupButton = widget.NewButton("Delete Selected", func() {
		dialog.ShowConfirm("Are you sure?", "Are you sure you want to delete these files?\n\n"+
			"Please be sure to back up any Non-Steam-Cloud save games before\n"+
			"deleting them using this tool, as any selected will be lost.",
			func(b bool) {
				if b {
					possibleLocations, err := getListOfDataAllDataLocations()
					if err != nil {
						CryoUtils.ErrorLog.Println(err)
						presentErrorInUI(err, CryoUtils.MainWindow)
					}

					removeGameData(removeList, possibleLocations)

					dialog.ShowInformation(
						"Success!",
						"Process completed!",
						CryoUtils.MainWindow,
					)
					w.Close()
				} else {
					w.Close()
				}
			}, w)
	})

	cleanAllUninstalled := widget.NewButton("Delete All Uninstalled", func() {
		dialog.ShowConfirm("Are you sure?", "Are you sure you want to delete these files?\n\n"+
			"Please be sure to back up any Non-Steam-Cloud save games before\n"+
			"deleting them using this tool, as any selected will be lost.",
			func(b bool) {
				if !b {
					w.Close()
				}

				locations, err := getListOfDataAllDataLocations()
				if err != nil {
					CryoUtils.ErrorLog.Println(err)
					presentErrorInUI(err, CryoUtils.MainWindow)
				}

				removeGameData(getUninstalledGamesData(), locations)

				dialog.ShowInformation(
					"Success!",
					"Process completed!",
					CryoUtils.MainWindow,
				)
				w.Close()

			}, w)

	})

	// Format the window
	cleanupMain := container.NewGridWithColumns(1, cleanupCard)
	cleanupButtonsGrid := container.NewGridWithColumns(2, cancelButton, cleanupButton)
	extraButtonGrid := container.NewGridWithColumns(1, cleanAllUninstalled)
	footerButtons := container.NewGridWithColumns(1, cleanupButtonsGrid, extraButtonGrid)
	cleanupLayout := container.NewBorder(nil, footerButtons, nil, nil, cleanupMain)
	w.SetContent(cleanupLayout)
	w.Resize(fyne.NewSize(300, 450))
	w.CenterOnScreen()
	w.RequestFocus()
	w.Show()
}

func getUninstalledGamesData() (uninstalled []string) {

	localGames, err := getLocalGameList()
	if err != nil {
		return nil
	}

	for key, game := range localGames {
		if game.IsInstalled == false {
			uninstalled = append(uninstalled, strconv.Itoa(key))
		}
	}

	return uninstalled

}

func removeGameData(removeList []string, locations []string) {

	CryoUtils.InfoLog.Println("Removing the following content:")
	for i := range removeList {
		for j := range locations {
			path := filepath.Join(locations[j], removeList[i])
			CryoUtils.InfoLog.Println(path)
			err := os.RemoveAll(path)
			if err != nil {
				CryoUtils.ErrorLog.Println(err)
				presentErrorInUI(err, CryoUtils.MainWindow)
			}
		}
	}
}

func swapSizeWindow() {
	// Create a new window
	w := CryoUtils.App.NewWindow("Change Swap Size")

	// Place a prompt near the top of the window
	prompt := canvas.NewText("Please choose the new swap file size in gigabytes:", nil)
	prompt.TextSize, prompt.TextStyle = 18, fyne.TextStyle{Bold: true}

	// Determine maximum available space for a swap file and construct a list of available sizes based on it
	availableSwapSizes, err := getAvailableSwapSizes()
	if err != nil {
		presentErrorInUI(err, w)
	}

	// Give the user a choice in swap file sizes
	var chosenSize int
	choice := widget.NewRadioGroup(availableSwapSizes, func(value string) {
		// Only grab the number at the beginning of the string, allows for suffixes.
		chosenSize, err = strconv.Atoi(strings.Split(value, " ")[0])
		if err != nil {
			presentErrorInUI(err, w)
		}
	})

	// Provide a button to submit the choice
	swapResizeButton := widget.NewButton("Resize Swap File", func() {
		progress := widget.NewProgressBarInfinite()
		d := dialog.NewCustom("Resizing Swap File, please be patient..."+
			"(This can take up to 30 minutes)", "Quit", progress,
			w,
		)
		d.Show()
		err = changeSwapSizeGUI(chosenSize)
		if err != nil {
			d.Hide()
			presentErrorInUI(err, w)
		} else {
			d.Hide()
			dialog.ShowInformation(
				"Success!",
				"Process completed! You can verify the file is resized by\n"+
					"running 'ls -lash /home/swapfile' or 'swapon -s' in Konsole.",
				CryoUtils.MainWindow,
			)
			CryoUtils.refreshSwapContent()
			w.Close()
		}
	})

	// Make a progress bar and hide it
	progress := widget.NewProgressBar()
	CryoUtils.SwapResizeProgressBar = progress
	progress.Hide()

	// Format the window
	swapVBox := container.NewVBox(prompt, choice, swapResizeButton)
	w.SetContent(swapVBox)
	w.Resize(fyne.NewSize(400, 300))
	w.CenterOnScreen()
	w.RequestFocus()
	w.Show()
}

// Note: Having a separate function for this is hacky, but necessary for progress bar functionality
func changeSwapSizeGUI(size int) error {
	// Disable swap temporarily
	renewSudoAuth()
	CryoUtils.InfoLog.Println("Disabling swap temporarily...")
	err := disableSwap()
	if err != nil {
		return err
	}
	// Resize the file
	renewSudoAuth()
	err = resizeSwapFile(size)
	if err != nil {
		return err
	}
	// Set permissions on file
	renewSudoAuth()
	err = setSwapPermissions()
	if err != nil {
		return err
	}
	// Initialize new swap file
	renewSudoAuth()
	err = initNewSwapFile()
	if err != nil {
		return err
	}
	return nil
}

func swappinessWindow() {
	// Create a new window
	w := CryoUtils.App.NewWindow("Change Swappiness")

	// Place a prompt near the top of the window
	prompt := canvas.NewText("Please choose the new swappiness.", nil)
	prompt.TextSize, prompt.TextStyle = 18, fyne.TextStyle{Bold: true}

	// Give the user a choice in swap file sizes
	var chosenSwappiness string
	choice := widget.NewRadioGroup(AvailableSwappinessOptions, func(value string) {
		chosenSwappiness = strings.Fields(value)[0]
	})

	// Provide a button to submit the choice
	swappinessChangeButton := widget.NewButton("Change Swappiness", func() {
		renewSudoAuth()
		err := ChangeSwappiness(chosenSwappiness)
		if err != nil {
			presentErrorInUI(err, w)
		} else {
			dialog.ShowInformation(
				"Success!",
				"Swappiness change completed!",
				CryoUtils.MainWindow,
			)
			CryoUtils.refreshSwappinessContent()
			w.Close()
		}
	})

	// Format the window
	swapVBox := container.NewVBox(prompt, choice, swappinessChangeButton)
	w.SetContent(swapVBox)
	w.CenterOnScreen()
	w.RequestFocus()
	w.Show()
}
