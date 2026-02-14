package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

type samplePayload struct {
	Name string `json:"name"`
	Val  int    `json:"val"`
}

func TestWriteReadPayload(t *testing.T) {
	dir := t.TempDir()
	fc := New(dir)

	payload := samplePayload{Name: "alpha", Val: 42}
	meta := map[string]string{"source": "unit"}

	entry, err := fc.Write("sample-key", payload, meta)
	if err != nil {
		t.Fatalf("write: %v", err)
	}
	if entry.FetchedAt == "" {
		t.Fatalf("fetchedAt is empty")
	}

	var out samplePayload
	readEntry, ok, stale, err := fc.ReadPayload("sample-key", time.Minute, &out)
	if err != nil {
		t.Fatalf("read payload: %v", err)
	}
	if !ok {
		t.Fatalf("cache not found")
	}
	if stale {
		t.Fatalf("cache is stale")
	}
	if out != payload {
		t.Fatalf("payload mismatch")
	}
	if readEntry.FetchedAt == "" {
		t.Fatalf("read fetchedAt is empty")
	}
}

func TestReadStaleByTTL(t *testing.T) {
	dir := t.TempDir()
	fc := New(dir)

	payload := samplePayload{Name: "beta", Val: 7}
	if _, err := fc.Write("ttl-key", payload, nil); err != nil {
		t.Fatalf("write: %v", err)
	}

	_, ok, stale, err := fc.Read("ttl-key", -1)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !ok {
		t.Fatalf("cache not found")
	}
	if stale {
		t.Fatalf("ttl <= 0 should not be stale")
	}

	_, ok, stale, err = fc.Read("ttl-key", 0)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !ok {
		t.Fatalf("cache not found")
	}
	if stale {
		t.Fatalf("ttl = 0 should not be stale")
	}
}

func TestDelete(t *testing.T) {
	dir := t.TempDir()
	fc := New(dir)

	if _, err := fc.Write("delete-key", samplePayload{Name: "x", Val: 1}, nil); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := fc.Delete("delete-key"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, ok, _, err := fc.Read("delete-key", time.Minute)
	if err != nil {
		t.Fatalf("read after delete: %v", err)
	}
	if ok {
		t.Fatalf("cache still exists")
	}
}

func TestWriteCreatesSafeFileName(t *testing.T) {
	dir := t.TempDir()
	fc := New(dir)

	key := "unsafe key/with:chars"
	if _, err := fc.Write(key, samplePayload{Name: "safe", Val: 9}, nil); err != nil {
		t.Fatalf("write: %v", err)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("unexpected file count: %d", len(files))
	}

	name := files[0].Name()
	if filepath.Ext(name) != ".json" {
		t.Fatalf("unexpected file extension: %s", name)
	}
}
