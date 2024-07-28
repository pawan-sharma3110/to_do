package handler

import (
	"database/sql"
	"strconv"
	"to_do/database"
	"to_do/models"

	"github.com/gofiber/fiber/v2"
)

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
		return c.Status(400).JSON(fiber.Map{"error": "Todo body is require"})
	}
	// Insert the new to-do item into the database
	query := `INSERT INTO todos (complete, body) VALUES ($1, $2) RETURNING id`
	err = db.QueryRow(query, todo.Complete, todo.Body).Scan(&todo.ID)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}
	// fiber.Map{"error": "Error inserting to-do item")
	// Return the created to-do item as JSON
	return c.Status(fiber.StatusCreated).JSON(*todo)
}

// func DeleteToDo(c *fiber.Ctx) error {

// }
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
	query := `UPDATE todos SET complete = $1, body = $2 WHERE id = $3 RETURNING id, complete, body`
	row := db.QueryRow(query, todo.Complete, todo.Body, todoID)
	// Scan the updated item
	updatedToDo := models.ToDo{}
	err = row.Scan(&updatedToDo.ID, &updatedToDo.Complete, &updatedToDo.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "To-do item not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating to-do item"})
	}

	// Return the updated to-do item as JSON
	return c.Status(fiber.StatusOK).JSON(updatedToDo)

}
func GetToDO(c *fiber.Ctx) error {
	db, err := database.DbIn()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error opening database connection"})
	}
	defer db.Close()
	// Fetch all to-do items
	rows, err := db.Query("SELECT id, complete, body FROM todos")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching to-do items",
		})
	}
	defer rows.Close()

	// Iterate through the rows and create a slice of ToDo
	var todos []models.ToDo
	for rows.Next() {
		var todo models.ToDo
		err := rows.Scan(&todo.ID, &todo.Complete, &todo.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error scanning to-do item",
			})
		}
		todos = append(todos, todo)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error iterating over to-do items",
		})
	}

	// Return the list of to-do items as JSON
	return c.JSON(todos)
}
