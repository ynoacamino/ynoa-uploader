package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var (
		R2_ENDPOINT           = "https://cde705ccd8ba0aec2e00e415023fef17.r2.cloudflarestorage.com"
		AWS_BUCKET            = os.Getenv("AWS_BUCKET")
		AWS_ACCESS_KEY_ID     = os.Getenv("AWS_ACCESS_KEY_ID")
		AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
		PORT                  = os.Getenv("PORT")
	)

	if AWS_BUCKET == "" || AWS_ACCESS_KEY_ID == "" || AWS_SECRET_ACCESS_KEY == "" {
		panic("Enviroment variables not found")
	}

	if PORT == "" {
		PORT = "3000"
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 300 * 1024 * 1024,
	})

	app.Use(cors.New())

	app.Static("/", "./public")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithCredentialsProvider(aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			AWS_ACCESS_KEY_ID,
			AWS_SECRET_ACCESS_KEY,
			"",
		))),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: R2_ENDPOINT,
			}, nil
		})),
	)
	if err != nil {
		panic("Error cargando la configuración de AWS: " + err.Error())
	}
	s3Client := s3.NewFromConfig(cfg)

	app.Post("/", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")

		maxSize := int64(3 * 1024 * 1024 * 1024)

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

		src, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "No se puede abrir el archivo",
			})
		}
		defer src.Close()

		key := fmt.Sprintf("uploads/%d_", time.Now().Unix())

		keyName := key + file.Filename

		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:             aws.String(AWS_BUCKET),
			Key:                aws.String(keyName),
			Body:               src,
			ACL:                types.ObjectCannedACLPublicRead,
			ContentDisposition: aws.String("inline"),
			ContentType:        aws.String(file.Header.Get("Content-Type")),
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Error con R2",
				"message": err.Error(),
			})
		}

		url := fmt.Sprintf("https://ynoa-uploader.ynoacamino.site/%s%s", key, url.PathEscape(file.Filename))

		return c.JSON(fiber.Map{
			"url": url,
		})
	})

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"response": "ok",
		})
	})

	app.Listen(":" + PORT)
}
