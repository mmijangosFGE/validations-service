package mongo

import (
	"context"
	"github.com/mmijangosFGE/validations-service/pkg/constants"
	"github.com/mmijangosFGE/validations-service/pkg/messages"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"sync"
	"time"
)

// DBConnector is a struct that implements the Connector interface.
// It uses the official MongoDB driver to establish a connection to the MongoDB database.
type DBConnector struct{}

// Connect is a method that establishes a connection to a MongoDB database using the provided
// context and client options. It returns a MongoDB client, which can be used
// to interact with the database, and any error encountered during the connection process.
func (c *DBConnector) Connect(
	ctx context.Context,
	opts *options.ClientOptions,
) (*mongo.Client, error) {
	return mongo.Connect(ctx, opts)
}

// Connector is an interface that represents an object that can establish a
// connection to a MongoDB database. It abstracts the actual mechanism of
// connecting to MongoDB, making it easier to mock in unit tests.
type Connector interface {
	Connect(
		context.Context,
		*options.ClientOptions,
	) (*mongo.Client, error)
}

// Connection struct holds the information related to MongoDB connection
type Connection struct {
	ClientChan chan *mongo.Client        // Channel to send connected MongoDB client to other parts of the application
	Connector  Connector                 // The connector interface
	State      constants.ConnectionState // Connection state (Connected or Disconnected)
	cancel     context.CancelFunc        // Context cancel function
	client     *mongo.Client             // MongoDB client instance
	stateMtx   sync.RWMutex              // Mutex for safely accessing the connection state
}

// EnsureConnection is a method that attempts to establish a connection to the MongoDB server,
// and if successful, pings the server to ensure that the connection is alive.
// If the maximum number of retries is reached, the application will log an error and exit.
func (c *Connection) EnsureConnection(url string) {
	retryDelay := constants.InitialRetryDelay
	retries := 0 // retries counter

	for {
		client, errConnect := c.Connector.Connect(
			context.Background(),
			options.Client().ApplyURI(url),
		)
		errPing := pingMongoDB(client)

		if errConnect != nil || errPing != nil {
			c.setState(constants.Disconnected)

			if retries++; retries >= constants.MaxRetries {
				// If the maximum number of retries is reached, log a message and exit the application
				log.Fatal(messages.MaximumNumberRetries)
			}

			// If there was a connection error, log the error, wait a while, and try again
			if errConnect != nil {
				log.Println(messages.ConnectToMongoDBFailed)
				if retryDelay *= 2; retryDelay > constants.MaxRetryDelay {
					retryDelay = constants.MaxRetryDelay
				}
				time.Sleep(retryDelay)
				continue
			}

			// If the ping failed, log the error, close the connection, wait a while, and try again
			if errPing != nil {
				log.Println(messages.PingToMongoDBFailed)
				c.closeMongoDB()
				c.setState(constants.Disconnected)
				time.Sleep(constants.MaxRetryDelay)
				continue
			}
		} else {
			// If the connection and ping were successful, set the client and state and break the loop
			c.setClient(client, nil)
			c.setState(constants.Connected)

			// Clear the ClientChan channel
			for len(c.ClientChan) > 0 {
				<-c.ClientChan // Delete the current value of the channel
			}
			c.ClientChan <- client
			log.Println(messages.ConnectToMongoDBSuccess)
			break
		}
	}
}

// MonitorConnection is a method that continuously checks the connection state and, if necessary,
// tries to reestablish the connection to the MongoDB server.
func (c *Connection) MonitorConnection(url string) {
	ticker := time.NewTicker(constants.ConnectionCheck)
	for {
		select {
		case <-ticker.C:
			state := c.getState()
			if state == constants.Disconnected {
				c.EnsureConnection(url)
			}
			if state == constants.Connected {
				// ping the server
				errPing := pingMongoDB(c.client)
				if errPing != nil {
					log.Println(messages.ConnectionLost)
					c.setState(constants.Disconnected)
					c.closeMongoDB() // Close MongoDB connection
					c.EnsureConnection(url)
				}
			}
		}
	}
}

// closeMongoDB is a method that closes the MongoDB connection and cancels the associated context
func (c *Connection) closeMongoDB() {
	// Lock the mutex to protect the state
	c.stateMtx.Lock()
	// Unlock the mutex after the function returns
	defer c.stateMtx.Unlock()

	if c.client == nil || c.cancel == nil {
		return
	}
	// Create a context with a timeout for disconnecting
	ctx, cancelClose := context.WithTimeout(context.Background(), 10*time.Second)
	// Cancel the context after the function returns
	defer cancelClose()
	// Disconnect the client
	if err := c.client.Disconnect(ctx); err != nil {
		log.Printf(messages.FailedToCloseConnection, err)
	}
	// Call the cancel function
	c.cancel()
}

// getState is a method that returns the current connection state (Connected or Disconnected)
func (c *Connection) getState() constants.ConnectionState {
	// Lock the mutex for read access to the state
	c.stateMtx.RLock()
	// Unlock the mutex after the function returns
	defer c.stateMtx.RUnlock()
	// Return the current connection state
	return c.State
}

// setClient is a method that sets the MongoDB client instance and the cancel function for the connection
func (c *Connection) setClient(
	client *mongo.Client,
	cancel context.CancelFunc,
) {
	// Lock the mutex to protect the state
	c.stateMtx.Lock()
	// Unlock the mutex after the function returns
	defer c.stateMtx.Unlock()
	// Assign the client instance
	c.client = client
	// Assign the cancel function
	c.cancel = cancel
}

// setState is a method that sets the connection state (Connected or Disconnected)
func (c *Connection) setState(state constants.ConnectionState) {
	// Lock the mutex to protect the state
	c.stateMtx.Lock()
	// Unlock the mutex after the function returns
	defer c.stateMtx.Unlock()
	// Set the connection state
	c.State = state
}

// pingMongoDB is a function that sends a ping to the MongoDB server to check if the connection is alive
var pingMongoDB = func(client *mongo.Client) error {
	// Create a context with a timeout for the ping operation
	ctx, cancel := context.WithTimeout(
		context.Background(),
		constants.ConnectionCheck,
	)
	// Cancel the context after the function returns
	defer cancel()
	// Send a ping to the MongoDB server and return any error
	return client.Ping(ctx, readpref.Primary())
}
