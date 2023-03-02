package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/moby/sys/mountinfo"
	"golang.org/x/sys/unix"
)

type Config struct {
	App                           fyne.App
	InfoLog                       *log.Logger
	ErrorLog                      *log.Logger
	SwapText                      *canvas.Text
	SwappinessText                *canvas.Text
	HugePagesText                 *canvas.Text
	ShMemText                     *canvas.Text
	CompactionProactivenessText   *canvas.Text
	DefragText                    *canvas.Text
	PageLockUnfairnessText        *canvas.Text
	VRAMText                      *canvas.Text
	SteamAPIResponse              map[int]string
	MainWindow                    fyne.Window
	SwapResizeProgressBar         *widget.ProgressBar
	MoveDataProgressBar           *widget.ProgressBar
	HomeContainer                 *fyne.Container
	GameDataContainer             *fyne.Container
	MemoryContainer               *fyne.Container
	SwapBar                       *fyne.Container
	MemoryBar                     *fyne.Container
	HugePagesButton               *widget.Button
	ShMemButton                   *widget.Button
	CompactionProactivenessButton *widget.Button
	DefragButton                  *widget.Button
	PageLockUnfairnessButton      *widget.Button
	VRAMButton                    *widget.Button
	UserPassword                  string
	SwapFileLocation              string
}

var CryoUtils Config

var stat unix.Statfs_t

func getFreeSpace(path string) (int64, error) {
	err := unix.Statfs(path, &stat)
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		return 0, fmt.Errorf("error getting free space")
	}
	return int64(stat.Bfree * uint64(stat.Bsize)), nil
}

func getDirectorySize(path string) int64 {
	var size int64
	_ = filepath.Walk(path, func(_ string, info os.FileInfo, _ error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

func isSymbolicLink(path string) bool {
	fi, err := os.Lstat(path)
	if err != nil {
		CryoUtils.ErrorLog.Println("Unable to determine if file was symlink:", path)
		panic(err)
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		return true
	}
	return false
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func waitForDeletion(path string, directory string) {
	for {
		if !doesDirectoryExist(path, directory) {
			break
		}
		time.Sleep(time.Second)
	}
}

// Checks the path variable until directory is no longer found, then exits.
func doesDirectoryExist(path string, directory string) bool {
	directories, _ := os.ReadDir(path)
	for _, dir := range directories {
		if dir.Name() == directory {
			return true
		}
	}
	return false
}

func doesFileExist(path string) bool {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func isSubPath(parent string, sub string) bool {
	subFolds := filepath.SplitList(sub)
	for i, fold := range filepath.SplitList(parent) {
		if subFolds[i] != fold {
			return false
		}
	}
	return true
}

// Write a file with a given string
func writeFile(path string, contents string) error {
	CryoUtils.InfoLog.Println("Writing", path)

	tempPath := filepath.Join(InstallDirectory, "temp.txt")
	// Try to remove tempfile just in case it exists for some reason
	_ = removeFile(tempPath)
	// Write to the CU install directory to avoid permissions issues
	f, err := os.Create(tempPath)
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		return err
	}

	defer f.Close()

	_, err = f.WriteString(contents)
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		return err
	}

	// Move the completed file to final location.
	_, err = exec.Command("sudo", "mv", tempPath, path).Output()
	if err != nil {
		return fmt.Errorf("error moving temp file to final location")
	}

	return nil
}

func removeFile(path string) error {
	CryoUtils.InfoLog.Println("Removing", path)
	_, err := exec.Command("sudo", "rm", path).Output()
	if err != nil {
		CryoUtils.ErrorLog.Println("Couldn't delete", path, ", likely missing.")
	}
	return nil
}

// getListOfDataAllDataLocations Get a list off all data locations (compat and shader data).
func getListOfDataAllDataLocations() ([]string, error) {
	drives, err := getListOfAttachedDrives()
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		return nil, err
	}

	var possibleLocations []string
	for x := range drives {
		if drives[x] == SteamDataRoot {
			possibleLocations = append(possibleLocations, SteamCompatRoot)
			possibleLocations = append(possibleLocations, SteamDataRoot)
		} else {
			compat := filepath.Join(drives[x], ExternalCompatRoot)
			shader := filepath.Join(drives[x], ExternalShaderRoot)
			possibleLocations = append(possibleLocations, compat)
			possibleLocations = append(possibleLocations, shader)
		}
	}

	return possibleLocations, nil
}

func getListOfAttachedDrives() ([]string, error) {
	drives := []string{SteamDataRoot}
	filter := mountinfo.PrefixFilter(MountDirectory)
	info, err := mountinfo.GetMounts(filter)
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		return nil, err
	}
	for x := range info {
		drives = append(drives, info[x].Mountpoint)
	}
	CryoUtils.InfoLog.Printf("Attached drives: %s", drives)
	return drives, nil
}

func removeElementFromStringSlice(str string, slice []string) []string {
	var newSlice []string
	for x := range slice {
		if str != slice[x] {
			newSlice = append(newSlice, slice[x])
		}
	}
	return newSlice
}

func getUnitStatus(param string) (string, error) {
	var output string
	cmd, err := exec.Command("sudo", "cat", UnitMatrix[param]).Output()
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		return "nil", err
	}
	// This is just to get the actual value in units which present as a list.
	if strings.Contains(string(cmd), "[") {
		slice := strings.Fields(string(cmd))
		for x := range slice {
			if strings.Contains(slice[x], "[") {
				output = strings.ReplaceAll(slice[x], "[", "")
				output = strings.ReplaceAll(output, "]", "")
			}
		}
	} else {
		output = strings.TrimSpace(string(cmd))
	}
	return output, nil
}

func writeUnitFile(param string, value string) error {
	path := filepath.Join(TmpFilesRoot, param+".conf")
	CryoUtils.InfoLog.Println("Writing", value, "to", path, "to preserve", param, "setting...")
	contents := strings.ReplaceAll(TemplateUnitFile, "PARAM", UnitMatrix[param])
	contents = strings.ReplaceAll(contents, "VALUE", value)
	err := writeFile(path, contents)
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		return err
	}
	return nil
}

func removeUnitFile(param string) error {
	path := filepath.Join(TmpFilesRoot, param+".conf")
	CryoUtils.InfoLog.Println("Removing", path, "to revert", param, "setting...")
	err := removeFile(path)
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		return err
	}
	return nil
}

func setUnitValue(param string, value string) error {
	CryoUtils.InfoLog.Println("Writing", value, "for param", param, "to memory.")
	// This mess is the only way I could find to push directly to unit files, without requiring
	// a sudo password on installation to change capabilities.
	echoCmd := exec.Command("echo", value)
	teeCmd := exec.Command("sudo", "tee", UnitMatrix[param])
	reader, writer := io.Pipe()
	var buf bytes.Buffer
	echoCmd.Stdout = writer
	teeCmd.Stdin = reader
	teeCmd.Stdout = &buf
	echoCmd.Start()
	teeCmd.Start()
	echoCmd.Wait()
	writer.Close()
	teeCmd.Wait()
	reader.Close()
	io.Copy(os.Stdout, &buf)

	return nil
}

func getHumanVRAMSize(size int) string {
	// Converts the VRAM size to human readable format.
	// The size argument is in MB.
	text := fmt.Sprintf("%dMB", size)
	if(size >= 1024) {
		text = fmt.Sprintf("%dGB", size / 1024)
	}

	return text
}