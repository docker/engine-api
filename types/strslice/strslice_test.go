package strslice

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStrSliceMarshalJSON(t *testing.T) {
	for _, testcase := range []struct {
		input    StrSlice
		expected string
	}{
		{input: nil, expected: ""},

		// MADNESS(stevvooe): No clue why nil would be "" but empty would be
		// "null". Maintaining compatibility. This is pretty bad.
		{StrSlice{}, "null"},
		{StrSlice{"/bin/sh", "-c", "echo"}, `["/bin/sh","-c","echo"]`},
	} {
		data, err := testcase.input.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != testcase.expected {
			t.Fatalf("%#v: expected %v, got %v", testcase.input, testcase.expected, string(data))
		}
	}
}

func TestStrSliceUnmarshalJSON(t *testing.T) {
	parts := map[string][]string{
		"":   {"default", "values"},
		"[]": {},
		`["/bin/sh","-c","echo"]`: {"/bin/sh", "-c", "echo"},
	}
	for json, expectedParts := range parts {
		strs := StrSlice{"default", "values"}
		if err := strs.UnmarshalJSON([]byte(json)); err != nil {
			t.Fatal(err)
		}

		actualParts := []string(strs)
		if !reflect.DeepEqual(actualParts, expectedParts) {
			t.Fatalf("%#v: expected %v, got %v", json, expectedParts, actualParts)
		}

	}
}

func TestStrSliceUnmarshalString(t *testing.T) {
	var e StrSlice
	echo, err := json.Marshal("echo")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(echo, &e); err != nil {
		t.Fatal(err)
	}

	if len(e) != 1 {
		t.Fatalf("expected 1 element after unmarshal: %q", e)
	}

	if e[0] != "echo" {
		t.Fatalf("expected `echo`, got: %q", e[0])
	}
}

func TestStrSliceUnmarshalSlice(t *testing.T) {
	var e StrSlice
	echo, err := json.Marshal([]string{"echo"})
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(echo, &e); err != nil {
		t.Fatal(err)
	}

	if len(e) != 1 {
		t.Fatalf("expected 1 element after unmarshal: %q", e)
	}

	if e[0] != "echo" {
		t.Fatalf("expected `echo`, got: %q", e[0])
	}
}

func TestStrSliceToString(t *testing.T) {
	for _, testcase := range []struct {
		input    StrSlice
		expected string
	}{
		{New(""), ""},
		{New("one"), "one"},
		{New("one", "two"), "one two"},
	} {
		toString := testcase.input.String()
		if toString != testcase.expected {
			t.Fatalf("Expected %v, got %v", testcase.expected, toString)
		}
	}
}
