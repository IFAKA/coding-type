package history

import (
	"math"
	"sort"
	"time"
)

// Stats holds computed aggregate statistics over all history entries.
type Stats struct {
	BestWPM      int
	AvgWPM       int
	AvgAccuracy  float64
	TotalSessions int
	Streak       int // consecutive calendar days with ≥1 session
	Recent       []Entry // last 10 entries, newest first
}

// Compute derives Stats from a slice of history entries.
func Compute(entries []Entry) Stats {
	if len(entries) == 0 {
		return Stats{}
	}

	// Sort newest first for recents
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.After(sorted[j].Timestamp)
	})

	best := 0
	totalWPM := 0
	totalAcc := 0.0
	for _, e := range entries {
		if e.WPM > best {
			best = e.WPM
		}
		totalWPM += e.WPM
		totalAcc += e.Accuracy
	}

	n := len(entries)
	avgWPM := int(math.Round(float64(totalWPM) / float64(n)))
	avgAcc := totalAcc / float64(n)

	recent := sorted
	if len(recent) > 10 {
		recent = recent[:10]
	}

	return Stats{
		BestWPM:       best,
		AvgWPM:        avgWPM,
		AvgAccuracy:   avgAcc,
		TotalSessions: n,
		Streak:        computeStreak(sorted),
		Recent:        recent,
	}
}

// computeStreak counts consecutive calendar days (going back from today)
// that had at least one session, based on entries sorted newest-first.
func computeStreak(sorted []Entry) int {
	if len(sorted) == 0 {
		return 0
	}

	// Build set of days that have sessions
	days := make(map[string]bool)
	for _, e := range sorted {
		day := e.Timestamp.Format("2006-01-02")
		days[day] = true
	}

	streak := 0
	today := time.Now()
	for {
		day := today.Format("2006-01-02")
		if !days[day] {
			break
		}
		streak++
		today = today.AddDate(0, 0, -1)
	}
	return streak
}

// AvgWPMForLanguage returns the average WPM for a specific language.
func AvgWPMForLanguage(entries []Entry, language string) int {
	total, count := 0, 0
	for _, e := range entries {
		if e.Language == language {
			total += e.WPM
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return int(math.Round(float64(total) / float64(count)))
}
