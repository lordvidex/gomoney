package main

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/gomoney/pkg/config"
	mygrpc "github.com/lordvidex/gomoney/server/internal/adapters/grpc"
	"github.com/lordvidex/gomoney/server/internal/adapters/postgres"
	"github.com/lordvidex/gomoney/server/internal/application"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	// read config
	c := config.New()

	// initialise the database connection
	conn, err := initDB(c)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.TODO())

	// run migrations
	_ = runMigrations(c)

	// driven adapters
	uRepo := postgres.NewUser(conn)
	aRepo := postgres.NewAccount(conn)

	// application
	app := application.New(uRepo, aRepo)

	// grpc driver
	server := grpc.NewServer()
	reflection.Register(server)
	handler := mygrpc.NewHandler(app)
	mygrpc.RegisterAccountServiceServer(server, handler)
	mygrpc.RegisterUserServiceServer(server, handler)

	// listen for incoming requests
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("server listening on port :8080")
	err = server.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}

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
