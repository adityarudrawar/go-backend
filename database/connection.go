package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)


func getConnectionString() string {
	host, port, user, password, dbname := getdbCred()
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return psqlInfo
}

func CreateConnection() *sql.DB {
	
	connectionString := getConnectionString()
	
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	return db
}

func CloseDB(db *sql.DB) error {
	return db.Close()
}

func getdbCred() (string, string, string, string, string) {

	err := godotenv.Load()
	if err != nil {
		log.Panic("Failed to load ENV")
	}
	POSTGRES_HOST := os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT := os.Getenv("POSTGRES_PORT")

	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DBNAME := os.Getenv("POSTGRES_DBNAME")

	return POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DBNAME
}