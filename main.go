/*
 * Secret Server
 *
 * This is an API of a secret service. You can save your secret by using the API. You can restrict the access of a secret after the certen number of views or after a certen period of time.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package main

import (
	"log"
	"net/http"
	"os"

	sw "github.com/jakule/codersranktask/go"
	"github.com/joho/godotenv"
)

//go:generate mockgen -source=go/storage.go -destination go/mocks/storage_mock.go

func main() {
	log.Printf("Server started")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Heroku thing...
	addr := ":8080"
	if addrEnv := os.Getenv("PORT"); addrEnv != "" {
		addr = ":" + addrEnv
	}

	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	err = sw.MigrateDB(dbConnStr)
	if err != nil {
		log.Fatal(err)
	}

	router := sw.NewRouter(dbConnStr)
	log.Fatal(http.ListenAndServe(addr, router))
}
