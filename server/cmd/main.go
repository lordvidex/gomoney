package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v4"
	mygrpc "github.com/lordvidex/gomoney/server/internal/adapters/grpc"
	"github.com/lordvidex/gomoney/server/internal/adapters/postgres"
	"github.com/lordvidex/gomoney/server/internal/application"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// initialise the database connection
	conn, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.TODO())

	// driven adapters
	repo := postgres.NewRepository(conn)

	// application
	app := application.New(repo)

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

func initDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.TODO(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.TODO())
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}
	return conn, nil

}
