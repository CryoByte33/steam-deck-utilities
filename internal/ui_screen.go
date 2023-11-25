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
	smallScreenSize  = 100
	mediumScreenSize = 200
)

const scaleEnvKey = "FYNE_SCALE"
const defaultScreenName = ":0"

type ScreenSizer struct {
	screen screen
}

type screen struct {
	name                          string
	widthPhysical, heightPhysical int
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
	return primaryScreen
}

func (s ScreenSizer) UpdateScaleForActiveMonitor() {

	if s.screen.widthPhysical <= smallScreenSize {
		s.setScale(0.25)
	} else if s.screen.widthPhysical > smallScreenSize && s.screen.widthPhysical <= mediumScreenSize {
		s.setScale(0.75)
	} else if s.screen.widthPhysical > mediumScreenSize {
		s.setScale(1)
	}
}

func (s ScreenSizer) setScale(f float32) {
	err := os.Setenv(scaleEnvKey, fmt.Sprintf("%f", f))
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
	}
}

func defaultScreen() screen {
	return screen{name: defaultScreenName, widthPhysical: 60, heightPhysical: 100}
}
