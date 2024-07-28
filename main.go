package main

import (
	"fmt"
	"log"
	"to_do/models"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Hello World")
	app := fiber.New()

	todos := []models.ToDo{}
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "Hello World"})
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
	log.Fatal(app.Listen(":4000"))
}
