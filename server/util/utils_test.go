package util

import (
	"fmt"
	"strings"
	"testing"
	"reflect"
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

	for i, preview := range result {
		for j, line := range strings.Split(preview.Preview, "\n") {
			fmt.Println(i, preview.Offset+j, line)
		}
	}

	if len(result) != 1 {
		t.Errorf("result length is %#v", len(result))
	}

	if result[0].Offset != 8 {
		t.Errorf("result offset is %#v", result[0].Offset)
	}

	if !reflect.DeepEqual(result[0].Hits, []int{10, 11, 13, 14, 16, 17}) {
		t.Errorf("result hits is %#v", result[0].Hits)
	}

	if len(strings.Split(result[0].Preview, "\n")) != 12 {
		t.Errorf("result preview lines is %#v", len(strings.Split(result[0].Preview, "\n")))
	}
}


func TestFilterPreviewText2(t *testing.T) {
	r := strings.NewReader(`
1
hit
2
3
4
5
hit
`)

	result := FilterTextPreview(r, func(line string) bool {
		return strings.Contains(line, "hit")
	}, 2, 2)

	for i, preview := range result {
		fmt.Println(preview.Hits)
		for j, line := range strings.Split(preview.Preview, "\n") {
			fmt.Println(i, preview.Offset+j, line)
		}
	}

}
