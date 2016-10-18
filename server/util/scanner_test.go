package util

import (
	"strings"
	"testing"
)

func TestLineScanner(t *testing.T) {
	r := strings.NewReader(`1
HIT
2
3
4
5
HIT
6
7
8
9
10
HIT
11
`)
	scanner := NewLineScanner(r, 1024, 2, 2)
	type Result struct {
		line    string
		ok      bool
		lineNum int
	}
	result := []Result{}

	for scanner.HasNext() {
		lineNum, line, ok := scanner.FindLine(func(line string) bool {
			return strings.Contains(line, "HIT")
		})
		result = append(result, Result{line: line, ok: ok, lineNum: lineNum})
	}

	actual := len(result)
	expected := 14
	if actual != expected {
		t.Errorf("got %v, want %v", actual, expected)
	}

	for i := range result {
		if i == 1 || i == 6 || i == 12 {
			if !result[i].ok {
				t.Errorf("got false, want true")
			}
		} else {
			if result[i].ok {
				t.Errorf("got true\nwant false")
			}
		}
	}
}
