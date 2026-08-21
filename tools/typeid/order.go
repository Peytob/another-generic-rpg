package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

// OrderedTypeId is the input of the generator: the target package with its
// sources and the ordered type names.
type OrderedTypeId struct {
	// ImportPath is the full import path of the package; together with a
	// type name it forms the hashed identity.
	ImportPath string
	// Dir is the package directory on disk.
	Dir string
	// Files is the sorted list of the package file names.
	Files []string
	// Sources maps a file name to its current content.
	Sources map[string]string
	// TypeNames lists the types whose IDs are ordered.
	TypeNames []string
}

// ReadOrder reads the generator input: the -type flag, the $GOFILE
// environment variable set by go generate, the package information obtained
// via go list and the package sources. Without -type all struct types
// declared in $GOFILE are ordered.
func ReadOrder() (OrderedTypeId, error) {
	typesFlag := flag.String("type", "", "comma-separated list of type names; defaults to all struct types in $GOFILE")

	flag.Parse()
	if flag.NArg() != 0 {
		return OrderedTypeId{}, fmt.Errorf("unexpected arguments: %s", strings.Join(flag.Args(), " "))
	}

	pkg, err := listPackage()
	if err != nil {
		return OrderedTypeId{}, err
	}

	goFile := os.Getenv("GOFILE")
	order, err := readSources(pkg, goFile)
	if err != nil {
		return OrderedTypeId{}, err
	}

	if *typesFlag != "" {
		order.TypeNames = splitTypes(*typesFlag)
		return order, nil
	}

	if goFile == "" {
		return order, fmt.Errorf("no -type flag given and $GOFILE is not set; run via go generate or pass -type")
	}

	names, err := structTypeNames(order.Sources[goFile])
	if err != nil {
		return order, fmt.Errorf("parse %s: %w", goFile, err)
	}

	if len(names) == 0 {
		return order, fmt.Errorf("no struct types in %s: %w", goFile, errNothingToDo)
	}

	order.TypeNames = names

	return order, nil
}

// readSources builds the order skeleton from the package info, making sure
// the go:generate file is part of the file set, and loads the file contents.
func readSources(pkg *packageInfo, goFile string) (OrderedTypeId, error) {
	order := OrderedTypeId{
		ImportPath: pkg.ImportPath,
		Dir:        pkg.Dir,
		Files:      slices.Clone(pkg.GoFiles),
	}
	if goFile != "" && !slices.Contains(order.Files, goFile) {
		order.Files = append(order.Files, goFile)
	}
	sort.Strings(order.Files)

	order.Sources = make(map[string]string, len(order.Files))
	for _, f := range order.Files {
		b, err := os.ReadFile(filepath.Join(pkg.Dir, f))
		if err != nil {
			return OrderedTypeId{}, err
		}
		order.Sources[f] = string(b)
	}
	return order, nil
}

func splitTypes(s string) []string {
	var names []string
	for _, t := range strings.Split(s, ",") {
		if t = strings.TrimSpace(t); t != "" && !slices.Contains(names, t) {
			names = append(names, t)
		}
	}
	return names
}

// structTypeNames returns the names of the non-generic struct types declared
// in the given source.
func structTypeNames(source string) ([]string, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "src.go", source, 0)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, s := range gd.Specs {
			ts, ok := s.(*ast.TypeSpec)
			if !ok || ts.TypeParams != nil {
				continue
			}
			if _, ok := ts.Type.(*ast.StructType); ok {
				names = append(names, ts.Name.Name)
			}
		}
	}
	return names, nil
}

type packageInfo struct {
	ImportPath string
	Dir        string
	Name       string
	GoFiles    []string
}

func listPackage() (*packageInfo, error) {
	cmd := exec.Command("go", "list", "-json")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("go list: %v: %s", err, strings.TrimSpace(stderr.String()))
	}
	var pkg packageInfo
	if err := json.Unmarshal(stdout.Bytes(), &pkg); err != nil {
		return nil, fmt.Errorf("go list: decode output: %w", err)
	}
	if pkg.ImportPath == "" || pkg.Dir == "" {
		return nil, fmt.Errorf("go list: empty package info")
	}
	return &pkg, nil
}
