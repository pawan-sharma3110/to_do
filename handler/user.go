package handler

import (
	"io"
	"time"
	"to_do/database"
	"to_do/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserHandler(c *fiber.Ctx) error {
	db, err := database.DbIn()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error opening database connection"})
	}
	defer db.Close()
	// Parse the multipart form
	if err := c.BodyParser(&struct{}{}); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Unable to parse form")
	}

	// Get the file from the form
	file, err := c.FormFile("profile_picture")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Unable to get file from form")
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to open file")
	}
	defer src.Close()

	// Read the file into a byte slice
	profilePicture, err := io.ReadAll(src)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to read file")
	}

	// Get other user details from the form
	user := models.User{
		ID:             uuid.New(),
		ProfilePicture: profilePicture,
		FullName:       c.FormValue("full_name"),
		Email:          c.FormValue("email"),
		PhoneNo:        c.FormValue("phone_no"),
		Password:       c.FormValue("password"),
		CreatedOn:      time.Now(),
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}
	user.Password = string(hashedPassword)

	// Insert the user into the database
	_, err = db.Exec(`
        INSERT INTO users (id, profile_picture, full_name, email, phone_no, password, created_on)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, user.ID, user.ProfilePicture, user.FullName, user.Email, user.PhoneNo, user.Password, user.CreatedOn)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}

	// Respond with the created user
	return c.Status(fiber.StatusCreated).JSON(user)
}
