//go:build linux

package tray

import "os"

// Available reports whether a graphical environment is present. On a headless
// system (no display and no session bus) it returns false so the server runs
// without a tray instead of spinning up a tray that can never show. When a
// display or session bus exists the tray still degrades gracefully if no
// StatusNotifierHost is watching.
func Available() bool {
	return os.Getenv("DISPLAY") != "" ||
		os.Getenv("WAYLAND_DISPLAY") != "" ||
		os.Getenv("DBUS_SESSION_BUS_ADDRESS") != ""
}
