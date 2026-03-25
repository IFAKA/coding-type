package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Entry records the result of one completed typing exercise.
type Entry struct {
	Timestamp  time.Time `json:"timestamp"`
	Language   string    `json:"language"`
	SnippetID  string    `json:"snippet_id"`
	SnippetTitle string  `json:"snippet_title"`
	WPM        int       `json:"wpm"`
	Accuracy   float64   `json:"accuracy"`
	DurationMs int64     `json:"duration_ms"`
	Errors     int       `json:"errors"`
}

type store struct {
	Entries []Entry `json:"entries"`
}

func configDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		home, err2 := os.UserHomeDir()
		if err2 != nil {
			return "", err2
		}
		base = filepath.Join(home, ".config")
	}
	dir := filepath.Join(base, "coding-type")
	return dir, os.MkdirAll(dir, 0755)
}

func storePath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "history.json"), nil
}

// Load reads all history entries from disk. Returns empty slice if no file exists.
func Load() ([]Entry, error) {
	path, err := storePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var s store
	if err := json.Unmarshal(data, &s); err != nil {
		return []Entry{}, nil
	}
	return s.Entries, nil
}

// Save appends a new entry to the history file.
func Save(e Entry) error {
	entries, err := Load()
	if err != nil {
		entries = []Entry{}
	}
	entries = append(entries, e)
	// Keep last 1000 entries max
	if len(entries) > 1000 {
		entries = entries[len(entries)-1000:]
	}
	data, err := json.MarshalIndent(store{Entries: entries}, "", "  ")
	if err != nil {
		return err
	}
	path, err := storePath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LastSeenMap returns a map of snippetID → last time it was played.
func LastSeenMap(entries []Entry) map[string]time.Time {
	m := make(map[string]time.Time)
	for _, e := range entries {
		if last, ok := m[e.SnippetID]; !ok || e.Timestamp.After(last) {
			m[e.SnippetID] = e.Timestamp
		}
	}
	return m
}
