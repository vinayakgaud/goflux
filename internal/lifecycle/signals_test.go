package lifecycle

import (
	"context"
	"testing"
	"time"
)

func TestSignalContextParentCancellation(t *testing.T) {
	parent, cancelParent := context.WithCancel(context.Background())
	defer cancelParent()

	ctx, cancel := signalContext(parent)
	defer cancel()

	cancelParent()

	select {
	case <-ctx.Done():
		//expected
	case <-time.After(time.Second):
		t.Fatal("expected signal context to be cancelled")
	}
}
