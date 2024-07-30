package main

import (
	"fmt"
	"log"
	"os"
	"to_do/database"
	"to_do/database/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	db, err := database.DbIn()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
	fmt.Println("Hello World")
	app := fiber.New()
	app.Static("/", "./public")
	app.Post("/todo/create", handler.CreateTodo)
	app.Get("/api/todos", handler.GetToDO)
	app.Patch("/api/todos/:id", handler.UpdateToDo)
	app.Patch("/api/todos/:id/complete", handler.CompleteToDo)
	app.Delete("/api/delete/todos/:id", handler.DeleteToDo)

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}
	fmt.Printf("Server start on port:%v", port)
	// Start the server
	log.Fatal(app.Listen(":" + port))

}
