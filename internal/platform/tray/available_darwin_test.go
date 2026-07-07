//go:build darwin

package tray

import "testing"

// TestAvailableDoesNotPanic exercises the CGSession query. Its result depends
// on whether the test host has a GUI session, so only its safety is asserted.
func TestAvailableDoesNotPanic(t *testing.T) {
	_ = Available()
}
