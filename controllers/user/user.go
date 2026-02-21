package user

import (
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v3"

	"backend/config"
	"backend/models"
	"backend/repository"
	"backend/utils"
)

var userRepo = repository.NewUserRepository()

func GetAllUsers(c fiber.Ctx) error {
	users, err := userRepo.GetAll(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to fetch users")
	}
	return utils.SuccessResponse(c, 200, "Users fetched", users)
}

func GetUserByID(c fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return utils.ErrorResponse(c, 400, "User ID is required")
	}

	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil {
		return utils.ErrorResponse(c, 404, "User not found")
	}

	return utils.SuccessResponse(c, 200, "User fetched", user)
}

func CreateUser(c fiber.Ctx) error {
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

	record, err := config.FirebaseAuth.CreateUser(c.Context(), params)
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

	if err := userRepo.Create(c.Context(), user); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to save user")
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

	user, err := userRepo.FindByID(c.Context(), userID)
	if err != nil || user == nil {
		return utils.ErrorResponse(c, 404, "User not found")
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	user.UpdatedAt = time.Now()

	if err := userRepo.Update(c.Context(), user); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to update user")
	}

	return utils.SuccessResponse(c, 200, "User updated successfully", user)
}
