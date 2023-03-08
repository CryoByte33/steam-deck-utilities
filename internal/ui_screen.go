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

/*
#cgo LDFLAGS: -lX11
#include <X11/Xlib.h>
*/
import "C"
import (
	"fmt"
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

type ScreenSizer struct {
	screen screen
}

type screen struct {
	name          string
	width, height int
}

func NewScreenSizer() *ScreenSizer {
	return &ScreenSizer{screen: getDefaultScreen()}
}

func getDefaultScreen() screen {
	display := C.XOpenDisplay(nil)
	if display == nil {
		return screen{name: defaultScreenName, width: defaultSteamDeckScreenWidth, height: defaultSteamDeckScreenHeight}
	}
	defer C.XCloseDisplay(display)

	displayScreen := C.XDefaultScreenOfDisplay(display)
	displayScreenName := C.GoString(C.XDisplayString(display))
	width := int(C.XWidthOfScreen(displayScreen))
	height := int(C.XHeightOfScreen(displayScreen))
	if width == 0 || height == 0 {
		return screen{name: defaultScreenName, width: defaultSteamDeckScreenWidth, height: defaultSteamDeckScreenHeight}
	}
	return screen{displayScreenName, width, height}
}

func (s ScreenSizer) UpdateScaleForActiveMonitor() {
	defaultScreen := getDefaultScreen()

	if defaultScreen.width <= smallScreen {
		s.setScale(0.25)
	} else if defaultScreen.width > smallScreen && defaultScreen.width <= mediumScreen {
		s.setScale(1)
	} else if defaultScreen.width > mediumScreen {
		s.setScale(2)
	}
}

func (s ScreenSizer) setScale(f float32) {
	err := os.Setenv(scaleEnvKey, fmt.Sprintf("%f", f))
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
	}
}
