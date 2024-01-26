package node

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSetValue(t *testing.T) {
	type testCase struct {
		target    interface{}
		after     interface{}
		fieldName string
		value     interface{}
	}

	type book struct {
		Title string
	}

	a := [2]int{0, 1}
	m := map[string]int{"a": 2}
	cases := []testCase{
		{new(book), &book{Title: "Dragon Book"}, "Title", "Dragon Book"},
		{&book{Title: "a"}, new(book), "Title", ""},
		{new([2]int), &a, "[1]", 1},
		{&map[string]int{}, &m, `["a"]`, 2},
	}
	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			err := setValue(c.target, c.fieldName, c.value)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(c.after, c.target) {
				t.Log(cmp.Diff(c.after, c.target))
				t.Fatal("not equal")
			}
		})
	}
}

func TestGetFirst(t *testing.T) {
	type testCase struct {
		fieldName string
		end       int
		first     interface{}
		t         tokenType
	}
	cases := []testCase{
		{"FieldName", 9, "FieldName", tokenTypeIdent},
		{"FieldName.", 9, "FieldName", tokenTypeIdent},
		{"FieldName[0]", 9, "FieldName", tokenTypeIdent},
		{"[0]", 3, 0, tokenTypeInt | tokenTypeKey},
		{`["key"]`, 7, "key", tokenTypeString | tokenTypeKey},
	}
	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			first, end, tt, err := getFirst(c.fieldName)
			if err != nil {
				t.Fatal(err)
			}
			if first != c.first || tt != c.t || end != c.end {
				t.Fatalf("expected (%T %v, %d, end %d) but got (%T %v, %d, end %d)", c.first, c.first, c.t, c.end, first, first, tt, end)
			}
		})
	}
}
