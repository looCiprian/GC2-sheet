package C2

import (
	"fmt"
	"os"
)

type FileSystem interface {
	// Download a file, returns file content and error if any
	pullFile(string) ([]byte, error)

	// Local path of the file to push, return error if any
	pushFile(name string, file *os.File) error
}

var ErrorUnableToPullFile = fmt.Errorf("an error occurred while pulling the file")
var ErrorUnableToPushFile = fmt.Errorf("an error occurred while pushing the file")
