package loaders

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env.test")
	content := "# comment\nKV1=aaa\nKV2=bbb\n"
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	m, err := ReadFile(p)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if m["KV1"] != "aaa" || m["KV2"] != "bbb" {
		t.Fatalf("map mismatch: %#v", m)
	}
}
