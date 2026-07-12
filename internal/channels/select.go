package channels

// mp3QualityRank orders SomaFM playlist quality levels, best first.
var mp3QualityRank = map[string]int{"highest": 0, "high": 1, "low": 2}

// SelectMP3PlaylistURL returns the best-quality MP3 playlist URL from a
// channel's playlists (highest > high > low > unknown), or "" if none.
func SelectMP3PlaylistURL(playlists []Playlist) string {
	bestURL := ""
	// The seed must exceed the unknown-quality rank below, or a channel
	// whose playlists all have unrecognized quality labels would select
	// nothing at all instead of falling back to its first entry.
	bestRank := len(mp3QualityRank) + 1
	for _, playlist := range playlists {
		if playlist.Format != "mp3" {
			continue
		}
		rank, ok := mp3QualityRank[playlist.Quality]
		if !ok {
			rank = len(mp3QualityRank)
		}
		if rank < bestRank {
			bestURL = playlist.URL
			bestRank = rank
		}
	}
	return bestURL
}
