package C2

import (
	"GC2-sheet/internal/configuration"
	"context"
	"crypto/tls"
	"golang.org/x/oauth2"
	"net/http"

	"google.golang.org/api/option"
	GHTTP "google.golang.org/api/transport/http"
)

type GoogleConnector struct {
	client *http.Client
}

func newGoogleHttpClient() (context.Context, *http.Client) {
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

	client := &http.Client{Transport: myTransport}

	return ctx, client
}

func NewGoogleConnector() (*GoogleConnector, error) {
	_, client := newGoogleHttpClient()

	connector := &GoogleConnector{
		client: client,
	}

	return connector, nil
}
