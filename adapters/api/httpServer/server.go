package httpServer

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mmijangosFGE/validations-service/adapters/api/routes"
	"github.com/mmijangosFGE/validations-service/pkg/constants"
	"github.com/mmijangosFGE/validations-service/pkg/messages"
	"github.com/mmijangosFGE/validations-service/pkg/responses"
	"go.mongodb.org/mongo-driver/mongo"
)

// Config - struct of server configs
type Config struct {
	Env       string
	JWTSecret string
	MongoDb   string
	Password  string
	Port      string
	Username  string
}

// Server - interface of config struct
type Server interface {
	Config() *Config
}

// Broker - struct of broker
type Broker struct {
	config *Config
}

// Config - Receiver function to config broker
func (b *Broker) Config() *Config {
	return b.config
}

// NewServer - constructor of server with config
func NewServer(
	config *Config,
) (
	*Broker,
	error,
) {
	// Validate .env variables of config broker
	if config.Env == "" {
		return nil, errors.New(messages.EnvIsRequired)
	}
	if config.Port == "" {
		return nil, errors.New(messages.PortIsRequired)
	}
	if config.MongoDb == "" {
		return nil, errors.New(messages.MongoDbRequired)
	}
	if config.JWTSecret == "" {
		return nil, errors.New(messages.JWSecretRequired)
	}
	// Create pointer space to broker
	broker := &Broker{
		config: config,
	}
	// return broker
	return broker, nil
}

// Start - start server with config
func (b *Broker) Start(
	mc *mongo.Client,
	s Server,
) {
	// Create fiber app
	app := fiber.New()
	// activate middlewares
	app.Use(
		// logger middleware
		logger.New(),
		// cors middleware
		cors.New(),
	)
	// Configure basic authMock middleware
	bAuth := basicauth.New(basicauth.Config{
		Users: map[string]string{
			b.config.Username: b.config.Password,
		},
		// Custom unauthorized handler
		Unauthorized: func(c *fiber.Ctx) error {
			response, ctx := responses.GeneralResponse(
				c,
				fiber.StatusUnauthorized,
				false,
				messages.Unauthorized,
			)
			return ctx.JSON(response)
		},
	})
	// Init routes
	routes.BindRoutes(
		app,
		&bAuth,
	)
	// Validate env to serve https on local
	var errServe error
	if b.config.Env == constants.EnvLocal {
		// Initialize server on specific port
		errServe = app.ListenTLS(
			b.config.Port,
			"certificates/localhost.pem",
			"certificates/localhost-key.pem",
		)
	} else {
		errServe = app.Listen(b.config.Port)
	}

	// Verify if exists error when initialize server
	if errServe != nil {
		log.Fatal(errServe)
	}
}
