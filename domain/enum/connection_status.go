package enum

type ConnectionStatus int

const (
	CONNECTED ConnectionStatus = iota
	CONNECTION_REQUEST
	BLOCKED
	NONE
)
