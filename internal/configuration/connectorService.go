package configuration

type ConnectorService int8

const (
	Google ConnectorService = iota
	Microsoft
)

var ConnectorName = map[ConnectorService]string{
	Google:    "Google",
	Microsoft: "Microsoft",
}

func (c ConnectorService) String() string {
	return ConnectorName[c]
}
