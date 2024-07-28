package main

import (
	"fmt"
	"log"
	"os"
	"to_do/models"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello World")
	app := fiber.New()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	port := os.Getenv("PORT")
	todos := []models.ToDo{}
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &models.ToDo{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is require"})
		}
		todo.ID = len(todos) + 1
		todos = append(todos, *todo)
		return c.Status(201).JSON(todo)
	})

	// Update to_do
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Complete = true
				return c.Status(200).JSON(todos[i])
			}
		}
		return c.Status(400).JSON(fiber.Map{"error": "Todo not found"})
	})
	// Delete to_do
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"success": "true"})

			}
		}
		return c.Status(400).JSON(fiber.Map{"error": "Todo not found"})

	})
	log.Fatal(app.Listen(":" + port))
}
