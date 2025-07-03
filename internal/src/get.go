package src

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Noskash/dostonproject/internal/ai"
	"github.com/Noskash/dostonproject/internal/models"
	"github.com/PuerkitoBio/goquery"
)

func Get_html(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.Request
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Fatal("Ошибка в json запросе ", err)
		}
		res, err := http.Get(request.Url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			log.Fatal("Неправильный статус", res.StatusCode)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal("Ошибка во время открытия", err)
		}
		doc.Find("script , style , nav , footer , header").Remove()
		var text []string
		doc.Find("p , h1 , h2 , h3 , arcticle , img , strong , section, code").Each(func(i int, s *goquery.Selection) {
			content := strings.TrimSpace(s.Text())
			text = append(text, content)
		})
		title := ""
		doc.Find("title").Each(func(i int, s *goquery.Selection) {
			s.SetText(title)
		})
		result := ai.Send_api_request(text, title)
		err = os.WriteFile("outputs/"+title, []byte(result), 0644)
		if err != nil {
			log.Fatal("Не удалось сохранить файл", err)
		}
		path := fmt.Sprintf("outputs/%s", title)
		db.Exec("INSERT INTO files(title , path) values($1 , $2)", title, path)
	}
}
