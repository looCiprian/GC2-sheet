package C2

import (
	"io"
	"os"

	"google.golang.org/api/drive/v3"
)

func downloadFile(clientDrive *drive.Service, fileId string, downloadPath string) error {

	resp, err := clientDrive.Files.Get(fileId).Download()
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	fileDownloaded, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	f, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	_, err = f.Write(fileDownloaded)
	if err != nil {
		return err
	}

	return nil
}
