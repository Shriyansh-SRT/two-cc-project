package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Shriyansh-SRT/two-cc-project/models"
	"github.com/Shriyansh-SRT/two-cc-project/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// MARK: Repository
type Repository struct {
	DB *gorm.DB
}

// MARK: Routes
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_articles", r.CreateArticle)
	api.Delete("/delete_articles/:id", r.DeleteArticle)
	api.Get("/get_articles/:id", r.GetArticleById)
	api.Get("/articles", r.GetArticles)
}

// MARK: Article struct
type Article struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

// MARK: CreateArticle
func (r *Repository) CreateArticle(c *fiber.Ctx) error {
	article := Article{}

	err := c.BodyParser(&article)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&article).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "an error occurred while creating the article"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "article created successfully", "data": article})

	return nil

}

// MARK: DeleteArticle
func (r *Repository) DeleteArticle(c *fiber.Ctx) error {
	articleModel := models.Articles{}

	id := c.Params("id")

	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Delete(articleModel, id)

	if err.Error != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete article"})
		return err.Error
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "article deleted successfully"})

	return nil
}

// MARK: GetArticleById
func (r *Repository) GetArticleById(c *fiber.Ctx) error {
	articleModel := models.Articles{}
	id := c.Params("id")

	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.First(&articleModel, id)

	if err.Error != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not find article"})
		return err.Error
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "article found successfully", "data": articleModel})

	return nil
}

// MARK: Get All Articles
func (r *Repository) GetArticles(c *fiber.Ctx) error {
	articleModels := &[]models.Articles{}

	err := r.DB.Find(articleModels).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "an error occurred while fetching the articles"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "articles fetched successfully", "data": articleModels})
	return nil
}

// MARK: Main
func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//MARK: Database connection
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	err = models.MigrateArticles(db)

	if err != nil {
		log.Fatal("Error migrating database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
