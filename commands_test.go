package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadChannels_FromCache(t *testing.T) {
	setCacheDir(t)

	// Pre-populate the cache
	require.NoError(t, writeChannelsToCache(&testChannelData))

	msg := loadChannels()

	loaded, ok := msg.(channelsLoadedMsg)
	require.True(t, ok, "expected channelsLoadedMsg, got %T", msg)
	assert.True(t, loaded.fromCache)
	assert.Equal(t, 2, len(loaded.channels.Channels))
}

func TestLoadChannels_FromNetwork(t *testing.T) {
	setCacheDir(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(testChannelData)
		_, _ = w.Write(data)
	}))
	defer server.Close()

	originalURL := somafmChannelsURL
	somafmChannelsURL = server.URL
	t.Cleanup(func() { somafmChannelsURL = originalURL })

	msg := loadChannels()

	loaded, ok := msg.(channelsLoadedMsg)
	require.True(t, ok, "expected channelsLoadedMsg, got %T", msg)
	assert.False(t, loaded.fromCache)
	assert.Equal(t, 2, len(loaded.channels.Channels))
}

func TestLoadChannels_NetworkError(t *testing.T) {
	setCacheDir(t)

	originalURL := somafmChannelsURL
	somafmChannelsURL = "http://invalid-host.example.com"
	t.Cleanup(func() { somafmChannelsURL = originalURL })

	msg := loadChannels()

	_, ok := msg.(errorMsg)
	assert.True(t, ok, "expected errorMsg, got %T", msg)
}

func TestRefreshChannels_Success(t *testing.T) {
	setCacheDir(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(testChannelData)
		_, _ = w.Write(data)
	}))
	defer server.Close()

	originalURL := somafmChannelsURL
	somafmChannelsURL = server.URL
	t.Cleanup(func() { somafmChannelsURL = originalURL })

	msg := refreshChannels()

	refreshed, ok := msg.(channelsRefreshedMsg)
	require.True(t, ok, "expected channelsRefreshedMsg, got %T", msg)
	assert.Equal(t, 2, len(refreshed.channels.Channels))
}

func TestRefreshChannels_Error(t *testing.T) {
	setCacheDir(t)

	originalURL := somafmChannelsURL
	somafmChannelsURL = "http://invalid-host.example.com"
	t.Cleanup(func() { somafmChannelsURL = originalURL })

	msg := refreshChannels()
	// refreshChannels silently ignores errors and returns nil
	assert.Nil(t, msg)
}

func TestTickChannelRefresh(t *testing.T) {
	cmd := tickChannelRefresh()
	assert.NotNil(t, cmd)
}
