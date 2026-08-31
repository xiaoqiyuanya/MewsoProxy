package webui

import (
	"io/fs"
	"testing"
)

func TestDistEmbedded(t *testing.T) {
	sub, err := fs.Sub(Dist, "dist")
	if err != nil {
		t.Fatalf("sub dist: %v", err)
	}
	if _, err := fs.Stat(sub, "index.html"); err != nil {
		t.Fatalf("index.html missing: %v", err)
	}
	assets, err := fs.Sub(Dist, "dist/assets")
	if err != nil {
		t.Fatalf("sub assets: %v", err)
	}
	entries, err := fs.ReadDir(assets, ".")
	if err != nil {
		t.Fatalf("read assets: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("assets dir is empty")
	}
}
