//go:build !linux

package main

import tea "github.com/charmbracelet/bubbletea"

// MPRISCmdSender is an interface for sending commands to the application.
type MPRISCmdSender interface {
	Send(msg tea.Msg)
}

// MPRIS is a stub for non-Linux platforms.
type MPRIS struct{}

// NewMPRIS returns nil on non-Linux platforms (MPRIS not supported).
func NewMPRIS() (*MPRIS, error) {
	return nil, nil
}

// SetSender is a no-op on non-Linux platforms.
func (m *MPRIS) SetSender(sender MPRISCmdSender) {}

// SetPlaying is a no-op on non-Linux platforms.
func (m *MPRIS) SetPlaying(station, track, artist string) {}

// SetStopped is a no-op on non-Linux platforms.
func (m *MPRIS) SetStopped() {}

// SetMetadata is a no-op on non-Linux platforms.
func (m *MPRIS) SetMetadata(station, track, artist string) {}

// Close is a no-op on non-Linux platforms.
func (m *MPRIS) Close() {}
