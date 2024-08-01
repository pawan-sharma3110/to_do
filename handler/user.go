package handler

import (
	"database/sql"
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

func LoginHandler(c *fiber.Ctx) error {
	db, err := database.DbIn()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error opening database connection"})
	}
	defer db.Close()

	// Parse the form data
	var form struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&form); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Unable to parse form")
	}

	// Fetch the user from the database
	var user models.User
	err = db.QueryRow(`
        SELECT id, profile_picture, full_name, email, password
        FROM users
        WHERE email = $1
    `, form.Email).Scan(&user.ID, &user.ProfilePicture, &user.FullName, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching user"})
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// Respond with user details
	response := struct {
		ID             uuid.UUID `json:"id"`
		FullName       string    `json:"full_name"`
		ProfilePicture []byte    `json:"profile_picture"`
	}{
		ID:             user.ID,
		FullName:       user.FullName,
		ProfilePicture: user.ProfilePicture,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
