package walletapp

import (
	"context"
	"testing"
)

func TestNewWalletApp(t *testing.T) {
	app := NewWalletApp()
	ctx := context.Background()

	app.Startup(ctx)

	if app.Ctx == nil {
		t.Error("Expected non-nil context, got nil")
	}

	if app.Shutdown {
		t.Error("Expected Shutdown to be false, got true")
	}

	if app.IsListening {
		t.Error("Expected IsListening to be false, got true")
	}
}
