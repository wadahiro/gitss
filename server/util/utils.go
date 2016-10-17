package util

import (
	"bufio"
	"container/ring"
	"fmt"
	"io"
	"strings"
)

type TextPreview struct {
	Offset  int    `json:"offset"`
	Preview string `json:"preview"`
	Hits    []int  `json:"hits"`
}

func FilterTextPreview(reader io.Reader, filter func(line string) bool, before int, after int) []*TextPreview {
	scanner := bufio.NewScanner(reader)
	lineNum := -1

	sourceList := []*TextPreview{}

	nextBeforeList := ring.New(before)
	var line string
	hasPrev := false
	eof := false

	for true {
		mode := 0
		offset := -1
		lines := []string{}
		hits := []int{}
		beforeList := ring.New(before)
		afterList := []string{}

		// fmt.Println("start")

		nextBeforeList.Do(func(v interface{}) {
			if v != nil {
				// fmt.Println("next", v)
				beforeList.Value = v.(string)
				beforeList = beforeList.Next()
			}
		})

	L:
		for true {
			if !scanner.Scan() {
				eof = true
				break L
			}
			if hasPrev {
				hasPrev = false
				// fmt.Println("hasPrev", lineNum, line)
			} else {
				lineNum++
				line = scanner.Text()
				fmt.Println(lineNum, line)
			}

			switch mode {
			case 0:
				if filter(line) {
					offset = lineNum
					lines = append(lines, line)
					hits = append(hits, lineNum)
					mode = 1
				} else {
					beforeList.Value = line
					beforeList = beforeList.Next()
				}
			case 1:
				if filter(line) {
					hits = append(hits, lineNum)
					if len(afterList) == 0 {
						lines = append(lines, line)
					} else {
						lines = append(lines, afterList...)
						afterList = nil
					}
				} else {
					if after < 1 {
						hasPrev = true
						break L
					}
					afterList = append(afterList, line)
					nextBeforeList.Value = line
					nextBeforeList = nextBeforeList.Next()

					// fmt.Println("nextBefore", line)

					if len(afterList) == after {
						// fmt.Println("break")
						break L
					}
				}
			}
		}
		if offset > -1 {
			beforeLines := []string{}
			beforeList.Do(func(v interface{}) {
				if v != nil {
					beforeLines = append(beforeLines, v.(string))
				}
			})
			source := &TextPreview{Offset: offset - len(beforeLines), Preview: strings.Join(append(beforeLines, append(lines, afterList...)...), "\n"), Hits: hits}
			sourceList = append(sourceList, source)

			fmt.Println("hitsNum", hits)
		}
		if eof {
			break
		}
	}

	return sourceList
}
