package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/oto/v2"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

// Player manages the audio playback for SomaFM streams.
type Player struct {
	ctx    *oto.Context
	player oto.Player
	stream io.ReadCloser
}

// NewPlayer initializes a new audio player with a default sample rate and channel count.
func NewPlayer() (*Player, error) {
	// Initialize oto context with standard audio parameters
	ctx, ready, err := oto.NewContext(44100, 2, oto.FormatSignedInt16LE)
	if err != nil {
		return nil, fmt.Errorf("failed to create oto context: %w", err)
	}
	// Wait for the audio context to be ready
	<-ready

	return &Player{ctx: ctx}, nil
}

// Play starts streaming and playing audio from the given URL.
// It closes any previously playing stream before starting a new one.
func (p *Player) Play(url string) error {
	// Close any existing player and stream to prevent resource leaks
	if p.player != nil {
		p.player.Close()
	}
	if p.stream != nil {
		p.stream.Close()
	}

	// Create a new HTTP request and set a User-Agent header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	// Execute the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch stream: %w", err)
	}

	// Check for successful HTTP status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode the MP3 stream from the response body
	decodedStream, err := mp3.DecodeWithSampleRate(44100, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to decode mp3: %w", err)
	}

	// Store the stream and create a new player, then start playback
	p.stream = resp.Body
	p.player = p.ctx.NewPlayer(decodedStream)
	p.player.Play()

	return nil
}

// Stop halts the current audio playback and closes the associated stream.
func (p *Player) Stop() {
	if p.player != nil {
		p.player.Close()
	}
	if p.stream != nil {
		p.stream.Close()
	}
}
