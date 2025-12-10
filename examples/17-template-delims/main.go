package main

import (
	"encoding/json"
	"net/http"
	"text/template"
)

// @title Template Delimiters API
// @version 1.0
// @host localhost:8080
// @BasePath /api
// @description API que usa template customizado com delimiters [[ ]]

type Message struct {
	Text string `json:"text"`
}

// GetMessage returns a templated message
// @Summary Get message
// @Tags messages
// @Produce json
// @Success 200 {object} Message
// @Router /message [get]
func GetMessage(w http.ResponseWriter, r *http.Request) {
	// Template com delimiters customizados
	tmpl := template.New("msg")
	tmpl.Delims("[[", "]]")

	tmpl, _ = tmpl.Parse("Hello [[.Name]]!")

	json.NewEncoder(w).Encode(Message{Text: "Templated response"})
}

func main() {
	http.HandleFunc("/api/message", GetMessage)
	http.ListenAndServe(":8080", nil)
}
