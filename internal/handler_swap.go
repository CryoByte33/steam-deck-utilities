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
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/afero"
)

type logger interface {
	Println(v ...any)
}

// NewSwap is a constructor for Swap.
func NewSwap(
	defaultSwapSizeBytes int64,
	availableSwapSizes []string,
	oldSwappinessUnitFile string,
	defaultSwapFileLocation string,
	fs afero.Fs,
	loggerInfo logger,
) (*Swap, error) {
	if defaultSwapSizeBytes == 0 {
		return nil, errors.New("defaultSwapSizeBytes is required")
	}
	if len(availableSwapSizes) == 0 {
		return nil, errors.New("availableSwapSizes is required")
	}
	if oldSwappinessUnitFile == "" {
		return nil, errors.New("oldSwappinessUnitFile is required")
	}
	if loggerInfo == nil {
		return nil, errors.New("info logger is required")
	}
	if defaultSwapFileLocation == "" {
		return nil, errors.New("default swap location is required")
	}
	if fs == nil {
		return nil, errors.New("fs is required")
	}
	swapFileLocation, err := getSwapFileLocation(fs, defaultSwapFileLocation)
	if err != nil {
		return nil, fmt.Errorf("getting swapfile location: %w", err)
	}

	return &Swap{
		defaultSwapSizeBytes:  defaultSwapSizeBytes,
		availableSwapSizes:    availableSwapSizes,
		oldSwappinessUnitFile: oldSwappinessUnitFile,
		swapFileLocation:      swapFileLocation,
		fs:                    fs,
		loggerInfo:            loggerInfo,
	}, nil
}

// Swap decorates all functionality related to swap/swappiness changes.
type Swap struct {
	defaultSwapSizeBytes  int64
	availableSwapSizes    []string
	oldSwappinessUnitFile string
	swapFileLocation      string
	fs                    afero.Fs
	loggerInfo            logger
}

// Get swap file location from the system (/proc/swaps)
// Sample output:
// Filename				Type		Size	Used	Priority
// /home/swapfile			file		8388604	0	-2
func getSwapFileLocation(fs afero.Fs, defaultSwapFileLocation string) (string, error) {
	filepath := "/proc/swaps"
	file, err := fs.Open(filepath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	// skip the first line (header)
	scanner.Scan()

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 3 && fields[0] != "Filename" {
			location := fields[0]
			// If swapfile is a partition then return no swapfile found
			if strings.HasPrefix(location, "/dev/") {
				return "", fmt.Errorf("no swapfile found in %s", filepath)
			}
			return location, nil
		}
	}

	if doesFileExist(defaultSwapFileLocation) {
		return defaultSwapFileLocation, nil
	}

	return "", fmt.Errorf("no swapfile found in %s", filepath)
}

// Get the current swap and swappiness values
func (s *Swap) getSwappinessValue() (int, error) {
	cmd, err := exec.Command("sysctl", "vm.swappiness").Output()
	if err != nil {
		return 100, fmt.Errorf("error getting current swappiness")
	}
	output := strings.Fields(string(cmd))
	s.loggerInfo.Println("Found a swappiness of", output[2])
	swappiness, _ := strconv.Atoi(output[2])

	return swappiness, nil
}

// Get current swap file size, in bytes.
func (s *Swap) getSwapFileSize() (int64, error) {
	info, err := s.fs.Stat(s.swapFileLocation)
	if err != nil {
		// Don't crash the program, just report the default size
		return s.defaultSwapSizeBytes, fmt.Errorf("error getting current swap file size")
	}
	s.loggerInfo.Println("Found a swap file with a size of", info.Size())
	return info.Size(), nil
}

// Get the available space for a swap file and return a slice of strings
func (s *Swap) getAvailableSwapSizes() ([]string, error) {
	// Get the free space in /home
	currentSwapSize, _ := s.getSwapFileSize()
	availableSpace, err := getFreeSpace("/home")
	if err != nil {
		return nil, fmt.Errorf("error getting available space in /home: %w", err)
	}

	// Loop through the range of available sizes and create a list of viable options for the current Deck.
	// This will always leave 1 as an available option, just in case.
	validSizes := []string{"1 - Default"}
	for _, size := range s.availableSwapSizes {
		intSize, _ := strconv.Atoi(size)
		byteSize := intSize * GigabyteMultiplier
		if int64(byteSize+SpaceOverhead) < (availableSpace + currentSwapSize) {
			if byteSize == int(currentSwapSize) {
				currentSizeString := fmt.Sprintf("%s - Current Size", size)
				validSizes = append(validSizes, currentSizeString)
			} else {
				validSizes = append(validSizes, size)
			}
		}
	}

	s.loggerInfo.Println("Available Swap Sizes:", validSizes)
	return validSizes, nil
}

// Disable swapping completely
func (s *Swap) disableSwap() error {
	s.loggerInfo.Println("Disabling swap temporarily...")
	_, err := exec.Command("sudo", "swapoff", "-a").Output()
	if err != nil {
		return fmt.Errorf("error disabling swap: %w", err)
	}
	return err
}

// Resize the swap file to the provided size, in GB.
func (s *Swap) resizeSwapFile(size int) error {
	locationArg := fmt.Sprintf("of=%s", s.swapFileLocation)
	countArg := fmt.Sprintf("count=%d", size)

	s.loggerInfo.Println("Resizing swap to", size, "GB...")
	// Use dd to write zeroes, reevaluate using Go directly in the future
	_, err := exec.Command("sudo", "dd", "if=/dev/zero", locationArg, "bs=1G", countArg, "status=progress").Output()
	if err != nil {
		return fmt.Errorf("error resizing %s: %w", s.swapFileLocation, err)
	}
	return nil
}

// Set swap permissions to a valid value.
func (s *Swap) setSwapPermissions() error {
	s.loggerInfo.Println("Setting permissions on", s.swapFileLocation, "to 0600...")
	_, err := exec.Command("sudo", "chmod", "600", s.swapFileLocation).Output()
	if err != nil {
		return fmt.Errorf("error setting permissions on %s: %w", s.swapFileLocation, err)
	}
	return nil
}

// Enable swapping on the newly resized file.
func (s *Swap) initNewSwapFile() error {
	s.loggerInfo.Println("Enabling swap on", s.swapFileLocation, "...")
	_, err := exec.Command("sudo", "mkswap", s.swapFileLocation).Output()
	if err != nil {
		return fmt.Errorf("error creating swap on %s: %w", s.swapFileLocation, err)
	}
	_, err = exec.Command("sudo", "swapon", s.swapFileLocation).Output()
	if err != nil {
		return fmt.Errorf("error enabling swap on %s: %w", s.swapFileLocation, err)
	}
	return nil
}

// ChangeSwappiness Set swappiness to the provided integer.
func (s *Swap) ChangeSwappiness(value string) error {
	s.loggerInfo.Println("Setting swappiness...")
	// Remove old swappiness file while we're at it
	_ = removeFile(s.oldSwappinessUnitFile)
	err := setUnitValue("swappiness", value)
	if err != nil {
		return err
	}

	if value == DefaultSwappiness {
		s.loggerInfo.Println("Removing swappiness unit to revert to default behavior...")
		err = removeUnitFile("swappiness")
		if err != nil {
			return err
		}
	} else {
		err = writeUnitFile("swappiness", value)
		if err != nil {
			return err
		}
		return nil
	}

	// Return no error if everything went as planned
	return nil
}
