package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestFilterPreviewText(t *testing.T) {
	r := strings.NewReader(`Title

testtest

add first improvement
add second improvement

add third improvement
add 4th improvement

add bug fix
add bug fix 2

add bug fix 3
fix typo

bug fix 4
fix typo



study1

test
testtest

hoge
`)

	result := FilterTextPreview(r, func(line string) bool {
		return strings.Contains(line, "fix")
	}, 2, 2)

	fmt.Printf("result: %#v", result[0].Hits)

	for i, t := range strings.Split(result[0].Preview, "\n") {
		fmt.Println(i, t)
	}
}
