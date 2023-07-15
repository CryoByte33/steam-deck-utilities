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
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Home tab for "recommended" and "default" buttons
func (app *Config) homeTab() *fyne.Container {
	welcomeText := canvas.NewText("Welcome to CryoUtilities!", White)
	welcomeText.TextSize = HeaderTextSize
	welcomeText.TextStyle.Bold = true

	subheadingText := canvas.NewText("Quick settings. Use the tabs at the top of the window to use "+
		"settings individually.", White)
	subheadingText.TextSize = SubHeadingTextSize

	availableSpace, err := getFreeSpace("/home")
	if err != nil {
		presentErrorInUI(err, app.MainWindow)
	}
	var chosenSize string
	if availableSpace < RecommendedSwapSizeBytes {
		availableSizes, _ := getAvailableSwapSizes()
		chosenSize = strings.Fields(availableSizes[len(availableSizes)-1])[0]
	} else {
		chosenSize = strconv.Itoa(RecommendedSwapSize)
	}

	actionText := widget.NewLabel(
		"Swap: " + chosenSize + "GB\n" +
			"Swappiness: " + RecommendedSwappiness + "\n" +
			"HugePages: Enabled\n" +
			"Compaction Proactivenes: " + RecommendedCompactionProactiveness + "\n" +
			"HugePage Defragmentation: Disabled\n" +
			"Page Lock Unfairness: " + RecommendedPageLockUnfairness + "\n" +
			"Shared Memory in Huge Pages: Enabled")

	recommendedButton := widget.NewButton("Recommended", func() {
		progressGroup := container.NewVBox(
			canvas.NewText("Applying recommended settings...", White),
			actionText,
			widget.NewProgressBarInfinite())
		modal := widget.NewModalPopUp(progressGroup, CryoUtils.MainWindow.Canvas())
		modal.Show()
		renewSudoAuth()
		err := UseRecommendedSettings()
		if err != nil {
			presentErrorInUI(err, CryoUtils.MainWindow)
		}
		modal.Hide()
		app.refreshAllContent()
		dialog.ShowInformation(
			"Success!",
			"Recommended settings applied!",
			CryoUtils.MainWindow,
		)
	})
	stockButton := widget.NewButton("Stock", func() {
		progressText := canvas.NewText("Reverting to stock settings...", White)
		progressBar := widget.NewProgressBarInfinite()
		progressGroup := container.NewVBox(progressText, progressBar)
		modal := widget.NewModalPopUp(progressGroup, CryoUtils.MainWindow.Canvas())
		modal.Show()
		renewSudoAuth()
		err := UseStockSettings()
		if err != nil {
			presentErrorInUI(err, CryoUtils.MainWindow)
		}
		modal.Hide()
		app.refreshAllContent()
		dialog.ShowInformation(
			"Success!",
			"Stock settings applied!",
			CryoUtils.MainWindow,
		)
	})

	recommendedSettings := widget.NewCard("Recommended Settings", "Set all settings to "+
		"CryoByte33's recommendations.", recommendedButton)
	stockSettings := widget.NewCard("Stock Settings", "Reset all settings to Valve defaults, excludes "+
		"'Game Data' tab/locations.", stockButton)

	homeVBox := container.NewVBox(
		welcomeText,
		subheadingText,
		recommendedSettings,
		stockSettings,
	)
	app.HomeContainer = homeVBox

	return homeVBox
}

// Swap tab for all swap-related tasks.
func (app *Config) swapTab() *fyne.Container {
	app.SwapText = canvas.NewText("Swap File Size: Unknown", Gray)
	app.SwappinessText = canvas.NewText("Swappiness: Unknown", Gray)
	// Main content including buttons to resize swap and change swappiness
	swapResizeButton := widget.NewButton("Resize", func() {
		swapSizeWindow()
		app.refreshSwapContent()
	})
	swappinessChangeButton := widget.NewButton("Change", func() {
		swappinessWindow()
		app.refreshSwappinessContent()
	})

	swapCard := widget.NewCard("Swap File", "Resize the swap file.", swapResizeButton)
	swappinessCard := widget.NewCard("Swappiness", "Change the swappiness value.", swappinessChangeButton)

	// Swap info gathering
	app.refreshSwapContent()
	app.refreshSwappinessContent()

	app.SwapBar = container.NewGridWithColumns(2,
		container.NewCenter(app.SwapText),
		container.NewCenter(app.SwappinessText))

	topBar := container.NewVBox(
		container.NewGridWithRows(1),
		container.NewGridWithRows(1, container.NewCenter(canvas.NewText("Current Swap Status:", White))),
		app.SwapBar,
	)

	swapVBox := container.NewVBox(
		swapCard,
		swappinessCard,
	)

	full := container.NewBorder(topBar, nil, nil, nil, swapVBox)

	return full
}

// Game Data tab to move and delete prefixes and shadercache.
func (app *Config) storageTab() *fyne.Container {
	// These can take a minute to come up, so create a loading bar to show things are happening.
	syncDataButton := widget.NewButton("Sync", func() {
		progressText := canvas.NewText("Calculating device status...", White)
		progressBar := widget.NewProgressBarInfinite()
		progressGroup := container.NewVBox(progressText, progressBar)
		modal := widget.NewModalPopUp(progressGroup, CryoUtils.MainWindow.Canvas())
		modal.Show()
		syncGameDataWindow()
		modal.Hide()
	})
	cleanupDataButton := widget.NewButton("Clean", func() {
		progressText := canvas.NewText("Calculating device status...", White)
		progressBar := widget.NewProgressBarInfinite()
		progressGroup := container.NewVBox(progressText, progressBar)
		modal := widget.NewModalPopUp(progressGroup, CryoUtils.MainWindow.Canvas())
		modal.Show()
		cleanupDataWindow()
		modal.Hide()
	})

	syncData := widget.NewCard("Sync Game Data", "Sync prefix and shaders to the device where the game "+
		"is installed", syncDataButton)
	cleanStaleData := widget.NewCard("Delete Game Data", "Delete prefixes and shaders for selected games.",
		cleanupDataButton)

	gameDataVBox := container.NewVBox(
		syncData,
		cleanStaleData,
	)
	app.GameDataContainer = gameDataVBox

	return gameDataVBox
}

