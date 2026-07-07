//go:build linux

package tray

import "testing"

func TestAvailableUsesEnv(t *testing.T) {
	for _, v := range []string{"DISPLAY", "WAYLAND_DISPLAY", "DBUS_SESSION_BUS_ADDRESS"} {
		t.Setenv(v, "")
	}
	if Available() {
		t.Fatal("Available() = true with no display or session bus, want false")
	}

	t.Setenv("WAYLAND_DISPLAY", "wayland-0")
	if !Available() {
		t.Fatal("Available() = false with WAYLAND_DISPLAY set, want true")
	}
}
