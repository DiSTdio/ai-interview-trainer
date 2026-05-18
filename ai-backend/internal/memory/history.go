package memory

import "sync"

type History struct {
	mu   sync.RWMutex
	data map[string][]string
}

func New() *History {
	return &History{
		data: make(map[string][]string),
	}
}

func (h *History) Get(userID string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return append([]string{}, h.data[userID]...)
}

func (h *History) Append(userID, msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.data[userID] = append(h.data[userID], msg)

	if len(h.data[userID]) > 10 {
		h.data[userID] = h.data[userID][len(h.data[userID])-10:]
	}
}
