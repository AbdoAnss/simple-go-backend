package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(mysql-service:%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Retry connection up to 5 times
	for i := 0; i < 5; i++ {
		err = DB.Ping()
		if err == nil {
			log.Println("Successfully connected to database")
			return
		}
		log.Printf("Failed to ping database, attempt %d/5: %v", i+1, err)
		time.Sleep(5 * time.Second)
	}

	log.Fatalf("Could not connect to database after 5 attempts")
}
