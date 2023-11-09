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
	Username string `validate:"required,min=3,max=32"`
	Password string `validate:"required,min=6,max=32"`
	Email    string `validate:"required,email"`
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
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	// Validate request
	if err := server.validator.Struct(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	args := db.ListUsersParams{
		Limit:  params.PageSize,
		Offset: (params.PageId - 1) * params.PageSize,
	}

	users, err := server.store.ListUsers(context.Background(), args)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
	}

	c.JSON(users)
	return nil

}

func (server *Server) createUser(c *fiber.Ctx) error {
	var req createUserRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))

	}

	if err := server.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))

	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	user, err := server.store.CreateUser(c.Context(), arg)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	c.Status(fiber.StatusCreated).JSON(user)
	return nil

}

type getUserRequest struct {
	UserID int `validate:"required,min=1"`
}

func (server *Server) getUserById(c *fiber.Ctx) error {
	// Parse URI params
	var req getUserRequest
	req.UserID, _ = strconv.Atoi(c.Params("id"))

	// Validate request
	if err := server.validator.Struct(req); err != nil {
		// fmt.Print("validation err : ", err)
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		// return err
	}

	// Get user from database
	user, err := server.store.GetUserById(context.Background(), int64(req.UserID))
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

type getUserByUsernameRequest struct {
	Username string `validate:"required,min=1"`
}

func (server *Server) getUserByUsername(c *fiber.Ctx) error {
	// Parse URI params
	var req getUserByUsernameRequest
	req.Username = c.Params("username")

	fmt.Println("username : ", req.Username)

	// Validate request
	if err := server.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		// return err
	}

	//check if empty or not a valid username

	if len(req.Username) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(fmt.Errorf("username cannot be empty")))
	}

	//check if its not a number
	if _, err := strconv.Atoi(req.Username); err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(fmt.Errorf("username cannot be a number")))
	}

	// Get user from database
	user, err := server.store.GetUserByUsername(context.Background(), req.Username)
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
	req.UserID, _ = strconv.Atoi(c.Params("id"))

	// Validate request
	if err := server.validator.Struct(req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	// Get user from database
	err := server.store.DeleteUser(context.Background(), int64(req.UserID))
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	// Return user as JSON response
	c.JSON(req.UserID)
	return nil
}
