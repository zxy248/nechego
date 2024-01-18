package middleware

import (
	"sync"
	"time"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Throttle struct {
	Duration time.Duration

	lastMessage map[int64]time.Time
	mu          sync.Mutex
	once        sync.Once
}

func (m *Throttle) Wrap(next tele.HandlerFunc) tele.HandlerFunc {
	m.once.Do(func() {
		m.lastMessage = map[int64]time.Time{}
	})
	return func(c tele.Context) error {
		id := c.Sender().ID

		m.mu.Lock()
		d := time.Since(m.lastMessage[id])
		m.lastMessage[id] = time.Now()
		m.mu.Unlock()

		if d < m.Duration {
			return nil
		}
		return next(c)
	}
}
