package api

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	auth := api.Group("/auth")
	auth.Post("/register", userHandler.HandleRegister)
	auth.Post("/login", userHandler.HandleLogin)

	users := v1.Group("/users")
	users.Get("/me", middleware.WithJWTAuth(userHandler.HandleMe, *userStore))

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Headers("Authorization"))
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}

			var parsedMsg string
			if err := json.Unmarshal(msg, &parsedMsg); err != nil {
				log.Println("json unmarshal:", err)
				continue
			}

			trimmedMsg := strings.TrimSpace(parsedMsg)
			if trimmedMsg == "chat" {
				log.Println("Chat")
			}

			log.Printf("Received message: '%s'", string(msg))

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	return app.Listen(s.addr)
}
