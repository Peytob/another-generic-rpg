package main

import (
	"log/slog"
	"os"
	"path/filepath"
)

// WriteResults stores the generation results: it writes every changed file
// back to the package directory and logs a summary.
func WriteResults(order OrderedTypeId, generated []GeneratedTypeId, changed map[string]string) error {
	for _, file := range sortedKeys(changed) {
		path := filepath.Join(order.Dir, file)
		if err := writeFile(path, changed[file]); err != nil {
			return err
		}
		slog.Info("typeid: updated", "path", path)
	}
	if len(changed) == 0 && len(generated) > 0 {
		slog.Info("typeid: all type IDs are up to date")
	}
	return nil
}

func writeFile(path, content string) error {
	mode := os.FileMode(0o644)
	if fi, err := os.Stat(path); err == nil {
		mode = fi.Mode().Perm()
	}
	return os.WriteFile(path, []byte(content), mode)
}
