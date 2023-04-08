package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/joskeiner/Go_Postgre-Gorm/models"
	"github.com/joskeiner/Go_Postgre-Gorm/storage"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title" `
	Publisher string `json:"publisher"`
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	// Ctx representa el contexto que mantiene la peticion http

	book := Book{}

	err := context.BodyParser(&book) // funcion para manejar la descodificacion de json si no puede descodificarlo retorna un error  ErrUnprocessableEntit

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request faild"})
		return err
	}

	err = r.DB.Create(&book).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "dould not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book has been added"})
	return nil

}

func (r *Repository) DelateBook(context *fiber.Ctx) error {
	bookModel := models.Books{}

	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id connot be empty"})
		return nil
	}

	err := r.DB.Delete(bookModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not delete book "})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book delete succesfully"})

	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModels := &models.Books{}

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	fmt.Println("the ID is ", id) //  erificar el id

	err := r.DB.Where("id = ?", id).First(bookModels).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could no get the book"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "book id fetched successfully",
			"data": bookModels})

	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {

	BookModels := &[]models.Books{}

	err := r.DB.Find(BookModels).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "books fetched successfully",
			"data":    BookModels,
		})
	return nil
}

// la variable DB es una estructra de gorm la cual tiene el statement , config, error , RowsAffected
// practicamente es la definicion de la base de datos
type Repository struct {
	DB *gorm.DB
}

// la variable app tiene el valor de App de fiber que es una struct inicial de fiber
func (r *Repository) SetupRoutes(app *fiber.App) {

	api := app.Group("/api") // Group lo que hace es agrupar estas funciones bajo el slash /api

	api.Post("/create_Books", r.CreateBook)
	api.Delete("/delate_Book/:id", r.DelateBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{

		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SslMode:  os.Getenv("DB_SSLMODE"),
		DbName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database ")
	}

	err = models.MigrateBook(db)

	if err != nil {
		log.Fatal("could not load the database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()

	r.SetupRoutes(app)

	app.Listen(":8080")

}
