package src

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/Noskash/dostonproject/internal/ai"
	"github.com/Noskash/dostonproject/internal/models"
	"github.com/PuerkitoBio/goquery"
)

func CleanTitle(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_\-]+`)
	clean := re.ReplaceAllString(name, "_")
	if len(clean) == 0 {
		return "default_output.txt"
	}
	return clean + ".html"
}

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
		var text string
		doc.Find("div , p , h1 , h2 , h3 , arcticle , img , strong , section, code , br").Each(func(i int, s *goquery.Selection) {
			content := strings.TrimSpace(s.Text())
			text += content
		})
		title := ""
		doc.Find("title").Each(func(i int, s *goquery.Selection) {
			title = strings.TrimSpace(s.Text())
		})
		result, err := ai.Send_api_request(text, title)
		if err != nil {
			log.Fatal("Ошибка отправки запроса ии", err)
		}
		ogTitle := ""
		doc.Find("meta[property='og:title']").Each(func(i int, s *goquery.Selection) {
			ogTitle, _ = s.Attr("content")
		})
		if ogTitle != "" {
			title = ogTitle
		}
		cleanTitle := CleanTitle(title)
		w.Header().Set("Content-type", "application/json")
		err = os.WriteFile("../internal/outputs/"+cleanTitle, []byte(result), 0644)
		if err != nil {
			log.Fatal("Не удалось сохранить файл", err)
		}
		if len(result) == 0 {
			log.Fatal("Пустой ответ", err)
		}

		path := fmt.Sprintf("outputs/%s", cleanTitle)
		_, err = db.Exec("INSERT INTO files(title , path) values($1 , $2)", cleanTitle, path)
		fmt.Printf(path)
		if err != nil {
			log.Fatal("Ошибка при вставке в бд", err)
		}

	}
}
