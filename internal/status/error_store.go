package status

import (
	"sort"
	"sync"
	"time"

	"github.com/rihow/FamilyDashboard/internal/models"
)

// ErrorStore tracks API fetch failures.
type ErrorStore struct {
	mu     sync.RWMutex
	errors map[string]models.ErrorInfo
	clock  func() time.Time
}

// NewErrorStore creates a new error store.
func NewErrorStore() *ErrorStore {
	return &ErrorStore{
		errors: map[string]models.ErrorInfo{},
		clock:  nowTokyo,
	}
}

// Set records an error for a source.
func (s *ErrorStore) Set(source, message string) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.errors[source] = models.ErrorInfo{
		Source:  source,
		Message: message,
		At:      s.clock().Format(time.RFC3339),
	}
}

// Clear clears the error for a source.
func (s *ErrorStore) Clear(source string) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.errors, source)
}

// List returns errors sorted by source for stable output.
func (s *ErrorStore) List() []models.ErrorInfo {
	if s == nil {
		return []models.ErrorInfo{}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.errors))
	for key := range s.errors {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := make([]models.ErrorInfo, 0, len(keys))
	for _, key := range keys {
		result = append(result, s.errors[key])
	}

	return result
}

// NowRFC3339 returns current time in Asia/Tokyo.
func NowRFC3339() string {
	return nowTokyo().Format(time.RFC3339)
}

var tokyoLocationOnce sync.Once
var tokyoLocation *time.Location

func nowTokyo() time.Time {
	tokyoLocationOnce.Do(func() {
		loc, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			tokyoLocation = time.UTC
			return
		}
		tokyoLocation = loc
	})

	return time.Now().In(tokyoLocation)
}
