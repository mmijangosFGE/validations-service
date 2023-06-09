package messages

const (
	ConnectToMongoDBFailed    = "Failed to connect to MongoDB, retrying..."
	ConnectToMongoDBSuccess   = "Connected MongoDB successfully"
	ConnectionLost            = "Connection lost, attempting to reconnect..."
	ExpectedConnectedState    = "Expected state to be Connected, got %v"
	ExpectedDisconnectedState = "Expected state to be Disconnected, got %v"
	FailedToCloseConnection   = "Failed to close MongoDB connection: %v"
	MaximumNumberRetries      = "Reached maximum number of retries, stopping attempts to reconnect"
	PingToMongoDBFailed       = "Failed to ping MongoDB, retrying connection..."
)
