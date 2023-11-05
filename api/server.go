package api

import (
	db "github.com/Shubham-Rasal/blog-backend/db/sqlc"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	store     db.Store
	router    *fiber.App
	validator *validator.Validate
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store:     store,
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}

	router := fiber.New()

	router.Get("/users", server.listUsers)
	router.Post("/users", server.createUser)
	router.Get("/users/:id", server.getUser)
	router.Delete("/users/:id", server.deleteUser)

	router.Post("/posts", server.createPost)
	router.Get("/posts/:id", server.getPost)
	router.Get("/posts", server.listPosts)
	router.Post("/posts/:id/like", server.likePost)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Listen(address)
}
