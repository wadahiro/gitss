package util

import (
	// "fmt"
	"io"
	// "log"
	"fmt"
	"strings"
)

type TextPreview struct {
	Offset   int `json:"offset"`
	previews []string
	Preview  string `json:"preview"`
	Hits     []int  `json:"hits"`
}

func FilterTextPreview(r io.Reader, filter func(line string) bool, before int, after int) []TextPreview {
	scanner := NewLineScanner(r, 1024, before, after)

	previews := []TextPreview{}

	for scanner.HasNext() {
		lineNum, line, ok := scanner.FindLine(filter)

		var beforePreview *TextPreview
		hasPreview := len(previews) > 0

		if hasPreview {
			beforePreview = &previews[len(previews)-1]
		}

		if ok {
			before := scanner.GetBefore()
			offset := lineNum - len(before)

			// Checking need to merge with previous preview

			if hasPreview {
				lastLineNum := beforePreview.Offset + len(beforePreview.previews) - 1
				lastHitNum := beforePreview.Hits[len(beforePreview.Hits)-1]

				// concat
				if (lineNum - lastHitNum) <= scanner.GetAfterSize() {
					appendPreview(beforePreview, lineNum, line)
					continue
				}

				if offset <= (lastLineNum + 1) {
					// reduce duplication
					start := (lastLineNum + 1) - offset
					// fmt.Println("before", len(before), lineNum, line, lastLineNum, offset, start)

					lines := append(before[start:], line)

					appendPreview(beforePreview, lineNum, lines...)
					continue
				}
				// fmt.Println("assert", lineNum, lastHitNum, offset)
			}

			// New preview
			lines := append(before, line)
			preview := TextPreview{Offset: offset, Hits: []int{lineNum}, previews: lines}

			previews = append(previews, preview)

		} else {
			if hasPreview {
				lastHitNum := beforePreview.Hits[len(beforePreview.Hits)-1]
				if (lineNum - lastHitNum) <= scanner.GetAfterSize() {
					beforePreview.previews = append(beforePreview.previews, line)
					// fmt.Println("-----", line, scanner.GetLineNum())
				}
			}
		}
	}

	for i := range previews {
		previews[i].Preview = strings.Join(previews[i].previews, "\n")
		previews[i].previews = nil
	}
	return previews
}

func appendPreview(preview *TextPreview, hitLinuNum int, line ...string) {
	preview.Hits = append(preview.Hits, hitLinuNum)
	preview.previews = append(preview.previews, line...)
}

func Must(e error) {
	if e != nil {
		panic(e)
	}
}
