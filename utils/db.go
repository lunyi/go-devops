package utils

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {

	var err error
	pwd := os.Getenv("DB_PWD")
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	port := os.Getenv("DB_PORT")
	db := os.Getenv("DB")
	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, pwd, db, host, port)
	DB, err = sql.Open("postgres", connString)

	CheckErr(err)
}