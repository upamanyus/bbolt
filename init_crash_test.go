package bbolt_test

import (
	"os"
	"path/filepath"
	"testing"

	bolt "go.etcd.io/bbolt"
)

// An init interrupted before its meta pages reached the disk leaves the other pages and zero
// meta pages (init writes and syncs them first). Open must initialize such a file again rather
// than refuse it forever.
func TestOpen_InterruptedInit(t *testing.T) {
	for _, size := range []int{4096, 4 * 4096} {
		path := filepath.Join(t.TempDir(), "db")
		if err := os.WriteFile(path, make([]byte, size), 0o600); err != nil {
			t.Fatal(err)
		}
		db, err := bolt.Open(path, 0o600, &bolt.Options{PageSize: 4096})
		if err != nil {
			t.Fatalf("size %d: Open: %v", size, err)
		}
		if err := db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucket([]byte("b"))
			if err != nil {
				return err
			}
			return b.Put([]byte("k"), []byte("v"))
		}); err != nil {
			t.Fatalf("size %d: Update: %v", size, err)
		}
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
		db, err = bolt.Open(path, 0o600, nil)
		if err != nil {
			t.Fatalf("size %d: reopen: %v", size, err)
		}
		if err := db.View(func(tx *bolt.Tx) error {
			if v := tx.Bucket([]byte("b")).Get([]byte("k")); string(v) != "v" {
				t.Fatalf("size %d: got %q", size, v)
			}
			return nil
		}); err != nil {
			t.Fatal(err)
		}
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
