package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	cp "github.com/otiai10/copy"
)

type StorageStatus struct {
	LeftCompatDirectories  []string
	LeftShaderDirectories  []string
	RightCompatDirectories []string
	RightShaderDirectories []string
}

type DataToMove struct {
	right     []string
	left      []string
	rightSize int64
	leftSize  int64
}

// Get a list of the directories inside the provided directory, ignoring symbolic links
func getDirectoryList(path string, includeSymlinks bool) ([]string, error) {
	var folderList []string
	// Get a list of all files in directory
	files, err := os.ReadDir(path)
	if err != nil {
		CryoUtils.ErrorLog.Println("Unable to list files in", path)
		return nil, err
	}

	for _, file := range files {
		// ui.CryoUtils.InfoLog.Println(file.Name())
		// Create a full path with the file and path names
		fullPath := filepath.Join(path, file.Name())
		if !includeSymlinks {
			// If the file is a directory AND is NOT a symlink, append the name of the folder to the list
			if file.IsDir() && !isSymbolicLink(fullPath) {
				folderList = append(folderList, file.Name())
			}
		} else {
			if file.IsDir() {
				folderList = append(folderList, file.Name())
			}
		}
	}

	return folderList, nil
}

// Update instance of StorageStatus with the current listings of directories.
func (s *StorageStatus) getStorageStatus(left string, right string) error {
	var err error
	if left == SteamDataRoot {
		s.LeftCompatDirectories, err = getDirectoryList(SteamCompatRoot, false)
		if err != nil {
			return err
		}
		s.LeftShaderDirectories, err = getDirectoryList(SteamShaderRoot, false)
		if err != nil {
			return err
		}
	} else {
		compat := filepath.Join(left, ExternalCompatRoot)
		shader := filepath.Join(left, ExternalShaderRoot)
		// Create the directories if they don't exist already
		_ = os.MkdirAll(compat, 0777)
		_ = os.MkdirAll(shader, 0777)
		s.LeftCompatDirectories, err = getDirectoryList(compat, false)
		if err != nil {
			return err
		}
		s.LeftShaderDirectories, err = getDirectoryList(shader, false)
		if err != nil {
			return err
		}
	}

	if right == SteamDataRoot {
		s.RightCompatDirectories, err = getDirectoryList(SteamCompatRoot, false)
		if err != nil {
			return err
		}
		s.RightShaderDirectories, err = getDirectoryList(SteamShaderRoot, false)
		if err != nil {
			return err
		}
	} else {
		compat := filepath.Join(right, ExternalCompatRoot)
		shader := filepath.Join(right, ExternalShaderRoot)
		// Create the directories if they don't exist already
		_ = os.MkdirAll(compat, 0777)
		_ = os.MkdirAll(shader, 0777)
		s.RightCompatDirectories, err = getDirectoryList(compat, false)
		if err != nil {
			return err
		}
		s.RightShaderDirectories, err = getDirectoryList(shader, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DataToMove) getSpaceNeeded(left string, right string) {
	var leftCompat, rightCompat, leftShader, rightShader string
	if left == SteamDataRoot {
		leftCompat = SteamCompatRoot
		leftShader = SteamShaderRoot
	} else {
		leftCompat = filepath.Join(left, ExternalCompatRoot)
		leftShader = filepath.Join(left, ExternalShaderRoot)
	}

	if right == SteamDataRoot {
		rightCompat = SteamCompatRoot
		rightShader = SteamShaderRoot
	} else {
		rightCompat = filepath.Join(right, ExternalCompatRoot)
		rightShader = filepath.Join(right, ExternalShaderRoot)
	}

	for x := range d.left {
		d.leftSize += getDirectorySize(filepath.Join(leftCompat, d.left[x]))
		d.leftSize += getDirectorySize(filepath.Join(leftShader, d.left[x]))
	}

	for x := range d.right {
		d.rightSize += getDirectorySize(filepath.Join(rightCompat, d.right[x]))
		d.rightSize += getDirectorySize(filepath.Join(rightShader, d.right[x]))
	}
}

// Populate a DataToMove object with the current queue of data needing to be moved.
func (d *DataToMove) getDataToMove(left string, right string) error {
	libraries, err := findDataFolders()
	if err != nil {
		return err
	}
	storage := new(StorageStatus)
	err = storage.getStorageStatus(left, right)
	if err != nil {
		return err
	}

	for i := range libraries {
		if isSubPath(left, libraries[i].Path) {
			CryoUtils.InfoLog.Println("Library location selected as left:", libraries[i].Path)
			// If the library is in the left-side parent directory
			for _, game := range libraries[i].InstalledGames {
				gameString := strconv.Itoa(game)
				CryoUtils.InfoLog.Println("Library contains:", gameString)
				// If game has files on the right side
				if contains(storage.RightCompatDirectories, gameString) && contains(storage.RightShaderDirectories, gameString) {
					d.right = append(d.right, gameString)
				}
			}
		} else if isSubPath(right, libraries[i].Path) {
			CryoUtils.InfoLog.Println("Library location selected as right:", libraries[i].Path)
			// If the library is in the right-side parent directory
			for _, game := range libraries[i].InstalledGames {
				gameString := strconv.Itoa(game)
				CryoUtils.InfoLog.Println("Library contains:", gameString)
				// If game is installed on the right side AND has files on the left side
				if contains(storage.LeftCompatDirectories, gameString) && contains(storage.LeftShaderDirectories, gameString) {
					d.left = append(d.left, gameString)
				}
			}
		} else {
			CryoUtils.InfoLog.Println("Library location not selected, skipping:", libraries[i].Path)
		}
	}
	return nil
}

// Move game data between each location as necessary
func moveGameData(data DataToMove, left string, right string) error {
	var progressPerMove = 1.0 / float64(len(data.right)+len(data.left))
	var leftCompatPath, leftShaderPath, rightCompatPath, rightShaderPath string

	if left == SteamDataRoot {
		leftCompatPath = SteamCompatRoot
		leftShaderPath = SteamShaderRoot
	} else {
		leftCompatPath = filepath.Join(left, ExternalCompatRoot)
		leftShaderPath = filepath.Join(left, ExternalShaderRoot)
	}

	if right == SteamDataRoot {
		rightCompatPath = SteamCompatRoot
		rightShaderPath = SteamShaderRoot
	} else {
		rightCompatPath = filepath.Join(right, ExternalCompatRoot)
		rightShaderPath = filepath.Join(right, ExternalShaderRoot)
	}

	// Moving to the left
	for _, directory := range data.right {
		leftCompatDir := filepath.Join(leftCompatPath, directory)
		leftShaderDir := filepath.Join(leftShaderPath, directory)
		rightCompatDir := filepath.Join(rightCompatPath, directory)
		rightShaderDir := filepath.Join(rightShaderPath, directory)

		// Remove any symlinks on the SSD in preparation for either moving to the SSD, or creating new symlinks
		steamCompatDir := filepath.Join(SteamCompatRoot, directory)
		if isSymbolicLink(steamCompatDir) {
			_ = os.Remove(steamCompatDir)
		}
		steamShaderDir := filepath.Join(SteamShaderRoot, directory)
		if isSymbolicLink(steamShaderDir) {
			_ = os.Remove(steamShaderDir)
		}

		// Copy the files
		CryoUtils.InfoLog.Println("Moving " + directory + " left...")
		err := cp.Copy(rightCompatDir, leftCompatDir)
		if err != nil {
			CryoUtils.ErrorLog.Println(err)
			return err
		}
		err = cp.Copy(rightShaderDir, leftShaderDir)
		if err != nil {
			CryoUtils.ErrorLog.Println(err)
			return err
		}

		// Remove the old files on the right
		CryoUtils.InfoLog.Println("Removing old " + rightCompatDir)
		err = os.RemoveAll(rightCompatDir)
		if err != nil {
			CryoUtils.ErrorLog.Println(err)
			return err
		}
		waitForDeletion(rightCompatPath, directory)
		CryoUtils.InfoLog.Println("Removing old " + rightShaderDir)
		err = os.RemoveAll(rightShaderDir)
		if err != nil {
			CryoUtils.ErrorLog.Println(err)
			return err
		}
		waitForDeletion(rightShaderPath, directory)

		// If the destination is NOT on the SSD, make symlinks
		if leftCompatPath != SteamCompatRoot {
			// Create symlinks on the SSD to the new location
			CryoUtils.InfoLog.Println("Creating symlink to new path on SSD...")
			err = os.Symlink(leftCompatDir, filepath.Join(SteamCompatRoot, directory))
			if err != nil {
				CryoUtils.ErrorLog.Println(err)
				return err
			}
			err = os.Symlink(leftShaderDir, filepath.Join(SteamShaderRoot, directory))
			if err != nil {
				CryoUtils.ErrorLog.Println(err)
				return err
			}
		}

		CryoUtils.MoveDataProgressBar.SetValue(CryoUtils.MoveDataProgressBar.Value + progressPerMove)
	}

	// Moving to the right
	for _, directory := range data.left {
		leftCompatDir := filepath.Join(leftCompatPath, directory)
		leftShaderDir := filepath.Join(leftShaderPath, directory)
		rightCompatDir := filepath.Join(rightCompatPath, directory)
		rightShaderDir := filepath.Join(rightShaderPath, directory)

		// Remove any symlinks on the SSD in preparation for either moving to the SSD, or creating new symlinks
		steamCompatDir := filepath.Join(SteamCompatRoot, directory)
		if isSymbolicLink(steamCompatDir) {
			_ = os.Remove(steamCompatDir)
		}
		steamShaderDir := filepath.Join(SteamShaderRoot, directory)
		if isSymbolicLink(steamShaderDir) {
			_ = os.Remove(steamShaderDir)
		}

		// Copy the files
		CryoUtils.InfoLog.Println("Moving " + directory + " right...")
		err := cp.Copy(leftCompatDir, rightCompatDir)
		if err != nil {
			CryoUtils.ErrorLog.Println(err)
			return err
		}
		err = cp.Copy(leftShaderDir, rightShaderDir)
		if err != nil {
			CryoUtils.ErrorLog.Println(err)
			return err
		}

		// Remove the old files on the left
		CryoUtils.InfoLog.Println("Removing old " + leftCompatDir)
		err = os.RemoveAll(leftCompatDir)
		if err != nil {
			CryoUtils.ErrorLog.Println(err)
			return err
		}
		waitForDeletion(leftCompatPath, directory)
		CryoUtils.InfoLog.Println("Removing old " + leftShaderDir)
		err = os.RemoveAll(leftShaderDir)
		if err != nil {
			CryoUtils.ErrorLog.Println(err)
			return err
		}
		waitForDeletion(leftShaderPath, directory)

		// If the destination is NOT on the SSD, make symlinks
		if rightCompatPath != SteamCompatRoot {
			// Create symlinks on the SSD to the new location
			CryoUtils.InfoLog.Println("Creating symlink to new path on SSD...")
			err = os.Symlink(rightCompatDir, filepath.Join(SteamCompatRoot, directory))
			if err != nil {
				CryoUtils.ErrorLog.Println(err)
				return err
			}
			err = os.Symlink(rightShaderDir, filepath.Join(SteamShaderRoot, directory))
			if err != nil {
				CryoUtils.ErrorLog.Println(err)
				return err
			}
		}
		CryoUtils.MoveDataProgressBar.SetValue(CryoUtils.MoveDataProgressBar.Value + progressPerMove)
	}
	return nil
}

// Confirm that all directories are in the proper locations post-move.
func (d *DataToMove) confirmDirectoryStatus(left string, right string) (bool, error) {
	var unmoved []string
	var dirs StorageStatus
	_ = dirs.getStorageStatus(left, right)

	for _, directory := range d.right {
		for _, x := range dirs.RightCompatDirectories {
			if x == directory {
				unmoved = append(unmoved, directory)
			}
		}
		for _, x := range dirs.RightShaderDirectories {
			if x == directory {
				unmoved = append(unmoved, directory)
			}
		}
	}

	for _, directory := range d.left {
		for _, x := range dirs.LeftCompatDirectories {
			if x == directory {
				unmoved = append(unmoved, directory)
			}
		}
		for _, x := range dirs.LeftShaderDirectories {
			if x == directory {
				unmoved = append(unmoved, directory)
			}
		}
	}

	if len(unmoved) == 0 {
		return true, nil
	} else {
		return false, fmt.Errorf("the following directories remain in the incorrect locations:\n"+
			"%s", unmoved)
	}
}

func getUninstalledGamesData() (uninstalled []string) {

	localGames, err := getLocalGameList()
	if err != nil {
		return nil
	}

	for key, game := range localGames {
		if key != 0 && key <= SteamGameMaxInteger && !game.IsInstalled {
			uninstalled = append(uninstalled, strconv.Itoa(key))
		}
	}

	return uninstalled

}
