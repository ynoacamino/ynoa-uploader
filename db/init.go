package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/ynoacamino/ynoa-uploader/config"
)

var Query *Queries
var conn *pgx.Conn

func InitDBConnection() {
	ctx := context.Background()

	var err error
	conn, err = pgx.Connect(ctx, config.DATABASE_URL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		log.Println("Connected to database")
	}

	Query = New(conn)
}

func CloseDBConnection() {
	if conn != nil {
		log.Println("Closing database connection")
		conn.Close(context.Background())
	}
}
