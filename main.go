package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/ynoacamino/ynoa-uploader/config"
	"github.com/ynoacamino/ynoa-uploader/db"
	"github.com/ynoacamino/ynoa-uploader/middlewares"
	"github.com/ynoacamino/ynoa-uploader/routes"
)

func main() {
	config.InitEnv()

	db.InitDBConnection()
	defer db.CloseDBConnection()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := mux.NewRouter().StrictSlash(true)

	app.Use(middlewares.ConncetionSecret)

	uploaderRouter := app.PathPrefix("/api").Subrouter()

	routes.SetUpUploaderRoutes(uploaderRouter)

	if err := http.ListenAndServe(":8000", app); err != nil {
		panic(err)
	}
}
