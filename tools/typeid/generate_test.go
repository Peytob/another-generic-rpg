package main

import (
	"strconv"
	"strings"
	"testing"
)

func runGen(t *testing.T, importPath string, sources map[string]string, typeNames ...string) ([]GeneratedTypeId, map[string]string, error) {
	t.Helper()
	order := OrderedTypeId{
		ImportPath: importPath,
		Sources:    sources,
		TypeNames:  typeNames,
	}
	return Generate(order)
}

const healthSource = `package ecscomponent

import (
	"engine/utils/ecs"
)

//go:generate go run typeid

type HealthComponent struct {
}

var HealthComponentType = ecs.ComponentTypeOfT[HealthComponent]()
`

func TestFNV1a64Vectors(t *testing.T) {
	cases := map[string]uint64{
		"":       0xcbf29ce484222325,
		"a":      0xaf63dc4c8601ec8c,
		"foobar": 0x85944171f73967e8,
	}
	for in, want := range cases {
		if got := fnv1a64(in); got != want {
			t.Errorf("fnv1a64(%q) = %#x, want %#x", in, got, want)
		}
	}
}

func TestComputeIDDistinct(t *testing.T) {
	a := computeID("pkg", "A")
	b := computeID("pkg", "B")
	c := computeID("other", "A")
	if a == b || a == c || b == c {
		t.Errorf("expected distinct ids, got %d %d %d", a, b, c)
	}
}

func TestGenerateWithoutOrders(t *testing.T) {
	generated, changed, err := Generate(OrderedTypeId{ImportPath: "p", Sources: map[string]string{"a.go": "package p\n"}})
	if err != nil {
		t.Fatal(err)
	}
	if generated != nil || changed != nil {
		t.Errorf("expected no results, got %v %v", generated, changed)
	}
}

func TestInsertsConstAfterType(t *testing.T) {
	const imp = "game/internal/gamestate/ecscomponent"
	generated, changed, err := runGen(t, imp, map[string]string{"health.go": healthSource}, "HealthComponent")
	if err != nil {
		t.Fatal(err)
	}
	if len(generated) != 1 {
		t.Fatalf("generated = %v, want one entry", generated)
	}
	id := generated[0]
	if id.TypeName != "HealthComponent" || id.ConstName != "HealthComponentTypeID" || id.File != "health.go" {
		t.Errorf("unexpected result %+v", id)
	}
	if id.ID != computeID(imp, "HealthComponent") {
		t.Errorf("id = %d, want %d", id.ID, computeID(imp, "HealthComponent"))
	}
	out, ok := changed["health.go"]
	if !ok {
		t.Fatal("health.go was not changed")
	}
	want := "const HealthComponentTypeID int64 = " + strconv.FormatInt(id.ID, 10)
	if !strings.Contains(out, want) {
		t.Errorf("output missing %q:\n%s", want, out)
	}
	if !strings.Contains(out, "// "+imp+".HealthComponent. DO NOT EDIT.") {
		t.Errorf("output missing marker comment:\n%s", out)
	}
	if !strings.Contains(out, "var HealthComponentType = ecs.ComponentTypeOfT[HealthComponent]()") {
		t.Errorf("output lost original declaration:\n%s", out)
	}
	if i, j := strings.Index(out, "type HealthComponent"), strings.Index(out, want); i < 0 || i > j {
		t.Errorf("const inserted before type declaration:\n%s", out)
	}
}

func TestIdempotent(t *testing.T) {
	_, changed, err := runGen(t, "pkg", map[string]string{"health.go": healthSource}, "HealthComponent")
	if err != nil {
		t.Fatal(err)
	}
	first := changed["health.go"]
	_, changed, err = runGen(t, "pkg", map[string]string{"health.go": first}, "HealthComponent")
	if err != nil {
		t.Fatal(err)
	}
	for name := range changed {
		t.Errorf("second run changed %s:\n%s", name, changed[name])
	}
}

func TestUpdatesExistingValue(t *testing.T) {
	src := `package p

type A struct{ X int }

// ATypeID is the generated type ID (FNV-1a 64-bit) of
// pkg.A. DO NOT EDIT.
const ATypeID int64 = 1
`
	_, changed, err := runGen(t, "pkg", map[string]string{"a.go": src}, "A")
	if err != nil {
		t.Fatal(err)
	}
	out := changed["a.go"]
	want := "const ATypeID int64 = " + strconv.FormatInt(computeID("pkg", "A"), 10)
	if !strings.Contains(out, want) {
		t.Errorf("value not updated:\n%s", out)
	}
	if strings.Contains(out, "int64 = 1\n") {
		t.Errorf("stale value remains:\n%s", out)
	}
	if !strings.Contains(out, "// ATypeID is the generated type ID") {
		t.Errorf("existing comment lost:\n%s", out)
	}
}

