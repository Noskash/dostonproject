package src

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func get_html(url string) string {
	res, err := http.Get(url)
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
	title := ""
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		s.SetText(title)
	})
	html, err := doc.Html()
	if err != nil {
		log.Fatal(err)
	}
	return html
}
