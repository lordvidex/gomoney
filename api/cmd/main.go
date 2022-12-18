package main

import (
	"log"

	"github.com/lordvidex/gomoney/api/internal/adapters/encryption"
	mgrpc "github.com/lordvidex/gomoney/api/internal/adapters/grpc"
	"github.com/lordvidex/gomoney/api/internal/adapters/handlers"
	"github.com/lordvidex/gomoney/api/internal/adapters/paseto"
	"github.com/lordvidex/gomoney/api/internal/adapters/redis"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/pkg/config"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//	@title			GoMoni API
//	@version		1.0
//	@description	This is the API for GoMoni, a simple money management application.
//
//	@contact.email	evans.dev99@gmail.com
//	@contact.name	Evans Owamoyo

//	@host		localhost:8000
//	@BasePath	/api
func main() {
	// read configs
	c := config.New()

	// create grpc service
	grpconn, err := connectGRPC(c)
	if err != nil {
		log.Fatal(err)
	}
	service := mgrpc.New(grpconn)

	// create token helper
	symmetricKey := c.Get("SYMMETRIC_KEY")
	if symmetricKey == "" {
		log.Fatal("env key 'SYMMETRIC_KEY' not set")
	}
	th := paseto.New([]byte(symmetricKey))

	// create redis client & userRepo
	client := redis.NewConnection(c)
	userRepo := redis.NewUserRepo(client)

	// bind application
	ph := encryption.NewBcryptPasswordHasher()
	app := application.New(userRepo, th, service, ph)

	// drive application
	restHandler := handlers.New(app)
	if err = restHandler.Listen(); err != nil {
		log.Fatal(err)
	}
}

func connectGRPC(c *config.Config) (*grpc.ClientConn, error) {
	server := c.Get("GRPC_SERVER")
	if server == "" {
		return nil, errors.New("key 'GRPC_SERVER' not set")
	}
	return grpc.Dial(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
