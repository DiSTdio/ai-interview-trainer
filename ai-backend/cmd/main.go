package main

import (
	"log"
	"net/http"
	"time"

	"ai-backend/internal/api"
	"ai-backend/internal/llm"
	"ai-backend/internal/memory"
	"ai-backend/internal/rate"
)

func main() {

	h := memory.New()
	l := rate.New(5, 10*time.Second)

	llmClient := llm.New(
		"http://host.docker.internal:11434",
		"qwen2.5-coder:3b",
	)

	handler := api.NewHandler(llmClient, h, l)

	http.Handle("/ai", handler)

	log.Println("🚀 :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
