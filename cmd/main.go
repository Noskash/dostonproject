package main

import (
	"log"
	"net/http"

	"github.com/Noskash/dostonproject/db"
	"github.com/Noskash/dostonproject/internal/src"
)

func main() {
	db, err := db.Connect_to_database()
	if err != nil {
		log.Fatal("Ошибка при подключении к базе данных", err)
	}
	http.HandleFunc("/cretehtml", src.Get_html(db))
	http.ListenAndServe(":8080", nil)
}
