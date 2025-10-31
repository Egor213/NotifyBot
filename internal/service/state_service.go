package service

import (
	"sync"

	"github.com/Egor213/notifyBot/internal/entity"
)

type UserState struct {
	State entity.StateType
	Data  map[any]any
}

type StateService struct {
	mu     sync.RWMutex
	states map[int64]*UserState
}

func NewStateService() *StateService {
	return &StateService{
		states: make(map[int64]*UserState),
	}
}

func (sm *StateService) SetState(chatID int64, state entity.StateType, data map[any]any) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.states[chatID] = &UserState{
		State: state,
		Data:  data,
	}
}

func (sm *StateService) GetState(chatID int64) *UserState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	if state, ok := sm.states[chatID]; ok {
		return state
	}
	return &UserState{State: entity.StateNone}
}

func (sm *StateService) ClearState(chatID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.states, chatID)
}
