package services

import (
	"sync"
	"time"
)

type TimerService struct {
	mu sync.RWMutex

	active bool
	endsAt int64
}

func NewTimerService() *TimerService {
	return &TimerService{}
}

func (s *TimerService) Start30Seconds() {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.active = true

	s.endsAt = time.Now().
		Add(30 * time.Second).
		Unix()
}

func (s *TimerService) State() (
	bool,
	int64,
) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.active {
		return false, 0
	}

	return true, s.endsAt
}

func (s *TimerService) SecondsLeft() int64 {

	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.active {
		return 0
	}

	left := s.endsAt - time.Now().Unix()

	if left < 0 {
		return 0
	}

	return left
}
