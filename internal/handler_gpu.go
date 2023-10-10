package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Get the current VRAM
func getVRAMValue() (int, error) {
	cmd, err := createCommand("glxinfo", "-B").Output()

	// Extract video memory
	re := regexp.MustCompile(`Video memory: [0-9]+`)
	match := re.FindStringSubmatch(string(cmd))

	if err != nil || match == nil {
		return 100, fmt.Errorf("error getting current VRAM")
	}

	output := strings.Split(match[0], " ")[2]
	CryoUtils.InfoLog.Println("Found a VRAM of", output)
	vram, _ := strconv.Atoi(output)

	return vram, nil
}
