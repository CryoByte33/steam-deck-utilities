package internal

import (
	"fmt"
	"strconv"
)

// ChangeSwapSizeCLI Change the swap file size to the specified size in GB
func ChangeSwapSizeCLI(size int) error {
	// Get the free space in /home
	availableSpace, err := getFreeSpace("/home")
	if err != nil {
		return fmt.Errorf("error getting available space in /home")
	}

	selectedSize := size * GigabyteMultiplier
	if int64(selectedSize+SpaceOverhead) > availableSpace {
		return fmt.Errorf("not enough space to resize swap to chosen size")
	}

	// Disable swap temporarily
	err = disableSwap()
	if err != nil {
		return err
	}

	// Resize the file
	err = resizeSwapFile(size)
	if err != nil {
		return err
	}

	// Set permissions on file
	err = setSwapPermissions()
	if err != nil {
		return err
	}

	// Initialize new swap file
	err = initNewSwapFile()
	if err != nil {
		return err
	}
	return nil
}

func UseRecommendedSettings() error {
	// Change swap
	CryoUtils.InfoLog.Println("Starting swap file resize...")
	availableSpace, err := getFreeSpace("/home")
	if err != nil {
		return err
	}
	if availableSpace < RecommendedSwapSizeBytes {
		size := 1
		availableSizes, err := getAvailableSwapSizes()
		if err != nil {
			return err
		}
		if len(availableSizes) != 1 {
			// Get the last entry in the availableSizes list
			size, err = strconv.Atoi(availableSizes[len(availableSizes)-1])
			if err != nil {
				return err
			}
			// Never create a swap file larger than 16GB automatically.
			if size > 16 {
				size = 16
			}
		}
		err = ChangeSwapSizeCLI(size)
		if err != nil {
			return err
		}
	} else {
		err = ChangeSwapSizeCLI(RecommendedSwapSize)
		if err != nil {
			return err
		}
	}
	CryoUtils.InfoLog.Println("Swap file resized, changing swappiness...")
	err = ChangeSwappiness(RecommendedSwappiness)
	if err != nil {
		return err
	}

	CryoUtils.InfoLog.Println("Swappiness changed, enabling HugePages...")
	err = SetHugePages()
	if err != nil {
		return err
	}

	CryoUtils.InfoLog.Println("HugePages enabled, setting compaction proactiveness...")
	err = SetCompactionProactiveness()
	if err != nil {
		return err
	}

	CryoUtils.InfoLog.Println("Compaction proactiveness changed, disabling hugePage defragmentation...")
	err = SetDefrag()
	if err != nil {
		return err
	}

	CryoUtils.InfoLog.Println("HugePage defragmentation disabled, setting page lock unfairness...")
	err = SetPageLockUnfairness()
	if err != nil {
		return err
	}

	CryoUtils.InfoLog.Println("Page lock unfairness changed, enabling Shared Memory...")
	err = SetShMem()
	if err != nil {
		return err
	}

	CryoUtils.InfoLog.Println("All settings configured!")
	return nil
}

func UseStockSettings() error {
	CryoUtils.InfoLog.Println("Resizing swap file to 1GB...")
	// Revert swap file size
	err := ChangeSwapSizeCLI(DefaultSwapSize)
	if err != nil {
		return err
	}

	CryoUtils.InfoLog.Println("Setting swappiness to 100...")
	// Revert swappiness
	err = ChangeSwappiness(DefaultSwappiness)
	if err != nil {
		return err
	}

	CryoUtils.InfoLog.Println("Disabling HugePages...")
	// Enable HugePages
	err = RevertHugePages()
	if err != nil {
		return err
	}

	CryoUtils.InfoLog.Println("Reverting compaction proactiveness...")
	err = RevertCompactionProactiveness()
	if err != nil {
		return err
	}

	CryoUtils.InfoLog.Println("Enabling hugePage defragmentation...")
	err = RevertDefrag()
	if err != nil {
		return err
	}

	CryoUtils.InfoLog.Println("Reverting page lock unfairness...")
	err = RevertPageLockUnfairness()
	if err != nil {
		return err
	}

	CryoUtils.InfoLog.Println("Disabling shared memory in hugepages...")
	err = RevertShMem()
	if err != nil {
		CryoUtils.InfoLog.Println("All settings reverted to default!")
	}

	return nil
}
