package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect_to_database() (*sql.DB, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Ошибка при загрузке .env файла")
	}
	Db_Name := os.Getenv("DB_NAME")
	Port := os.Getenv("PORT")
	SSLMODE := os.Getenv("SSLMODE")
	Password := os.Getenv("PASSWORD")
	User := os.Getenv("User")
	connectStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=%s", User, Password, Port, Db_Name, SSLMODE)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных")
	}
	return db, nil
}
