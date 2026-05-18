package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"ai-backend/internal/llm"
	"ai-backend/internal/memory"
	"ai-backend/internal/model"
	"ai-backend/internal/prompt"
	"ai-backend/internal/rate"
)

type Handler struct {
	llm     *llm.Client
	history *memory.History
	limiter *rate.Limiter
}

func NewHandler(llm *llm.Client, h *memory.History, l *rate.Limiter) *Handler {
	return &Handler{llm: llm, history: h, limiter: l}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req model.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", 400)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id required", 400)
		return
	}

	if !h.limiter.Allow(req.UserID) {
		http.Error(w, "rate limit", 429)
		return
	}

	histSlice := h.history.Get(req.UserID)
	histSlice = append(histSlice, "User: "+req.Input)

	if len(histSlice) > 10 {
		histSlice = histSlice[len(histSlice)-10:]
	}

	var p string

	switch req.Mode {

	case "chat":
		p = prompt.BuildChat(
			req.Input,
			req.Judge,
		)

	case "interview":
		p = prompt.BuildInterview(
			strings.Join(histSlice, "\n"),
			req.Judge,
			req.Topic,
		)

	default:
		p = prompt.BuildChat(
			req.Input,
			false,
		)
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "stream not supported", 500)
		return
	}

	// log
	fmt.Println("TOPIC:", req.Topic)
	fmt.Println("MODE:", req.Mode)
	fmt.Println("PROMPT:", p)

	answer, err := h.llm.Stream(p, w, flusher)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var score *int
	if req.Mode == "interview" && req.Judge {
		score = extractScore(answer)
	}

	if score != nil {
		fmt.Fprintf(w, "event: score\n")
		fmt.Fprintf(w, "data: %d\n\n", *score)
		flusher.Flush()
	}

	h.history.Append(req.UserID, "User: "+req.Input)
	h.history.Append(req.UserID, "AI: "+answer)

	flusher.Flush()
}

func extractScore(text string) *int {
	re := regexp.MustCompile(`(?i)(?:score[:\s]*)?(\d{1,2})\s*/?\s*10?`)
	match := re.FindStringSubmatch(text)

	if len(match) < 2 {
		return nil
	}

	var score int
	fmt.Sscanf(match[1], "%d", &score)
	return &score
}
