package internal

import (
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io"
	"os/exec"
)

// Renews sudo auth for GUI mode
func renewSudoAuth() {
	// Do a really basic command to renew sudo auth
	cmd := exec.Command("sudo", "-S", "--", "echo")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
	}

	go func() {
		defer func(stdin io.WriteCloser) {
			err := stdin.Close()
			if err != nil {
				CryoUtils.ErrorLog.Println(err)
				return
			}
		}(stdin)
		_, err := io.WriteString(stdin, CryoUtils.UserPassword)
		if err != nil {
			CryoUtils.ErrorLog.Println(err)
			return
		}
	}()

	_, err = cmd.CombinedOutput()
	if err != nil {
		CryoUtils.ErrorLog.Println(err)
	}
}

// Test that the sudo password works
func testAuth(password string) error {
	progress := widget.NewProgressBarInfinite()
	d := dialog.NewCustom("Testing Authentication", "Quit", progress,
		CryoUtils.MainWindow,
	)
	d.Show()
	// Do a really basic command to renew sudo auth
	cmd := exec.Command("sudo", "-S", "--", "echo")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		d.Hide()
		return err
	}

	go func() {
		defer func(stdin io.WriteCloser) {
			err := stdin.Close()
			if err != nil {
				d.Hide()
				CryoUtils.ErrorLog.Println(err)
			}
		}(stdin)
		_, err := io.WriteString(stdin, password)
		if err != nil {
			d.Hide()
			CryoUtils.ErrorLog.Println(err)
		}
	}()

	_, err = cmd.CombinedOutput()
	if err != nil {
		d.Hide()
		return err
	}

	d.Hide()

	return nil
}
