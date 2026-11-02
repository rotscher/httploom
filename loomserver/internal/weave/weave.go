// Package weave contains the domain logic for the Loom API.
package weave

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"sync"
	"time"

	"github.com/rotscher/httploom/loomserver/internal/api"
)

// DefaultWidth is used when a request does not specify a width.
const DefaultWidth = 80

// Wrap breaks text into lines of at most width characters without splitting
// words. Existing newlines are preserved as line breaks.
func Wrap(text string, width int) string {
	if width < 1 {
		width = DefaultWidth
	}

	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	wrapped := make([]string, 0, len(lines))
	for _, line := range lines {
		words := strings.Fields(line)
		if len(words) == 0 {
			wrapped = append(wrapped, "")
			continue
		}

		var b strings.Builder
		col := 0
		for i, word := range words {
			if i == 0 {
				b.WriteString(word)
				col = len(word)
				continue
			}
			if col+1+len(word) > width {
				b.WriteString("\n")
				b.WriteString(word)
				col = len(word)
				continue
			}
			b.WriteString(" ")
			b.WriteString(word)
			col += 1 + len(word)
		}
		wrapped = append(wrapped, b.String())
	}

	return strings.Join(wrapped, "\n")
}

// Store is an in-memory collection of weaves. It is safe for concurrent use.
type Store struct {
	mu    sync.RWMutex
	items map[string]api.Weave
	order []string
}

// NewStore returns an empty Store.
func NewStore() *Store {
	return &Store{items: make(map[string]api.Weave)}
}

// Create weaves text to width, stores it and returns the resulting record.
func (s *Store) Create(text string, width int) api.Weave {
	w := api.Weave{
		Id:        newID(),
		Text:      Wrap(text, width),
		Width:     width,
		CreatedAt: time.Now().UTC(),
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[w.Id] = w
	s.order = append(s.order, w.Id)

	return w
}

// Get returns the weave with the given id.
func (s *Store) Get(id string) (api.Weave, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.items[id]
	return w, ok
}

// List returns all weaves in creation order.
func (s *Store) List() []api.Weave {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]api.Weave, 0, len(s.order))
	for _, id := range s.order {
		items = append(items, s.items[id])
	}
	return items
}

func newID() string {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand failing is not expected on supported platforms; fall
		// back to a timestamp-based id rather than failing the request.
		return hex.EncodeToString([]byte(time.Now().UTC().Format("150405.000000")))
	}
	return hex.EncodeToString(b[:])
}
