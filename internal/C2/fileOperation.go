package C2

import (
	"os"
)

func safeFile(path string, fileContent []byte) error {

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = f.Write(fileContent)
	if err != nil {
		return err
	}

	return nil
}
