package events

import (
	"sync"
)

type Store struct {
	mu sync.Mutex
	events []Event
}

func NewStore() *Store {
	return &Store{
		events: []Event{},
	}
}

func (s *Store) Add(event Event) Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = append(s.events, event)

	return event
}

func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = []Event{}
}

func (s *Store) GetAll() []Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	cEvents := append([]Event{}, s.events...)

	return cEvents
}