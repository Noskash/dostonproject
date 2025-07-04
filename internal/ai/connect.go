package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

const apiURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent"

type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Part `json:"parts"`
	Role  string `json:"role"`
}

type Request struct {
	Contents []Content `json:"contents"`
}

var loadEnvOnce sync.Once

func Send_api_request(html string, title string) (string, error) {
	loadEnvOnce.Do(func() {
		_ = godotenv.Load("../.env")
	})

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API_KEY пуст")
	}

	prompt := fmt.Sprintf("Перепиши данный тебе файл на 70%%, но так чтобы смысл не менялся и текст был максимум на 30%% меньше оригинального. Верни HTML-страницу с заголовком \"%s\". Вот содержание: %s", title, html)

	requestBody := Request{
		Contents: []Content{
			{
				Role: "user",
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("ошибка сериализации запроса: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL+"?key="+apiKey, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	fmt.Println(html)
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка от API (%d): %s", resp.StatusCode, string(respBytes))
	}

	// Распарсим только нужную часть
	var parsed struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return "", fmt.Errorf("ошибка при декодировании JSON: %v", err)
	}

	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("пустой ответ от модели")
	}

	text := parsed.Candidates[0].Content.Parts[0].Text

	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```html")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	return text, nil
}
