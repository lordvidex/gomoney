package main

import (
	"context"
	"github.com/lordvidex/gomoney/server/internal/adapters"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/lordvidex/gomoney/pkg/config"
	mygrpc "github.com/lordvidex/gomoney/pkg/grpc"
	myhandler "github.com/lordvidex/gomoney/server/internal/adapters/grpc"
	"github.com/lordvidex/gomoney/server/internal/adapters/postgres"
	"github.com/lordvidex/gomoney/server/internal/application"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	log.Printf("pprof listening on port 8081")
	go func() {
		_ = http.ListenAndServe(":8081", nil)
	}()
}
func main() {
	appCtx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	// read config
	c := config.New()

	// initialise the database connection
	conn, err := postgres.NewConn(appCtx, c)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(appCtx)

	// driven adapters
	uRepo := postgres.NewUser(conn)
	aRepo := postgres.NewAccount(conn)
	locker := adapters.NewLocker(appCtx, time.Minute)

	// application
	app := application.New(uRepo, aRepo, locker)

	// grpc driver
	server := grpc.NewServer()
	reflection.Register(server)
	handler := myhandler.NewHandler(app)
	mygrpc.RegisterAccountServiceServer(server, handler)
	mygrpc.RegisterUserServiceServer(server, handler)
	mygrpc.RegisterTransactionServiceServer(server, handler)

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
