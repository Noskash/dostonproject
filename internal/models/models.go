package models

type ApiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ApiRequest struct {
	Model    string       `json:"model"`
	Messages []ApiMessage `json:"Messages"`
}

type Request struct {
	Url string `json:"url"`
}
