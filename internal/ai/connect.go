package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Noskash/dostonproject/internal/models"
	"github.com/joho/godotenv"
)

const groqURL = "https://api.groq.com/openai/v1/chat/completions"

func Send_api_request(html []string, title string) string {
	err := godotenv.Load("dostonproject/.env")
	if err != nil {
		log.Fatal("Ошибка при загрузке .env файла", err)
	}
	req_string := fmt.Sprintf("на основе данного тебе html файла верни новый html файл весь контент перефразируй так чтобы он отличался от оригинального на 70% но чтобы сам смысл статьи не менялся \n %s	", html)
	api := os.Getenv("API_KEY")
	req := models.ApiRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []models.ApiMessage{{
			Role:    "user",
			Content: req_string,
		},
		},
	}
	req_body, _ := json.Marshal(req)
	request, _ := http.NewRequest("POST", groqURL, bytes.NewBuffer(req_body))
	request.Header.Set("Authorization", "Bearer "+api)
	request.Header.Set("Content-type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Ошибка чтения json ответа", err)
	}

	var apiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &apiResponse); err != nil {
		log.Fatal("Ошибка при распаковке json файла", err)
	}
	return apiResponse.Choices[0].Message.Content
}
