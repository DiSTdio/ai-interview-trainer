package rate

import (
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
	users  map[string][]time.Time
	Limit  int
	Window time.Duration
}

func New(limit int, window time.Duration) *Limiter {
	return &Limiter{
		users:  make(map[string][]time.Time),
		Limit:  limit,
		Window: window,
	}
}

func (rl *Limiter) Allow(userID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	timestamps := rl.users[userID]

	var valid []time.Time
	for _, t := range timestamps {
		if now.Sub(t) < rl.Window {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.Limit {
		rl.users[userID] = valid
		return false
	}

	valid = append(valid, now)
	rl.users[userID] = valid

	return true
}
