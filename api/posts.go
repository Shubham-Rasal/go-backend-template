package api

import (
	"database/sql"
	"strconv"

	db "github.com/Shubham-Rasal/blog-backend/db/sqlc"
	"github.com/gofiber/fiber/v2"
)

type createPostRequest struct {
	Title  string `validate:"required"`
	Body   string `validate:"required"`
	UserID int32  `validate:"required"`
	Status string `validate:"required"`
}

func (server *Server) createPost(c *fiber.Ctx) error {
	var req createPostRequest
	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	if err := server.validator.Struct(req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	arg := db.CreatePostParams{
		Title:  req.Title,
		Body:   req.Body,
		UserID: req.UserID,
		Status: req.Status,
	}

	user, err := server.store.CreatePost(c.Context(), arg)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	c.JSON(user)
	return nil

}

func (server *Server) getPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	post, err := server.store.GetPost(c.Context(), int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			c.Status(fiber.StatusNotFound).JSON(errorResponse(err))
			return err
		}
		c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	c.JSON(post)
	return nil
}

func (server *Server) listPosts(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	arg := db.ListPostsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	posts, err := server.store.ListPosts(c.Context(), arg)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	c.JSON(posts)
	return nil
}

func (server *Server) likePost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(errorResponse(err))
		return err
	}

	err = server.store.LikePost(c.Context(), int64(id))
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(errorResponse(err))
		return err
	}

	c.JSON(true)
	return nil
}
