package main

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/gomoney/api/internal/adapters/api"
	mgrpc "github.com/lordvidex/gomoney/api/internal/adapters/grpc"
	"github.com/lordvidex/gomoney/api/internal/adapters/memory"
	"github.com/lordvidex/gomoney/api/internal/adapters/paseto"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/pkg/config"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

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
	symmetricKey := []byte("01234567890123456789012345678912")
	th := paseto.New(symmetricKey)

	// create repository
	//dbconn, err := initDB(c)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer dbconn.Close(context.TODO())
	//err = runMigrations(c)

	//TODO: add persistent repo
	repo := memory.New()

	// bind application
	app := application.New(repo, th, service)

	// drive application
	restHandler := api.New(app)
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
func initDB(c *config.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.TODO(), c.Get("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}
	return conn, nil
}

func runMigrations(c *config.Config) error {
	m, err := migrate.New("file:///migrations", c.Get("DATABASE_URL"))
	if err != nil {
		return errors.Wrap(err, "failed to run migrations")
	}
	return m.Up()
}
