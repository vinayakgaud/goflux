package lifecycle

import (
	"sync/atomic"
)

type ServerState uint8

const (
	serverStateStarting ServerState = iota
	serverStateRunning
	serverStateDraining
	serverStateStopped
)

type stateStore struct {
	value atomic.Uint32
}

func newStateStore() *stateStore {
	s := &stateStore{}
	s.value.Store(uint32(serverStateStarting))
	return s
}

func (s *stateStore) current() ServerState {
	return ServerState(s.value.Load())
}

func (s *stateStore) transition(toState ServerState) bool {
	for {
		currentState := ServerState(s.value.Load())

		if !isValidTransition(currentState, toState) {
			return false
		}

		if s.value.CompareAndSwap(
			uint32(currentState),
			uint32(toState),
		) {
			return true
		}
	}
}

func isValidTransition(fromState ServerState, toState ServerState) bool {
	switch fromState {
	case serverStateStarting:
		return toState == serverStateRunning || toState == serverStateStopped

	case serverStateRunning:
		return toState == serverStateDraining

	case serverStateDraining:
		return toState == serverStateStopped

	case serverStateStopped:
		return false

	default:
		return false
	}
}
