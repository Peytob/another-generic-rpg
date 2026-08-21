package main

import (
	"strings"
	"testing"
)

func TestSplitTypes(t *testing.T) {
	got := splitTypes(" A , B,,C ")
	want := []string{"A", "B", "C"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("splitTypes = %v, want %v", got, want)
	}
	if got := splitTypes(" ,, "); len(got) != 0 {
		t.Errorf("splitTypes of blanks = %v, want empty", got)
	}
	if got := splitTypes("A, B, A, A"); strings.Join(got, "|") != "A|B" {
		t.Errorf("splitTypes duplicates = %v, want [A B]", got)
	}
}

func TestStructTypeNames(t *testing.T) {
	src := `package p

type A struct{}

type I interface{ Foo() }

type Alias = A

type G[T any] struct{}

type MyInt int
`
	names, err := structTypeNames(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 1 || names[0] != "A" {
		t.Errorf("structTypeNames = %v, want [A]", names)
	}
}

func TestStructTypeNamesInvalidSource(t *testing.T) {
	if _, err := structTypeNames("not go source"); err == nil {
		t.Error("expected parse error")
	}
}

func TestReadOrderNothingToDo(t *testing.T) {
	t.Skip("ReadOrder requires process-wide flag and environment state")
}
