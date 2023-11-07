package api

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	db "github.com/Shubham-Rasal/blog-backend/db/sqlc"
	"github.com/gofiber/fiber/v2"
)

type createUserRequest struct {
	Username string `validate:"required"`
	Role     string `validate:"required,oneof=admin user"`
}

type ListUsersQueryParams struct {
	PageId   int32 `validate:"required,min=1"`
	PageSize int32 `validate:"required,min=5,max=10"`
}

func (server *Server) listUsers(c *fiber.Ctx) error {
	// TODO: Implement
	var params ListUsersQueryParams
	err := c.QueryParser(&params)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	// Validate request
	if err := server.validator.Struct(params); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	args := db.ListUsersParams{
		Limit:  params.PageSize,
		Offset: (params.PageId - 1) * params.PageSize,
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

type getUserRequest struct {
	ID int `validate:"required,min=1"`
}

func (server *Server) getUser(c *fiber.Ctx) error {
	// Parse URI params
	var req getUserRequest
	req.ID, _ = strconv.Atoi(c.Params("id"))

	// Validate request
	if err := server.validator.Struct(req); err != nil {
		fmt.Print("validation err : ", err)
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		// return err
	}

	// Get user from database
	user, err := server.store.GetUser(context.Background(), int64(req.ID))
	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(errorResponse(err))
		}

		c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	// Return user as JSON response
	c.JSON(user)
	return nil
}

func (server *Server) deleteUser(c *fiber.Ctx) error {
	// Parse URI params
	var req getUserRequest
	req.ID, _ = strconv.Atoi(c.Params("id"))

	// Validate request
	if err := server.validator.Struct(req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	// Get user from database
	err := server.store.DeleteUser(context.Background(), int64(req.ID))
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	// Return user as JSON response
	c.JSON(req.ID)
	return nil
}
