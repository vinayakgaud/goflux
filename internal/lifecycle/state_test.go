package lifecycle

import (
	"errors"
	"testing"

	"github.com/vinayakgaud/goflux/internal/goerror"
)

func TestStateStoreInitialState(t *testing.T) {
	storeState := newStateStore()

	if got := storeState.current(); got != serverStateStarting {
		t.Fatalf("got %v, want %v", got, serverStateStarting)
	}
}

func TestStateStoreTransition(t *testing.T) {
	storeState := newStateStore()

	if err := storeState.transition(serverStateRunning); err != nil {
		t.Fatal("expected transition to succedd")
	}

	if got := storeState.current(); got != serverStateRunning {
		t.Fatalf("got %v, want %v", got, serverStateRunning)
	}
}

func TestInvalidStateStoretransition(t *testing.T) {
	storeState := newStateStore()

	if err := storeState.transition(serverStateDraining); err == nil {
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

func TestInvalidStateTransitionError(t *testing.T) {
	storeState := newStateStore()

	err := storeState.transition(serverStateDraining)

	if err == nil {
		t.Fatal("expected error")
	}

	var goErr *goerror.Error

	if !errors.As(err, &goErr) {
		t.Fatal("expected goerror.Error")
	}

	if goErr.Code != goerror.CodeInvalidStateTransition {
		t.Fatalf(
			"got %v, want %v",
			goErr.Code,
			goerror.CodeInvalidStateTransition,
		)
	}
}
