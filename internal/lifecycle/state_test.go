package lifecycle

import (
	"testing"
)

func TestStateStoreInitialState(t *testing.T) {
	storeState := newStateStore()

	if got := storeState.current(); got != serverStateStarting {
		t.Fatalf("got %v, want %v", got, serverStateStarting)
	}
}

func TestStateStoreTransition(t *testing.T) {
	storeState := newStateStore()

	if ok := storeState.transition(serverStateRunning); !ok {
		t.Fatal("expected transition to succedd")
	}

	if got := storeState.current(); got != serverStateRunning {
		t.Fatalf("got %v, want %v", got, serverStateRunning)
	}
}

func TestInvalidStateStoretransition(t *testing.T) {
	storeState := newStateStore()

	if ok := storeState.transition(serverStateDraining); ok {
		t.Fatal("expected transition to fail")
	}

	if got := storeState.current(); got != serverStateStarting {
		t.Fatalf("got %v, want %v", got, serverStateStarting)
	}
}

func TestStateTransition(t *testing.T) {
	testStruct := []struct {
		name      string
		fromState ServerState
		toState   ServerState
		allowed   bool
	}{
		{
			name:      "starting to running",
			fromState: serverStateStarting,
			toState:   serverStateRunning,
			allowed:   true,
		},
		{
			name:      "running to draining",
			fromState: serverStateRunning,
			toState:   serverStateDraining,
			allowed:   true,
		},
		{
			name:      "draining to stopped",
			fromState: serverStateDraining,
			toState:   serverStateStopped,
			allowed:   true,
		},
		{
			name:      "starting to draining",
			fromState: serverStateStarting,
			toState:   serverStateDraining,
			allowed:   false,
		},
		{
			name:      "running to stopped",
			fromState: serverStateRunning,
			toState:   serverStateStopped,
			allowed:   false,
		},
		{
			name:      "stopped to running",
			fromState: serverStateStopped,
			toState:   serverStateRunning,
			allowed:   false,
		},
	}

	for _, tt := range testStruct {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidTransition(tt.fromState, tt.toState); got != tt.allowed {
				t.Fatalf("got %v, want %v", got, tt.allowed)
			}
		})
	}
}
