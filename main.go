package main

import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	print("ping")
	var (
		CLOUD_NAME = os.Getenv("CLOUD_NAME")
		API_KEY    = os.Getenv("API_KEY")
		API_SECRET = os.Getenv("API_SECRET")
		PORT       = os.Getenv("PORT")
	)

	if CLOUD_NAME == "" {
		panic("Enviroment variables not found")
	}

	if PORT == "" {
		PORT = "3000"
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 300 * 1024 * 1024
	})

	app.Use(cors.New())

	app.Static("/", "./public")

	app.Post("/", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")

		maxSize := int64(300 * 1024 * 1024)

		if err != nil {
			return c.JSON(fiber.Map{
				"error": "multipart form error",
			})
		}

		if file.Size > maxSize {
			return c.JSON(fiber.Map{
				"error": "max size",
			})
		}

		cld, err := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
		if err != nil {
			return err
		}

		res, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{})
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"url": res.SecureURL,
		})
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"response": "ok",
		})
	})

	app.Listen(":" + PORT)
}
