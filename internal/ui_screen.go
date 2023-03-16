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
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
	"os"
)

const (
	smallScreen  = 1280
	mediumScreen = 1920
)

const (
	defaultSteamDeckScreenWidth  = 1280
	defaultSteamDeckScreenHeight = 800
)

const scaleEnvKey = "FYNE_SCALE"
const defaultScreenName = ":0"
const baselineDPI = 120.0

type ScreenSizer struct {
	screen screen
}

type screen struct {
	name                                         string
	width, height, widthPhysical, heightPhysical int
}

func NewScreenSizer() *ScreenSizer {
	return &ScreenSizer{screen: getDefaultScreen()}
}

func getDefaultScreen() screen {

	var primaryScreen screen

	err := glfw.Init()
	if err != nil {
		return defaultScreen()
	}
	defer glfw.Terminate()

	monitor := glfw.GetPrimaryMonitor()
	primaryScreen.name = monitor.GetName()
	primaryScreen.widthPhysical, primaryScreen.heightPhysical = monitor.GetPhysicalSize()
	primaryScreen.width = monitor.GetVideoMode().Width
	primaryScreen.height = monitor.GetVideoMode().Height
	if primaryScreen.width == 0 || primaryScreen.height == 0 {
		return defaultScreen()
	}
	return primaryScreen
}

func (s ScreenSizer) UpdateScaleForActiveMonitor() {

	fyneAlreadyOverrides := overrideFyneScale(s.screen)

	if fyneAlreadyOverrides {
		s.setScale(1)
	} else if s.screen.width <= smallScreen {
		s.setScale(0.25)
	} else if s.screen.width > smallScreen && s.screen.width <= mediumScreen {
		s.setScale(1)
	} else if s.screen.width > mediumScreen {
		s.setScale(2)
	}
}

func overrideFyneScale(defaultScreen screen) bool {
	dpi := float32(defaultScreen.width) / (float32(defaultScreen.widthPhysical) / 25.4)

	if dpi > 1000 || dpi < 10 {
		dpi = baselineDPI
	}

	scale := float32(float64(dpi) / baselineDPI)
	if scale < 1.0 {
		return true
	}
	return false
}

func (s ScreenSizer) setScale(f float32) {
	err := os.Setenv(scaleEnvKey, fmt.Sprintf("%f", f))
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
	}
}

func defaultScreen() screen {
	return screen{name: defaultScreenName, width: defaultSteamDeckScreenWidth, height: defaultSteamDeckScreenHeight}
}
