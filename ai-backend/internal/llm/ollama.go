package llm

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"ai-backend/internal/model"
)

type Client struct {
	URL   string
	Model string
}

func New(url string, model string) *Client {
	return &Client{
		URL:   url,
		Model: model,
	}
}

func (c *Client) Stream(prompt string, w http.ResponseWriter, flusher http.Flusher) (string, error) {

	reqBody := model.OllamaRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: true,
	}

	data, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		c.URL+"/api/generate",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	var full strings.Builder

	for {
		var chunk model.OllamaResponse
		if err := dec.Decode(&chunk); err != nil {
			break
		}

		full.WriteString(chunk.Response)

		_, _ = w.Write([]byte("data: " + chunk.Response + "\n\n"))
		flusher.Flush()
	}

	final := full.String()

	_, _ = w.Write([]byte("event: done\n"))
	_, _ = w.Write([]byte("data: " + final + "\n\n"))
	flusher.Flush()

	return final, nil

}
