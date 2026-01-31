package main

import (
	"fmt"
	"io"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"

// Player manages the audio playback for SomaFM streams.
type Player struct {
	ctx            *oto.Context
	player         *oto.Player
	stream         io.ReadCloser
	bufferedStream *BufferedStream
}

// NewPlayer initializes a new audio player with a default sample rate and channel count.
func NewPlayer() (*Player, error) {
	// Initialize oto context with standard audio parameters
	op := &oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}
	ctx, ready, err := oto.NewContext(op)
	if err != nil {
		return nil, fmt.Errorf("failed to create oto context: %w", err)
	}
	// Wait for the audio context to be ready
	<-ready

	return &Player{ctx: ctx}, nil
}

// Play starts streaming and playing audio from the given URL.
// It closes any previously playing stream before starting a new one.
// Returns a channel that emits buffer state updates.
func (p *Player) Play(url string) (<-chan BufferStats, error) {
	// Close any existing player and stream to prevent resource leaks
	if p.player != nil {
		p.player.Close()
		p.player = nil
	}
	if p.bufferedStream != nil {
		p.bufferedStream.Close()
		p.bufferedStream = nil
	}
	if p.stream != nil {
		p.stream.Close()
		p.stream = nil
	}

	// Create a new buffered stream
	bs := NewBufferedStream(url)

	// Start the buffered stream (performs initial connection)
	statsChan, err := bs.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start buffered stream: %w", err)
	}

	// Decode the MP3 stream from the buffered stream
	decodedStream, err := mp3.DecodeWithSampleRate(44100, bs)
	if err != nil {
		bs.Close()
		return nil, fmt.Errorf("failed to decode mp3: %w", err)
	}

	// Store the buffered stream and create a new player, then start playback
	p.bufferedStream = bs
	p.player = p.ctx.NewPlayer(decodedStream)
	p.player.Play()

	return statsChan, nil
}

// Stop halts the current audio playback and closes the associated stream.
func (p *Player) Stop() {
	if p.player != nil {
		p.player.Close()
		p.player = nil
	}
	if p.bufferedStream != nil {
		p.bufferedStream.Close()
		p.bufferedStream = nil
	}
	if p.stream != nil {
		p.stream.Close()
		p.stream = nil
	}
}
