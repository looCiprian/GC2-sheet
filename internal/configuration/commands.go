package configuration

type SpreadSheet struct {
	DriveId       string
	SpreadSheetId string
	CommandSheet  Sheet
}
type Sheet struct {
	Name              string // sheet name
	CommandsExecution []Commands
}

type Commands struct {
	RangeIn  string // Example !A
	RangeOut string
	RangeId  int // Example 1
	Input    string
	Output   string
}
