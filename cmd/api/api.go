package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guilhermemena/agenda-zap-server/cmd/api/middleware"
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
	v1 := api.Group("/v1")

	auth := api.Group("/auth")
	auth.Post("/register", userHandler.HandleRegister)
	auth.Post("/login", userHandler.HandleLogin)

	users := v1.Group("/users")
	users.Get("/me", middleware.WithJWTAuth(userHandler.HandleMe, *userStore))

	return app.Listen(s.addr)
}
