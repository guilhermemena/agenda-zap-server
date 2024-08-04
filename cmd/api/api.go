package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guilhermemena/agenda-zap-server/handlers/user"
	"github.com/guilhermemena/agenda-zap-server/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIServer struct {
	addr string
	db   *pgxpool.Pool
}

func NewAPIServer(addr string, db *pgxpool.Pool) *APIServer {
	return &APIServer{addr: addr, db: db}
}

func (s *APIServer) Run() error {
	var userStore = storage.NewUserStorage(s.db)
	var userHandler = user.NewUserHandler(*userStore)

	app := fiber.New()
	api := app.Group("/api")
	auth := api.Group("/auth")

	auth.Post("/register", userHandler.CreateUser)

	return app.Listen(s.addr)
}
