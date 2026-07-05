package main

import (
	"strings"
	"testing"

	"somatui/internal/channels"
	"somatui/internal/protocol"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatChannelList_MarksFavoritesAndKeepsOrder(t *testing.T) {
	payload := protocol.ChannelsPayload{
		Channels: []channels.Channel{
			{ID: "dronezone", Title: "Drone Zone", Genre: "ambient"},
			{ID: "groovesalad", Title: "Groove Salad", Genre: "ambient|electronica"},
			{ID: "secretagent", Title: "Secret Agent", Genre: "lounge"},
		},
		Favorites: []string{"dronezone"},
	}

	out := formatChannelList(payload)
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	require.Len(t, lines, 3)

	// Catalog order is preserved (the server already sorts favorites first).
	assert.Contains(t, lines[0], "dronezone")
	assert.Contains(t, lines[1], "groovesalad")
	assert.Contains(t, lines[2], "secretagent")

	// Favorites are starred, everything else is not.
	assert.True(t, strings.HasPrefix(lines[0], "* "), "favorite must be marked: %q", lines[0])
	assert.True(t, strings.HasPrefix(lines[1], "  "), "non-favorite must not be marked: %q", lines[1])

	// Title and genre columns are present.
	assert.Contains(t, lines[0], "Drone Zone")
	assert.Contains(t, lines[0], "ambient")
	assert.Contains(t, lines[2], "Secret Agent")
	assert.Contains(t, lines[2], "lounge")
}

func TestFormatChannelList_ScriptFriendlyFields(t *testing.T) {
	payload := protocol.ChannelsPayload{
		Channels: []channels.Channel{
			{ID: "groovesalad", Title: "Groove Salad", Genre: "ambient"},
		},
	}

	// Scripts consume the list with awk-style field splitting, so the ID
	// must lead the line, preceded only by the star on favorites.
	fields := strings.Fields(formatChannelList(payload))
	require.NotEmpty(t, fields)
	assert.Equal(t, "groovesalad", fields[0], "unstarred lines lead with the ID")

	payload.Favorites = []string{"groovesalad"}
	fields = strings.Fields(formatChannelList(payload))
	require.GreaterOrEqual(t, len(fields), 2)
	assert.Equal(t, "*", fields[0])
	assert.Equal(t, "groovesalad", fields[1])
}

func TestFormatChannelList_EmptyCatalog(t *testing.T) {
	assert.Empty(t, formatChannelList(protocol.ChannelsPayload{}))
}
