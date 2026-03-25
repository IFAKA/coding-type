package snippets

import (
	"embed"
	"encoding/json"
	"math/rand"
	"sort"
	"time"
)

//go:embed data/*.json
var dataFS embed.FS

var allSnippets []Snippet

func init() {
	files := []string{
		"data/python.json",
		"data/javascript.json",
		"data/typescript.json",
		"data/go.json",
		"data/cpp.json",
		"data/nextjs.json",
	}
	for _, f := range files {
		b, err := dataFS.ReadFile(f)
		if err != nil {
			continue
		}
		var s []Snippet
		if err := json.Unmarshal(b, &s); err != nil {
			continue
		}
		allSnippets = append(allSnippets, s...)
	}
}

// All returns all loaded snippets.
func All() []Snippet {
	return allSnippets
}

// Filter returns snippets matching the given language and difficulty.
// Empty string matches all values for that field.
func Filter(language, difficulty string) []Snippet {
	var result []Snippet
	for _, s := range allSnippets {
		if language != "" && s.Language != language {
			continue
		}
		if difficulty != "" && s.Difficulty != difficulty {
			continue
		}
		result = append(result, s)
	}
	return result
}

// Pick selects a snippet weighted by recency (snippets seen longer ago score higher).
// seenAt maps snippet ID → last time it was played (zero value = never played).
func Pick(language, difficulty string, seenAt map[string]time.Time) *Snippet {
	pool := Filter(language, difficulty)
	if len(pool) == 0 {
		// Fallback: ignore difficulty filter
		pool = Filter(language, "")
	}
	if len(pool) == 0 {
		return nil
	}

	now := time.Now()
	type scored struct {
		s     Snippet
		score float64
	}

	scored_ := make([]scored, len(pool))
	for i, s := range pool {
		last, ok := seenAt[s.ID]
		if !ok {
			// Never seen — maximum score
			scored_[i] = scored{s, 1e9}
		} else {
			scored_[i] = scored{s, now.Sub(last).Hours()}
		}
	}

	// Sort descending by score
	sort.Slice(scored_, func(i, j int) bool {
		return scored_[i].score > scored_[j].score
	})

	// Weighted random pick from top half (or at least 3)
	topN := len(scored_) / 2
	if topN < 3 {
		topN = len(scored_)
	}

	// Assign weights: higher score = higher weight
	totalWeight := 0.0
	for i := 0; i < topN; i++ {
		totalWeight += scored_[i].score
	}

	r := rand.Float64() * totalWeight
	cum := 0.0
	for i := 0; i < topN; i++ {
		cum += scored_[i].score
		if r <= cum {
			return &scored_[i].s
		}
	}
	return &scored_[0].s
}
