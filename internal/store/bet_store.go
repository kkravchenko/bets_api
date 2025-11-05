// store/bet_store.go
package store

import (
	"sync"
	"tz.api/internal/entity"
)

type BetStore struct {
	mu   sync.RWMutex
	bets map[string]entity.Bet
	list []entity.Bet
}

func NewBetStore() *BetStore {
	return &BetStore{
		bets: make(map[string]entity.Bet),
		list: make([]entity.Bet, 0),
	}
}

func (s *BetStore) Create(bet entity.Bet) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bets[bet.ID] = bet
	s.list = append(s.list, bet)
}

func (s *BetStore) GetAll() []entity.Bet {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copied := make([]entity.Bet, len(s.list))
	copy(copied, s.list)
	return copied
}

func (s *BetStore) GetByID(id string) (entity.Bet, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	bet, ok := s.bets[id]
	return bet, ok
}
