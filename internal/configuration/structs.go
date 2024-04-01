package configuration

// Sheet must be implemented for every types e.g. Google, Microsoft
type Sheet struct {
	Name                     string // sheet name
	CommandsExecution        []*Commands
	RangeTickerConfiguration string
}

const DefaultTickerDuration = 10

type Commands struct {
	RangeIn  string // cell A
	RangeOut string // cell B
	RangeLog string // cell C
	RangeId  int    // row number starting from 1
	Input    string // command to execute
	Output   string // command output
}
