// Command typeid generates stable numeric type identifiers for named Go
// types. The identifier of type T declared in a package with import path P is
// the FNV-1a 64-bit hash of "P.T" reinterpreted as int64.
//
// Usage in a package file (the typeid module must be listed in the workspace
// go.work):
//
//	//go:generate go run typeid
//	//go:generate go run typeid -type Foo,Bar
//
// Without -type, all struct types declared in $GOFILE (the file containing
// the go:generate directive) are processed. For every type T the generator
// maintains the constant
//
//	const TTypeID int64 = <id>
//
// in the same file right after the type declaration, or updates the existing
// constant in place wherever in the package it is declared. Repeated runs are
// idempotent.
//
// The implementation is split into three stages: ReadOrder reads the input,
// Generate computes the identifiers and WriteResults stores them.
package main

import (
	"errors"
	"log/slog"
	"os"
)

// errNothingToDo signals that the order contains no work; it is not a
// failure.
var errNothingToDo = errors.New("nothing to do")

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	if err := run(); err != nil {
		slog.Error("typeid: failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	order, err := ReadOrder()
	if err != nil {
		if errors.Is(err, errNothingToDo) {
			slog.Info("typeid: nothing to do", "detail", err.Error())
			return nil
		}
		return err
	}

	generated, changed, err := Generate(order)
	if err != nil {
		return err
	}

	return WriteResults(order, generated, changed)
}
