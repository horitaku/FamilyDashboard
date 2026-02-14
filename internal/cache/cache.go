package cache

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Entry はキャッシュファイルの中身を表すのです。
type Entry struct {
	Payload   json.RawMessage   `json:"payload"`
	FetchedAt string            `json:"fetchedAt"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// FileCache はJSONファイルキャッシュを扱うのです。
type FileCache struct {
	dir   string
	clock func() time.Time
}

// New はキャッシュ管理者をつくるのです。
func New(dir string) *FileCache {
	return &FileCache{
		dir:   dir,
		clock: nowTokyo,
	}
}

// Write はペイロードを保存して、保存したEntryを返すのです。
func (fc *FileCache) Write(key string, payload any, meta map[string]string) (Entry, error) {
	if fc == nil {
		return Entry{}, errors.New("cache is nil")
	}

	if err := os.MkdirAll(fc.dir, 0o755); err != nil {
		return Entry{}, err
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Entry{}, err
	}

	entry := Entry{
		Payload:   payloadBytes,
		FetchedAt: fc.clock().Format(time.RFC3339),
		Meta:      meta,
	}

	entryBytes, err := json.Marshal(entry)
	if err != nil {
		return Entry{}, err
	}

	path := fc.filePath(key)
	tmpFile, err := os.CreateTemp(fc.dir, safeFilePrefix(key))
	if err != nil {
		return Entry{}, err
	}

	if _, err := tmpFile.Write(entryBytes); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
		return Entry{}, err
	}

	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return Entry{}, err
	}

	if err := os.Rename(tmpFile.Name(), path); err != nil {
		_ = os.Remove(tmpFile.Name())
		return Entry{}, err
	}

	return entry, nil
}

// Read はキャッシュを読み取り、存在・期限切れの状態も返すのです。
func (fc *FileCache) Read(key string, ttl time.Duration) (Entry, bool, bool, error) {
	if fc == nil {
		return Entry{}, false, false, errors.New("cache is nil")
	}

	path := fc.filePath(key)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Entry{}, false, false, nil
		}
		return Entry{}, false, false, err
	}

	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return entry, true, true, err
	}

	fetchedAt, err := time.Parse(time.RFC3339, entry.FetchedAt)
	if err != nil {
		return entry, true, true, err
	}

	if ttl <= 0 {
		return entry, true, false, nil
	}

	stale := fc.clock().Sub(fetchedAt) > ttl
	return entry, true, stale, nil
}

// ReadPayload はキャッシュを読み取り、payloadを型に詰めるのです。
func (fc *FileCache) ReadPayload(key string, ttl time.Duration, out any) (Entry, bool, bool, error) {
	entry, ok, stale, err := fc.Read(key, ttl)
	if !ok || err != nil {
		return entry, ok, stale, err
	}

	if out == nil {
		return entry, ok, stale, errors.New("output is nil")
	}

	if err := json.Unmarshal(entry.Payload, out); err != nil {
		return entry, ok, true, err
	}

	return entry, ok, stale, nil
}

// Delete は指定キーのキャッシュを削除するのです。
func (fc *FileCache) Delete(key string) error {
	if fc == nil {
		return errors.New("cache is nil")
	}

	path := fc.filePath(key)
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}

func (fc *FileCache) filePath(key string) string {
	return filepath.Join(fc.dir, safeFileName(key)+".json")
}

func safeFilePrefix(key string) string {
	return safeFileName(key) + "-*.tmp"
}

func safeFileName(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return "cache"
	}

	var builder strings.Builder
	for _, r := range key {
		if r >= 'a' && r <= 'z' {
			builder.WriteRune(r)
			continue
		}
		if r >= 'A' && r <= 'Z' {
			builder.WriteRune(r)
			continue
		}
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
			continue
		}
		if r == '-' || r == '_' {
			builder.WriteRune(r)
			continue
		}
		builder.WriteRune('_')
	}

	clean := builder.String()
	if clean == "" {
		return "cache"
	}
	return clean
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
