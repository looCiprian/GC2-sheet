package C2

import "fmt"

var ErrorUnableToPullCommandAndTicker = fmt.Errorf(
	"an error occurred while pulling command and ticker from remote source",
)
var ErrorUnableToPushCommand = fmt.Errorf(
	"an error occurred while pushing command to remote source",
)

type CommandExecutor interface {
	pullCommandAndTicker(index int) (string, int, error)
	pushOutput(index int, output string) error
}

const DefaultTickerDuration = 10
