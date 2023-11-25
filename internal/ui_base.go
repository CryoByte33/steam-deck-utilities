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
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func InitUI(swap *Swap) {
	// Create a Fyne application
	screenSizer := NewScreenSizer()
	screenSizer.UpdateScaleForActiveMonitor()
	fyneApp := app.NewWithID("io.cryobyte.cryoutilities")
	CryoUtils.App = fyneApp
	CryoUtils.App.SetIcon(ResourceIconPng)

	// Show and run the app
	title := "CryoUtilities " + CurrentVersionNumber
	CryoUtils.MainWindow = fyneApp.NewWindow(title)
	CryoUtils.makeUI(swap)
	CryoUtils.MainWindow.CenterOnScreen()
	CryoUtils.MainWindow.ShowAndRun()
}

func (app *Config) makeUI(swap *Swap) {
	app.authUI(swap)

	// Show a disclaimer that I'm not responsible for damage.
	dialog.ShowConfirm("Disclaimer",
		"This script was made by CryoByte33 to resize the swapfile on a Steam Deck.\n\n"+
			"Disclaimer: I am in no way responsible to damage done to any\n"+
			"device this is executed on, all liability lies with the user.\n\n"+
			"Do you accept these terms?",
		func(b bool) {
			if !b {
				presentErrorInUI(errors.New("terms not accepted"), CryoUtils.MainWindow)
				CryoUtils.MainWindow.Close()
			} else {
				CryoUtils.InfoLog.Println("Terms accepted, continuing...")
			}
		},
		CryoUtils.MainWindow,
	)

	// Create and size a Fyne window
	CryoUtils.MainWindow.Resize(fyne.NewSize(700, 410))
	CryoUtils.MainWindow.SetFixedSize(true)
	CryoUtils.MainWindow.SetMaster()
}

func (app *Config) mainUI(swap *Swap) {
	// Create heading section
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), app.homeTab(swap)),
		container.NewTabItemWithIcon("Swap", theme.MailReplyAllIcon(), app.swapTab(swap)),
		container.NewTabItemWithIcon("Memory", theme.ComputerIcon(), app.memoryTab()),
		container.NewTabItemWithIcon("Storage", theme.StorageIcon(), app.storageTab()),
		container.NewTabItemWithIcon("VRAM", theme.ViewFullScreenIcon(), app.vramTab()),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	finalContent := container.NewVBox(tabs)
	app.MainWindow.SetContent(finalContent)
}

func (app *Config) authUI(swap *Swap) {
	// Refactor this, duplicated code.
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.OnSubmitted = func(s string) {
		CryoUtils.InfoLog.Println("Testing password...")
		err := testAuth(s)
		if err != nil {
			CryoUtils.InfoLog.Println("Password invalid, asking again...")
			dialog.ShowInformation("Incorrect password", "Incorrect password, please try again.",
				CryoUtils.MainWindow)
		} else {
			CryoUtils.InfoLog.Println("Password valid, continuing...")
			CryoUtils.UserPassword = s
			app.mainUI(swap)
		}
	}
	passwordButton := widget.NewButton("Submit", func() {
		CryoUtils.InfoLog.Println("Testing password...")
		err := testAuth(passwordEntry.Text)
		if err != nil {
			CryoUtils.InfoLog.Println("Password invalid, asking again...")
			dialog.ShowInformation("Incorrect password", "Incorrect password, please try again.",
				CryoUtils.MainWindow)
		} else {
			CryoUtils.InfoLog.Println("Password valid, continuing...")
			CryoUtils.UserPassword = passwordEntry.Text
			app.mainUI(swap)
		}
	})
	passwordVBox := container.NewVBox(passwordEntry, passwordButton)
	passwordContainer := widget.NewCard("Enter your sudo/deck password.", "Enter your sudo/deck password.", passwordVBox)

	//  Add container to window

	app.MainWindow.SetContent(passwordContainer)
}
