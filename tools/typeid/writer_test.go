package main

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestMain silences the default logger so the tests do not pollute go test
// output.
func TestMain(m *testing.M) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	os.Exit(m.Run())
}

func TestWriteResultsWritesChangedFiles(t *testing.T) {
	dir := t.TempDir()
	order := OrderedTypeId{Dir: dir}
	generated := []GeneratedTypeId{{TypeName: "A", ConstName: "ATypeID", ID: 42, File: "a.go"}}
	changed := map[string]string{"a.go": "package p\n"}

	if err := WriteResults(order, generated, changed); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "a.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "package p\n" {
		t.Errorf("a.go = %q, want %q", b, "package p\n")
	}
}

func TestWriteResultsPreservesFileMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission bits are not fully supported on Windows")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	if err := os.WriteFile(path, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := WriteResults(OrderedTypeId{Dir: dir}, nil, map[string]string{"a.go": "new"}); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Errorf("mode = %v, want 0600", fi.Mode().Perm())
	}
}
