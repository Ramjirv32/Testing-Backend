package controllers

import (
	"context"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v3"

	"backend/config"
	"backend/models"
	"backend/utils"
)

func GetAllUsers(c fiber.Ctx) error {
	return utils.SuccessResponse(c, 200, "Users fetched", []models.User{})
}

func GetUserByID(c fiber.Ctx) error {
	userID := c.Params("id")

	user := &models.User{
		ID:        userID,
		Email:     "user@example.com",
		Name:      "John Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return utils.SuccessResponse(c, 200, "User fetched", user)
}

func CreateUser(c fiber.Ctx) error {
	ctx := context.Background()
	var req models.CreateUserRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if req.Email == "" || req.Name == "" {
		return utils.ErrorResponse(c, 400, "Email and name are required")
	}

	uuidv7 := utils.GenerateUUIDv7()

	params := (&auth.UserToCreate{}).
		Email(req.Email).
		DisplayName(req.Name).
		UID(uuidv7)

	record, err := config.FirebaseAuth.CreateUser(ctx, params)
	if err != nil {
		return utils.ErrorResponse(c, 400, "Failed to create user")
	}

	seqID := utils.GetNextSeqID()

	user := &models.User{
		ID:        record.UID,
		SeqID:     seqID,
		Email:     req.Email,
		Name:      req.Name,
		Phone:     req.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return utils.SuccessResponse(c, 201, "User created successfully", user)
}

func UpdateUser(c fiber.Ctx) error {
	userID := c.Params("id")
	var req models.UpdateUserRequest

	if err := c.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request body")
	}

	if userID == "" {
		return utils.ErrorResponse(c, 400, "User ID is required")
	}

	user := &models.User{
		ID:        userID,
		Name:      req.Name,
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		UpdatedAt: time.Now(),
	}

	return utils.SuccessResponse(c, 200, "User updated successfully", user)
}
