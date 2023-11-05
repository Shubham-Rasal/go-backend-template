package api

import (
	"context"

	db "github.com/Shubham-Rasal/blog-backend/db/sqlc"
	"github.com/gofiber/fiber/v2"
)

type createUserRequest struct {
	Username string `validate:"required"`
	Role     string `validate:"required,oneof=admin user"`
}

func (server *Server) listUsers(c *fiber.Ctx) error {
	// TODO: Implement

	args := db.ListUsersParams{
		Limit:  10,
		Offset: 0,
	}

	users, err := server.store.ListUsers(context.Background(), args)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	c.JSON(users)
	return nil

}

func errorResponse(err error) fiber.Map {
	return fiber.Map{"error": err.Error()}
}

func (server *Server) createUser(c *fiber.Ctx) error {
	var req createUserRequest
	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	if err := server.validator.Struct(req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Role:     req.Role,
	}

	user, err := server.store.CreateUser(c.Context(), arg)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	c.JSON(user)
	return nil

}

func (server *Server) getUser(c *fiber.Ctx) error {
	return nil
}

func (server *Server) deleteUser(c *fiber.Ctx) error {
	// TODO: Implement
	return nil
}
