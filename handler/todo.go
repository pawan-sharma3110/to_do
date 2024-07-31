package handler

import (
	"database/sql"
	"strconv"
	"to_do/database"
	"to_do/models"

	"github.com/gofiber/fiber/v2"
)

// CreateTodo creates a new to-do item
func CreateTodo(c *fiber.Ctx) error {
	db, err := database.DbIn()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error opening database connection"})
	}
	defer db.Close()

	todo := &models.ToDo{}
	if err := c.BodyParser(todo); err != nil {
		return err
	}
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
	}
	// Insert the new to-do item into the database
	query := `INSERT INTO todos (complete, body) VALUES ($1, $2) RETURNING id`
	err = db.QueryRow(query, todo.Complete, todo.Body).Scan(&todo.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error inserting to-do item"})
	}

	// Return the created to-do item as JSON
	return c.Status(fiber.StatusCreated).JSON(*todo)
}

// DeleteToDo deletes a to-do item
func DeleteToDo(c *fiber.Ctx) error {
	db, err := database.DbIn()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error opening database connection"})
	}
	defer db.Close()
	id := c.Params("id")
	todoID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}
	// Delete the to-do item from the database
	query := `DELETE FROM todos WHERE id = $1`
	result, err := db.Exec(query, todoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting to-do item"})
	}

	// Check if any row was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error getting affected rows"})
	}
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "To-do item not found"})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "To-do item deleted successfully"})
}

// UpdateToDo updates an existing to-do item
func UpdateToDo(c *fiber.Ctx) error {
	db, err := database.DbIn()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error opening database connection"})
	}
	defer db.Close()
	id := c.Params("id")
	todoID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Parse the request body into the ToDo struct
	todo := new(models.ToDo)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request body"})
	}

	// Validate the input
	if todo.Body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Body is required"})
	}
	// Update the existing to-do item in the database
	query := `UPDATE todos SET complete = $1, body = $2 WHERE id = $3`
	result, err := db.Exec(query, todo.Complete, todo.Body, todoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating to-do item"})
	}

	// Check if any row was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error getting affected rows"})
	}
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "To-do item not found"})
	}

	// Return the updated to-do item as JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "To-do item updated successfully"})
}

// GetToDO gets all to-do items
func GetToDO(c *fiber.Ctx) error {
	db, err := database.DbIn()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error opening database connection"})
	}
	defer db.Close()

	// Query the database for all to-do items
	query := `SELECT id, complete, body FROM todos`
	rows, err := db.Query(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching to-do items"})
	}
	defer rows.Close()

	// Iterate over the rows and add them to the to-do list
	var todos []models.ToDo
	for rows.Next() {
		var todo models.ToDo
		err := rows.Scan(&todo.ID, &todo.Complete, &todo.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning to-do item"})
		}
		todos = append(todos, todo)
	}

	// Check for any error that occurred during iteration
	if err := rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error iterating over to-do items"})
	}

	// Return the list of to-do items as JSON
	return c.Status(fiber.StatusOK).JSON(todos)
}

// CompleteToDo toggles the "complete" status of a to-do item
func CompleteToDo(c *fiber.Ctx) error {
	db, err := database.DbIn()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error opening database connection"})
	}
	defer db.Close()
	id := c.Params("id")
	todoID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Fetch the current status of the to-do item
	var currentComplete bool
	err = db.QueryRow(`SELECT complete FROM todos WHERE id = $1`, todoID).Scan(&currentComplete)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "To-do item not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching to-do item status"})
	}

	// Toggle the "complete" status
	newComplete := !currentComplete

	// Update the "complete" status in the database
	_, err = db.Exec(`UPDATE todos SET complete = $1 WHERE id = $2`, newComplete, todoID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating to-do item status"})
	}

	// Return the updated status as JSON
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"complete": newComplete})
}
