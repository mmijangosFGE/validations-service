package mongo

import (
	"github.com/mmijangosFGE/validations-service/pkg/constants"
	"github.com/mmijangosFGE/validations-service/pkg/messages"
	"testing"
)

// TestSetState is a unit test function for the SetState method of the Connection struct.
// It verifies that the connection state is set correctly and returns the expected state.
func TestSetState(t *testing.T) {
	// Create a new Connection instance
	connection := &Connection{}

	// Set the state to "Connected"
	connection.setState(constants.Connected)

	// Check if the state is set correctly
	if connection.getState() != constants.Connected {
		t.Errorf(
			messages.ExpectedConnectedState,
			connection.getState(),
		)
	}

	// Set the state to "Disconnected"
	connection.setState(constants.Disconnected)

	// Check if the state is set correctly
	if connection.getState() != constants.Disconnected {
		t.Errorf(
			messages.ExpectedDisconnectedState,
			connection.getState(),
		)
	}
}
