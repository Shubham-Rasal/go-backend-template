package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	db "github.com/Shubham-Rasal/blog-backend/db/sqlc"
	"github.com/gofiber/fiber/v2"
)

type createAccountRequest struct {
	Username string `validate:"required"`
	Role     string `validate:"required,oneof=admin user"`
	UserId   int32  `validate:"required"`
}

type ListAccountsQueryParams struct {
	PageId   int32 `validate:"required,min=1"`
	PageSize int32 `validate:"required,min=5,max=10"`
}

func (server *Server) listAccounts(c *fiber.Ctx) error {
	// TODO: Implement
	var params ListAccountsQueryParams
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

	args := db.ListAccountsParams{
		Limit:  params.PageSize,
		Offset: (params.PageId - 1) * params.PageSize,
	}

	users, err := server.store.ListAccounts(context.Background(), args)
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

func (server *Server) createAccount(c *fiber.Ctx) error {
	var req createAccountRequest
	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	if err := server.validator.Struct(req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	arg := db.CreateAccountParams{
		Username: req.Username,
		Role:     req.Role,
		UserID:   req.UserId,
	}

	user, err := server.store.CreateAccount(c.Context(), arg)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	c.JSON(user)
	return nil

}

type getAccountRequest struct {
	UserID int `validate:"required,min=1"`
}

func (server *Server) getAccount(c *fiber.Ctx) error {
	// Parse URI params
	var req getAccountRequest
	req.UserID, _ = strconv.Atoi(c.Params("id"))

	// Validate request
	if err := server.validator.Struct(req); err != nil {
		log.Print("validation error : ", err)
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		// return err
	}

	//get the username from header
	username := c.Locals("username").(string)

	caller, err := server.store.GetUserByUsername(context.Background(), username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
	}

	if caller.ID != int64(req.UserID) {
		return c.Status(fiber.StatusUnauthorized).JSON(errorResponse(fmt.Errorf("unauthorized: you are not allowed to access this account")))
	}

	// Get user from database
	account, err := server.store.GetAccount(context.Background(), int32(req.UserID))
	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(errorResponse(err))
		}

		c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	// Return user as JSON response
	c.JSON(account)
	return nil
}

func (server *Server) deleteAccount(c *fiber.Ctx) error {
	// Parse URI params
	var req getAccountRequest
	req.UserID, _ = strconv.Atoi(c.Params("id"))

	// Validate request
	if err := server.validator.Struct(req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	// Get user from database
	err := server.store.DeleteAccount(context.Background(), int32(req.UserID))
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	// Return user as JSON response
	c.JSON(req.UserID)
	return nil
}
