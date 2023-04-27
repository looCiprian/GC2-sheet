package C2

import (
	"os"
)

func Exit() {

	path, _ := os.Executable()
	os.Remove(path)
	os.Exit(0)

}
