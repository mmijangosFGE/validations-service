package constants

import "time"

// ConnectionState - enum of connection state
type ConnectionState int

// Constants to manage connection state
const (
	Disconnected ConnectionState = iota
	Connected
	ConnectionCheck   = 5 * time.Second
	InitialRetryDelay = 2 * time.Second
	MaxRetries        = 4
	MaxRetryDelay     = 10 * time.Second
)
