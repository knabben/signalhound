package tui

import (
	"os"
)

// Helper function to detect Wayland
// Wayland is a display server protocol that is intended to replace the
// X Window System (X11) on Linux and other Unix-like operating systems.
func isWayland() bool {
	// Check common Wayland environment variables
	waylandDisplay := os.Getenv("WAYLAND_DISPLAY")
	xdgSessionType := os.Getenv("XDG_SESSION_TYPE")
	return waylandDisplay != "" || xdgSessionType == "wayland"
}

// Helper function to detect WSL
// WSL (Windows Subsystem for Linux) is a compatibility layer for running
// Linux binary executables natively on Windows.
func isWSL() bool {
	_, exists := os.LookupEnv("WSL_DISTRO_NAME")
	return exists
}
