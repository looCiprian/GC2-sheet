package C2

import "fmt"

var ErrorUnableToPullCommandAndTicker = fmt.Errorf(
	"an error occurred while pulling command and ticker from remote source",
)
var ErrorUnableToPushCommand = fmt.Errorf(
	"an error occurred while pushing command to remote source",
)

type CommandExecutor interface {
	pullCommandAndTicker() (string, int, error)
	pushOutput(*Command) error
	getLastCommand() *Command
	appendEmptyCommand()
}

type Command struct {
	Ticker int    // Ticker delay for polling
	RowId  int    // row number starting from 1
	Input  string // command to execute
	Output string // command output
}

const DefaultTickerDuration = 10

func NewCommand(rowId int) *Command {
	return &Command{
		Ticker: DefaultTickerDuration,
		RowId:  rowId,
		Input:  "",
		Output: "",
	}
}
