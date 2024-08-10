package C2

import (
	"GC2-sheet/internal/configuration"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type MicrosoftToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	ExpireTime   int64
}

func (token *MicrosoftToken) isExpired() bool {
	var now = time.Now().Unix()
	return token.ExpireTime <= now
}

var ErrorUnableToGetToken = fmt.Errorf("an error occurred while getting a Microsoft token")

func newMicrosoftToken() (*MicrosoftToken, error) {
	tenantID := configuration.GetOptionsMicrosoftTenantID()
	clientId := configuration.GetOptionsMicrosoftClientID()
	clientSecret := configuration.GetOptionsMicrosoftClientSecret()

	scope := "https://graph.microsoft.com/.default"
	grantType := "client_credentials"

	URL := "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0/token"

	body := url.Values{}
	body.Add("client_id", clientId)
	body.Add("scope", scope)
	body.Add("client_secret", clientSecret)
	body.Add("grant_type", grantType)

	resp, err := http.Post(URL, "application/x-www-form-urlencoded", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorUnableToGetToken, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrorUnableToGetToken
	}

	var token *MicrosoftToken
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorUnableToGetToken, err)
	}

	// subtracting 3 seconds for buffer
	token.ExpireTime = time.Now().Unix() + int64(token.ExpiresIn) - 3

	return token, nil
}

type MicrosoftConnector struct {
	client *http.Client
	token  *MicrosoftToken
}

var ErrorUnableToCreateMicrosoftConnector = fmt.Errorf("an error occurred while creating Microsoft connector")

func newMicrosoftConnector() (*MicrosoftConnector, error) {
	connector := &MicrosoftConnector{
		client: newMicrosoftHttpClient(),
		token:  nil,
	}

	_, err := connector.getValidToken()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrorUnableToCreateMicrosoftConnector, err)
	}

	return connector, nil
}

func (connector *MicrosoftConnector) getValidToken() (*MicrosoftToken, error) {
	var err error
	if connector.token == nil || connector.token.isExpired() {
		connector.token, err = newMicrosoftToken()
	}

	if err != nil {
		return nil, err
	}

	return connector.token, nil
}

func (connector *MicrosoftConnector) do(request *http.Request) (*http.Response, error) {
	token, err := connector.getValidToken()
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", token.AccessToken)

	return connector.client.Do(request)
}

func newMicrosoftHttpClient() *http.Client {

	proxyUrl := configuration.GetOptionsProxy()

	transport := &http.Transport{}

	if proxyUrl != nil {
		transport.Proxy = http.ProxyURL(proxyUrl)
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{
		Transport: transport,
	}

	return client
}
