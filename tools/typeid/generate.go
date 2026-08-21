package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"hash/fnv"
	"io"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// GeneratedTypeId is the result of processing a single ordered type.
type GeneratedTypeId struct {
	// TypeName is the ordered type name.
	TypeName string
	// ConstName is the name of the maintained constant, TypeName + "TypeID".
	ConstName string
	// ID is the FNV-1a 64-bit hash of ImportPath.TypeName as int64.
	ID int64
	// File is the package file that declares the type.
	File string
}

// Generate runs the main logic: for every ordered type of the order it
// computes the identifier and collects the edits that keep the corresponding
// const declaration up to date. It returns the generated identifiers
// together with the updated contents of every changed file.
func Generate(order OrderedTypeId) ([]GeneratedTypeId, map[string]string, error) {
	if len(order.TypeNames) == 0 {
		return nil, nil, nil
	}

	g, err := newGenerator(order.ImportPath, order.Sources)
	if err != nil {
		return nil, nil, err
	}

	generated := make([]GeneratedTypeId, 0, len(order.TypeNames))
	for _, typeName := range order.TypeNames {
		id, err := g.ensureType(typeName)
		if err != nil {
			return nil, nil, err
		}
		generated = append(generated, id)
	}

	changed, err := g.apply()
	if err != nil {
		return nil, nil, err
	}

	return generated, changed, nil
}

const typeIDSuffix = "TypeID"

type edit struct {
	start int
	end   int
	text  string
}

type generator struct {
	importPath string
	fset       *token.FileSet
	names      []string
	files      map[string]*ast.File
	sources    map[string]string
	edits      map[string][]edit
}

func newGenerator(importPath string, sources map[string]string) (*generator, error) {
	g := &generator{
		importPath: importPath,
		fset:       token.NewFileSet(),
		files:      make(map[string]*ast.File, len(sources)),
		sources:    sources,
		edits:      make(map[string][]edit),
	}

	for _, name := range sortedKeys(sources) {
		f, err := parser.ParseFile(g.fset, name, sources[name], parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", name, err)
		}
		g.names = append(g.names, name)
		g.files[name] = f
	}

	return g, nil
}

// ensureType brings the const declaration of the type ID in sync with the
// computed value and reports the generated identifier.
func (g *generator) ensureType(typeName string) (GeneratedTypeId, error) {
	owner, err := g.typeOwner(typeName)
	if err != nil {
		return GeneratedTypeId{}, err
	}

	result := GeneratedTypeId{
		TypeName:  typeName,
		ConstName: typeName + typeIDSuffix,
		ID:        computeID(g.importPath, typeName),
		File:      owner,
	}

	value := strconv.FormatInt(result.ID, 10)
	file, expr, err := g.findValue(result.ConstName)
	if err != nil {
		return GeneratedTypeId{}, err
	}

	if expr != nil {
		start, end := g.offset(expr.Pos()), g.offset(expr.End())
		if g.sources[file][start:end] != value {
			g.edits[file] = append(g.edits[file], edit{start: start, end: end, text: value})
		}

		return result, nil
	}

	decl := g.typeDecl(owner, typeName)
	if decl == nil {
		return GeneratedTypeId{}, fmt.Errorf("type %s: declaration not found in %s", typeName, owner)
	}

	text := "\n" +
		"// " + result.ConstName + " is the generated type ID (FNV-1a 64-bit) of\n" +
		"// " + g.importPath + "." + typeName + ". DO NOT EDIT.\n" +
		"const " + result.ConstName + " int64 = " + value + "\n"

	off := g.offsetAfter(decl.End())
	g.edits[owner] = append(g.edits[owner], edit{start: off, end: off, text: text})

	return result, nil
}

// apply returns the content of every changed file with all edits applied,
// formatted through gofmt.
func (g *generator) apply() (map[string]string, error) {
	changed := make(map[string]string)

	for file, edits := range g.edits {
		content := applyEdits(g.sources[file], edits)
		if content == g.sources[file] {
			continue
		}

		formatted, err := format.Source([]byte(content))
		if err != nil {
			return nil, fmt.Errorf("format %s: %w", file, err)
		}

		changed[file] = string(formatted)
	}

	return changed, nil
}

func (g *generator) typeOwner(typeName string) (string, error) {
	var owners []string
	for _, name := range g.names {
		if hasNamedType(g.files[name], typeName) {
			owners = append(owners, name)
		}
	}

	switch len(owners) {
	case 1:
		return owners[0], nil
	case 0:
		return "", fmt.Errorf("type %s not found in package %s", typeName, g.importPath)
	default:
		return "", fmt.Errorf("type %s declared in more than one file: %s", typeName, strings.Join(owners, ", "))
	}
}

func (g *generator) typeDecl(file, typeName string) *ast.GenDecl {
	for _, decl := range g.files[file].Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}

		for _, s := range gd.Specs {
			if ts, ok := s.(*ast.TypeSpec); ok && ts.Name.Name == typeName && ts.TypeParams == nil {
				return gd
			}
		}
	}

	return nil
}

// findValue locates the const declaration of name and returns the file and
// the expression currently assigned to the name, so the generator can
// maintain the value in place. A nil expression means the name is not
// declared yet. Specs without an explicit per-name value (iota groups,
// implicit repetition, multi-assign) and var declarations cannot be
// maintained and are reported as errors.
func (g *generator) findValue(name string) (string, ast.Expr, error) {
	for _, file := range g.names {
		for _, decl := range g.files[file].Decls {
			gd, ok := decl.(*ast.GenDecl)

			if !ok || (gd.Tok != token.CONST && gd.Tok != token.VAR) {
				continue
			}

			for _, s := range gd.Specs {
				vs, ok := s.(*ast.ValueSpec)
				if !ok {
					continue
				}

				idx := slices.IndexFunc(vs.Names, func(n *ast.Ident) bool { return n.Name == name })
				if idx < 0 {
					continue
				}

				if gd.Tok != token.CONST {
					return "", nil, fmt.Errorf("%s is declared as var, expected a const declaration", name)
				}

				if len(vs.Values) != len(vs.Names) {
					return "", nil, fmt.Errorf("%s is declared in a const spec without an explicit value (iota or implicit repetition), which cannot be maintained", name)
				}

				return file, vs.Values[idx], nil
			}
		}
	}

	return "", nil, nil
}

func (g *generator) offset(p token.Pos) int {
	return g.fset.File(p).Offset(p)
}

func (g *generator) offsetAfter(p token.Pos) int {
	tf := g.fset.File(p)

	if line := tf.Line(p); line+1 <= tf.LineCount() {
		return tf.Offset(tf.LineStart(line + 1))
	}

	return tf.Size()
}

func hasNamedType(f *ast.File, typeName string) bool {
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, s := range gd.Specs {
			if ts, ok := s.(*ast.TypeSpec); ok && ts.Name.Name == typeName && ts.TypeParams == nil {
				return true
			}
		}
	}
	return false
}

func applyEdits(content string, edits []edit) string {
	sorted := slices.Clone(edits)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].start > sorted[j].start })

	for _, e := range sorted {
		content = content[:e.start] + e.text + content[e.end:]
	}

	return content
}

func fnv1a64(s string) uint64 {
	h := fnv.New64a()
	_, _ = io.WriteString(h, s)
	return h.Sum64()
}

func computeID(importPath, typeName string) int64 {
	return int64(fnv1a64(importPath + "." + typeName))
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