func TestMultipleTypesPerFile(t *testing.T) {
	src := `package p

type A struct{}

type B struct{}
`
	_, changed, err := runGen(t, "p", map[string]string{"x.go": src}, "A", "B")
	if err != nil {
		t.Fatal(err)
	}
	out := changed["x.go"]
	aConst := "const ATypeID int64 = " + strconv.FormatInt(computeID("p", "A"), 10)
	bConst := "const BTypeID int64 = " + strconv.FormatInt(computeID("p", "B"), 10)
	ai, bi := strings.Index(out, aConst), strings.Index(out, bConst)
	if ai < 0 || bi < 0 {
		t.Fatalf("consts missing:\n%s", out)
	}
	ta, tb := strings.Index(out, "type A "), strings.Index(out, "type B ")
	if !(ta < ai && ai < tb && tb < bi) {
		t.Errorf("unexpected layout:\n%s", out)
	}
}

func TestUpdatesMultiNameConst(t *testing.T) {
	src := `package p

type A struct{}

type B struct{}

const ATypeID, BTypeID int64 = 1, 2
`
	_, changed, err := runGen(t, "p", map[string]string{"x.go": src}, "A", "B")
	if err != nil {
		t.Fatal(err)
	}
	out := changed["x.go"]
	want := "const ATypeID, BTypeID int64 = " +
		strconv.FormatInt(computeID("p", "A"), 10) + ", " +
		strconv.FormatInt(computeID("p", "B"), 10)
	if !strings.Contains(out, want) {
		t.Errorf("values not updated per name:\n%s", out)
	}
	if strings.Contains(out, "int64 = 1,") || strings.Contains(out, ", 2\n") {
		t.Errorf("stale values remain:\n%s", out)
	}
}

func TestIotaConstRejected(t *testing.T) {
	src := `package p

const (
	XTypeID int64 = iota
	ATypeID
)

type A struct{}
`
	_, _, err := runGen(t, "p", map[string]string{"a.go": src}, "A")
	if err == nil || !strings.Contains(err.Error(), "explicit value") {
		t.Errorf("expected unmaintainable-const error, got %v", err)
	}
}

func TestVarDeclRejected(t *testing.T) {
	src := "package p\n\ntype A struct{}\n\nvar ATypeID int64 = 7\n"
	_, _, err := runGen(t, "p", map[string]string{"a.go": src}, "A")
	if err == nil || !strings.Contains(err.Error(), "declared as var") {
		t.Errorf("expected var error, got %v", err)
	}
}

func TestUnknownType(t *testing.T) {
	_, _, err := runGen(t, "p", map[string]string{"a.go": "package p\n\ntype A struct{}\n"}, "Missing")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not-found error, got %v", err)
	}
}

func TestTypeInMultipleFiles(t *testing.T) {
	sources := map[string]string{
		"a.go": "package p\n\ntype A struct{}\n",
		"b.go": "package p\n\ntype A struct{}\n",
	}
	_, _, err := runGen(t, "p", sources, "A")
	if err == nil || !strings.Contains(err.Error(), "more than one file") {
		t.Errorf("expected multiple-files error, got %v", err)
	}
}

func TestConstInAnotherFileUpdated(t *testing.T) {
	sources := map[string]string{
		"a.go": "package p\n\ntype A struct{}\n",
		"b.go": "package p\n\nconst ATypeID int64 = 7\n",
	}
	_, changed, err := runGen(t, "p", sources, "A")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := changed["a.go"]; ok {
		t.Errorf("a.go must stay unchanged:\n%s", changed["a.go"])
	}
	out, ok := changed["b.go"]
	if !ok {
		t.Fatal("b.go must be updated")
	}
	want := "const ATypeID int64 = " + strconv.FormatInt(computeID("p", "A"), 10)
	if !strings.Contains(out, want) || strings.Contains(out, "= 7") {
		t.Errorf("b.go not updated correctly:\n%s", out)
	}
}

func TestNonStructTypeInExplicitMode(t *testing.T) {
	src := "package p\n\ntype MyInt int\n"
	_, changed, err := runGen(t, "p", map[string]string{"m.go": src}, "MyInt")
	if err != nil {
		t.Fatal(err)
	}
	want := "const MyIntTypeID int64 = " + strconv.FormatInt(computeID("p", "MyInt"), 10)
	if !strings.Contains(changed["m.go"], want) {
		t.Errorf("const missing:\n%s", changed["m.go"])
	}
}
