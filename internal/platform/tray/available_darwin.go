//go:build darwin

package tray

/*
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation
#include <CoreGraphics/CoreGraphics.h>
#include <CoreFoundation/CoreFoundation.h>

// hasGUISession reports whether the caller is running in a Quartz GUI session.
// CGSessionCopyCurrentDictionary is safe to call from any context (including
// daemons and SSH sessions) and returns NULL when there is no window server
// session, so it never itself pulls in AppKit or touches the menu bar.
static int hasGUISession(void) {
    CFDictionaryRef d = CGSessionCopyCurrentDictionary();
    if (d == NULL) {
        return 0;
    }
    CFRelease(d);
    return 1;
}
*/
import "C"

// Available reports whether a GUI session is present. It is false on a headless
// Mac (no console login), where starting the tray would attach AppKit to a
// missing window server and could hang or abort — so the server runs without
// a tray instead.
func Available() bool {
	return C.hasGUISession() != 0
}
