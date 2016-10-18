package util

import (
	"bufio"
	"container/ring"
	"fmt"
	"io"
)

type LineScanner struct {
	reader      *bufio.Reader
	beforeLines *ring.Ring
	before      int
	after       int
	debug       bool
	lineNum     int
}

func NewLineScanner(reader io.Reader, bufferSize int, before int, after int) *LineScanner {
	r := bufio.NewReaderSize(reader, bufferSize)
	beforeLines := ring.New(before + 1)
	return &LineScanner{reader: r, debug: false, beforeLines: beforeLines, lineNum: -1, before: before, after: after}
}

func (l *LineScanner) Peek() (b []byte, err error) {
	return l.reader.Peek(1)
}

func (l *LineScanner) GetBeforeSize() int {
	return l.before
}

func (l *LineScanner) GetAfterSize() int {
	return l.after
}

func (l *LineScanner) GetBefore() []string {
	list := []string{}
	l.beforeLines.Do(func(v interface{}) {
		// fmt.Println("do ", v)
		if v != nil {
			list = append(list, v.(string))
		}
	})
	return list[0 : len(list)-1]
}

func (l *LineScanner) FindLine(filter func(line string) bool) (int, string, bool) {
	return l.searchLine(filter, nil)
}

func (l *LineScanner) searchLine(filter func(line string) bool, chunked *string) (int, string, bool) {
	b, isPrefix, err := l.reader.ReadLine()

	if err == io.EOF {
		return -1, "", false
	}

	line := string(b)

	// first time parsing the line
	if !isChunked(chunked) {
		l.lineNum++
		if l.debug {
			fmt.Println(l.lineNum, line)
		}

		beforeLine := line
		if isPrefix {
			beforeLine += "..."
		}
		l.beforeLines.Value = beforeLine
		l.beforeLines = l.beforeLines.Next()
	}

	// filtered case
	if filter(line) {
		if isChunked(chunked) {
			line = "... " + line
		}
		if isPrefix {
			line = line + " ..."
		}

		return l.lineNum, line, true
	}

	// not filtered case
	if !isPrefix {
		if isChunked(chunked) {
			line = *chunked + " ..."
		}
		return l.lineNum, line, false
	}

	// search chunked line
	return l.searchLine(filter, &line)
}

func isChunked(chunked *string) bool {
	return chunked != nil
}

// func (l *LineScanner) Next() string {
// 	b, isPrefix, _ := l.reader.ReadLine()
// 	line := string(b)
// 	if isPrefix {
// 		line += " ..."

// 		for true {
// 			_, isPrefix, err := l.reader.ReadLine()
// 			if !isPrefix || err != nil {
// 				break
// 			}
// 		}
// 	}
// 	return line
// }

func (l *LineScanner) HasNext() bool {
	_, err := l.reader.Peek(1)
	if l.debug && err != nil {
		fmt.Println(err)
	}
	return err == nil
}
