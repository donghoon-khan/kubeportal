package integration

import "testing"

func TestNewIntegrationManager(t *testing.T) {
	iManager := NewIntegrationManager(nil)
	if iManager == nil {
		t.Errorf("Failed to create integration manager.")
	}
}
