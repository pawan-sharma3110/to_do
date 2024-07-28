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
	app.Get("/todos", handler.GetToDO)
	app.Patch("/api/todos/:id", handler.UpdateToDo)
	// app.Delete("/api/todos/:id")
	app.Post("/todo/create", handler.CreateTodo)

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	port := os.Getenv("PORT")
	
	log.Fatal(app.Listen(":" + port))
}
