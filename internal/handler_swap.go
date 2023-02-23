package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Get swap file location from the system (/proc/swaps)
// Sample output:
// Filename				Type		Size	Used	Priority
// /home/swapfile			file		8388604	0	-2
func getSwapFileLocation() (string, error) {
	file, err := os.Open("/proc/swaps")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// skip the first line (header)
	scanner.Scan()

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 3 && fields[0] != "Filename" {
			location := fields[0]
			// If swapfile is a partition then return no swapfile found
			if strings.HasPrefix(location, "/dev/") {
				return "", fmt.Errorf("no swapfile found")
			}
			return location, nil
		}
	}

	return "", fmt.Errorf("no swapfile found")
}

// Get the current swap and swappiness values
func getSwappinessValue() (int, error) {
	cmd, err := exec.Command("sysctl", "vm.swappiness").Output()
	if err != nil {
		return 100, fmt.Errorf("error getting current swappiness")
	}
	output := strings.Fields(string(cmd))
	CryoUtils.InfoLog.Println("Found a swappiness of", output[2])
	swappiness, _ := strconv.Atoi(output[2])

	return swappiness, nil
}

// Get current swap file size, in bytes.
func getSwapFileSize() (int64, error) {
	if doesFileExist(BTRFSSwapFileLocation) {
		CryoUtils.SwapFileLocation = BTRFSSwapFileLocation
	} else {
		CryoUtils.SwapFileLocation = DefaultSwapFileLocation
	}

	info, err := os.Stat(CryoUtils.SwapFileLocation)
	if err != nil {
		// Don't crash the program, just report the default size
		return DefaultSwapSizeBytes, fmt.Errorf("error getting current swap file size")
	}
	CryoUtils.InfoLog.Println("Found a swap file with a size of", info.Size())
	return info.Size(), nil
}

// Get the available space for a swap file and return a slice of strings
func getAvailableSwapSizes() ([]string, error) {
	// Get the free space in /home
	currentSwapSize, _ := getSwapFileSize()
	availableSpace, err := getFreeSpace("/home")
	if err != nil {
		return nil, fmt.Errorf("error getting available space in /home")
	}

	// Loop through the range of available sizes and create a list of viable options for the current Deck.
	// This will always leave 1 as an available option, just in case.
	validSizes := []string{"1 - Default"}
	for _, size := range AvailableSwapSizes {
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

	CryoUtils.InfoLog.Println("Available Swap Sizes:", validSizes)
	return validSizes, nil
}

// Disable swapping completely
func disableSwap() error {
	CryoUtils.InfoLog.Println("Disabling swap temporarily...")
	_, err := exec.Command("sudo", "swapoff", "-a").Output()
	if err != nil {
		return fmt.Errorf("error disabling swap")
	}
	return err
}

// Resize the swap file to the provided size, in GB.
func resizeSwapFile(size int) error {
	locationArg := fmt.Sprintf("of=%s", CryoUtils.SwapFileLocation)
	countArg := fmt.Sprintf("count=%d", size)

	CryoUtils.InfoLog.Println("Resizing swap to", size, "GB...")
	// Use dd to write zeroes, reevaluate using Go directly in the future
	_, err := exec.Command("sudo", "dd", "if=/dev/zero", locationArg, "bs=1G", countArg, "status=progress").Output()
	if err != nil {
		return fmt.Errorf("error resizing %s", CryoUtils.SwapFileLocation)
	}
	return nil
}

// Set swap permissions to a valid value.
func setSwapPermissions() error {
	CryoUtils.InfoLog.Println("Setting permissions on", CryoUtils.SwapFileLocation, "to 0600...")
	_, err := exec.Command("sudo", "chmod", "600", CryoUtils.SwapFileLocation).Output()
	if err != nil {
		return fmt.Errorf("error setting permissions on %s", CryoUtils.SwapFileLocation)
	}
	return nil
}

// Enable swapping on the newly resized file.
func initNewSwapFile() error {
	CryoUtils.InfoLog.Println("Enabling swap on", CryoUtils.SwapFileLocation, "...")
	_, err := exec.Command("sudo", "mkswap", CryoUtils.SwapFileLocation).Output()
	if err != nil {
		return fmt.Errorf("error creating swap on %s", DefaultSwapFileLocation)
	}
	_, err = exec.Command("sudo", "swapon", CryoUtils.SwapFileLocation).Output()
	if err != nil {
		return fmt.Errorf("error enabling swap on %s", DefaultSwapFileLocation)
	}
	return nil
}

// ChangeSwappiness Set swappiness to the provided integer.
func ChangeSwappiness(value string) error {
	CryoUtils.InfoLog.Println("Setting swappiness...")
	// Remove old swappiness file while we're at it
	_ = removeFile(OldSwappinessUnitFile)
	err := setUnitValue("swappiness", value)
	if err != nil {
		return err
	}

	if value == DefaultSwappiness {
		CryoUtils.InfoLog.Println("Removing swappiness unit to revert to default behavior...")
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