// Tab for non-swap, memory-related tweaks.
func (app *Config) memoryTab() *fyne.Container {
	app.HugePagesText = canvas.NewText("Huge Pages (THP)", Red)
	app.ShMemText = canvas.NewText("Shared Memory in THP", Red)
	app.CompactionProactivenessText = canvas.NewText("Compaction Proactiveness", Red)
	app.DefragText = canvas.NewText("Defrag", Red)
	app.PageLockUnfairnessText = canvas.NewText("Page Lock Unfairness", Red)

	CryoUtils.HugePagesButton = widget.NewButton("Enable HugePages", func() {
		renewSudoAuth()
		err := ToggleHugePages()
		if err != nil {
			presentErrorInUI(err, CryoUtils.MainWindow)
		}
		app.refreshHugePagesContent()
	})

	CryoUtils.ShMemButton = widget.NewButton("Enable Shared Memory in THP", func() {
		renewSudoAuth()
		err := ToggleShMem()
		if err != nil {
			presentErrorInUI(err, CryoUtils.MainWindow)
		}
		app.refreshShMemContent()
	})

	CryoUtils.CompactionProactivenessButton = widget.NewButton("Set Compaction Proactiveness", func() {
		renewSudoAuth()
		err := ToggleCompactionProactiveness()
		if err != nil {
			presentErrorInUI(err, CryoUtils.MainWindow)
		}
		app.refreshCompactionProactivenessContent()
	})

	CryoUtils.DefragButton = widget.NewButton("Disable Huge Page Defragmentation", func() {
		renewSudoAuth()
		err := ToggleDefrag()
		if err != nil {
			presentErrorInUI(err, CryoUtils.MainWindow)
		}
		app.refreshDefragContent()
	})

	CryoUtils.PageLockUnfairnessButton = widget.NewButton("Set Page Lock Unfairness", func() {
		renewSudoAuth()
		err := TogglePageLockUnfairness()
		if err != nil {
			presentErrorInUI(err, CryoUtils.MainWindow)
		}
		app.refreshPageLockUnfairnessContent()
	})

	app.refreshHugePagesContent()
	app.refreshCompactionProactivenessContent()
	app.refreshShMemContent()
	app.refreshDefragContent()
	app.refreshPageLockUnfairnessContent()

	app.MemoryBar = container.NewGridWithColumns(5,
		container.NewCenter(app.HugePagesText),
		container.NewCenter(app.ShMemText),
		container.NewCenter(app.CompactionProactivenessText),
		container.NewCenter(app.DefragText),
		container.NewCenter(app.PageLockUnfairnessText))
	topBar := container.NewVBox(
		container.NewGridWithRows(1),
		container.NewGridWithRows(1, container.NewCenter(canvas.NewText("Current Tweak Status:", White))),
		app.MemoryBar,
	)

	hugePagesCard := widget.NewCard("Huge Pages", "Toggle huge pages", app.HugePagesButton)
	shMemCard := widget.NewCard("Shared Memory in THP", "Toggle shared memory in THP", app.ShMemButton)
	compactionProactivenessCard := widget.NewCard("Compaction Proactiveness", "Set compaction proactiveness", app.CompactionProactivenessButton)
	defragCard := widget.NewCard("Huge Page Defragmentation", "Toggle huge page defragmentation", app.DefragButton)
	pageLockUnfairnessCard := widget.NewCard("Page Lock Unfairness", "Set page lock unfairness", app.PageLockUnfairnessButton)

	memoryVBox := container.NewVBox(
		hugePagesCard,
		shMemCard,
		compactionProactivenessCard,
		defragCard,
		pageLockUnfairnessCard,
	)
	scroll := container.NewScroll(memoryVBox)
	full := container.NewBorder(topBar, nil, nil, nil, scroll)

	return full
}

func (app *Config) vramTab() *fyne.Container {
	app.VRAMText = canvas.NewText("Current VRAM size: Unknown", Gray)

	// Get VRAM value
	app.refreshVRAMContent()

	textHowTo := widget.NewLabel("1. Turn off the Steam Deck\n\n" +
		"2. Press and hold the volume up button, press the power button, then release both\n\n" +
		"3. Navigate to Setup Utility -> Advanced -> UMA Frame Buffer Size")

	textRecommended := widget.NewLabelWithStyle("4G is the recommended setting for most situations", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	textWarning := widget.NewLabel("Please be aware that some games (RDR2) may experience degraded performance.")

	textVBox := container.NewVBox(
		textHowTo,
		textRecommended,
		textWarning,
	)

	vramCard := widget.NewCard("Minimum VRAM", "How to change the minimum VRAM:", textVBox)

	vramBAR := container.NewGridWithColumns(1,
		container.NewCenter(app.VRAMText))
	topBar := container.NewVBox(
		container.NewGridWithRows(1),
		container.NewGridWithRows(1, container.NewCenter(canvas.NewText("Current Tweak Status:", White))),
		vramBAR,
	)

	vramVBOX := container.NewVBox(
		vramCard,
	)
	scroll := container.NewScroll(vramVBOX)
	full := container.NewBorder(topBar, nil, nil, nil, scroll)

	return full
}
