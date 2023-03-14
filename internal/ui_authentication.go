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
	"os/exec"
	"time"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Renews sudo auth for GUI mode
func renewSudoAuth() {
	// Do a really basic command to renew sudo auth
	cmd := exec.Command("sudo", "-S", "--", "echo")
	//Sudo will exit immediately if it's the correct password, but will hang for a moment if it isn't.
	cmd.WaitDelay = 500 * time.Millisecond
	stdin, err := cmd.StdinPipe()
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		return
	}
	err = cmd.Start()
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		return
	}
	_, err = stdin.Write([]byte(CryoUtils.UserPassword + "\n"))
	if err != nil {
		cmd.Process.Kill()
		CryoUtils.ErrorLog.Println(err)
		return
	}
	stdin.Close()
	err = cmd.Wait()
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
		return
	}
}

// Test that the sudo password works
func testAuth(password string) error {
	progress := widget.NewProgressBarInfinite()
	d := dialog.NewCustom("Testing Authentication", "Quit", progress,
		CryoUtils.MainWindow,
	)
	d.Show()
	defer d.Hide()
	// Do a really basic command to renew sudo auth
	cmd := exec.Command("sudo", "-S", "--", "echo")
	//Sudo will exit immediately if it's the correct password, but will hang for a moment if it isn't.
	cmd.WaitDelay = 500 * time.Millisecond
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	_, err = stdin.Write([]byte(password + "\n"))
	if err != nil {
		cmd.Process.Kill()
		return err
	}
	stdin.Close()
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
