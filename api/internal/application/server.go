package application

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/lordvidex/gomoney/api/internal/token"
	"google.golang.org/grpc"
)

var (
	serverAddr = flag.String("server_addr", "localhost:50051", "The server address in the format of host:port")
)

type Server struct {
	grpc  *grpc.ClientConn
	app   *fiber.App
	maker token.Maker
}

func NewServer() (*Server, error) {
	app := fiber.New()
	grpcServer, err := connectGRPC()

	if err != nil {
		return &Server{}, err
	}

	// TODO: Load symmetricKey from env
	symmetricKey := []byte("01234567890123456789012345678912")
	maker := token.NewPasetoMaker(symmetricKey)

	server := &Server{
		grpc:  grpcServer,
		app:   app,
		maker: maker,
	}

	server.setupRoutes()
	return server, nil
}

func (server *Server) setupRoutes() {
	apiAuth := server.app.Group("/", server.authMiddleware)

	// Unauthenticated routes
	server.app.Post("/login", server.login)

	// Authenticated EndPoints
	// -	Users EndPoint
	apiAuth.Get("/api/users", server.getUser)

	// -	Accounts EndPoint
	apiAuth.Post("/api/accounts", server.createAccount)
	apiAuth.Get("/api/accounts", server.getAccounts)

	// -	Transactions EndPoint
	apiAuth.Post("/api/transactions", server.createTransaction)
	apiAuth.Get("/api/transactions/:id", server.getTransactions)
}

func connectGRPC() (*grpc.ClientConn, error) {
	grpc, err := grpc.Dial(*serverAddr)
	if err != nil {
		return nil, err
	}

	return grpc, nil
}

func (server *Server) Run() {
	server.app.Listen(":8000")
}
