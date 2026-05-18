package model

type Request struct {
	UserID string `json:"user_id"`
	Mode   string `json:"mode"`
	Input  string `json:"input"`
	Judge  bool   `json:"judge"`
	Topic  string `json:"topic"`
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

type Response struct {
	Answer string `json:"answer"`
	Score  *int   `json:"score,omitempty"`
}

type JudgeResult struct {
	Score          int      `json:"score"`
	Mistakes       []string `json:"mistakes"`
	ImprovedAnswer string   `json:"improved_answer"`
	NextQuestion   string   `json:"next_question"`
}
