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
