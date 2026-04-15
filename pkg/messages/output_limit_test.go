package messages

import (
	"reflect"
	"testing"
)

func TestOutputLimit_NoLimit(t *testing.T) {
	input := "line1\nline2"
	result := OutputLimit(input, 5)

	if result != input {
		t.Errorf("expected same string, got: %s", result)
	}
}

func TestOutputLimit_WithLimit(t *testing.T) {
	input := "line1\nline2\nline3\nline4"
	result := OutputLimit(input, 2)

	expected := "line1\nline2\n...[+2 lines]"

	if result != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestOutputLimit_ExactLimit(t *testing.T) {
	input := "line1\nline2"
	result := OutputLimit(input, 2)

	if result != input {
		t.Errorf("expected exact match, got: %s", result)
	}
}

func TestLimitArgs(t *testing.T) {
	args := map[string]string{
		"a": "line1\nline2\nline3",
		"b": "short",
	}

	result := LimitArgs(args, 2)

	expected := map[string]string{
		"a": "line1\nline2\n...[+1 lines]",
		"b": "short",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestLimitArgs_DoesNotMutate(t *testing.T) {
	args := map[string]string{
		"a": "line1\nline2\nline3",
	}

	_ = LimitArgs(args, 2)

	// оригинал должен остаться без изменений
	if args["a"] != "line1\nline2\nline3" {
		t.Errorf("original map was mutated")
	}
}
