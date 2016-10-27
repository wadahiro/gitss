package indexer

import (
	"bytes"
	"unicode/utf8"

	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
)

type PathHierarchyTokenizer struct {
}

func (t *PathHierarchyTokenizer) Tokenize(input []byte) analysis.TokenStream {
	path := bytes.Split(input, []byte("/"))

	rv := make(analysis.TokenStream, 0, 256)

	for i := range path {
		term := bytes.Join(path[0:i+1], []byte("/"))

		rv = append(rv, &analysis.Token{
			Term:     term,
			Position: i + 1,
			Start:    0,
			End:      len(term),
			Type:     analysis.AlphaNumeric,
		})
	}

	return rv
}

type FullRefTokenizer struct {
}

func (t *FullRefTokenizer) Tokenize(input []byte) analysis.TokenStream {

	rv := make(analysis.TokenStream, 0, 1024)

	offset := 0
	start := 0
	end := 0
	count := 0
	phase := 0

	for currRune, size := utf8.DecodeRune(input[offset:]); currRune != utf8.RuneError; currRune, size = utf8.DecodeRune(input[offset:]) {
		isToken := notFullRefToken(currRune, phase)
		if isToken {
			end = offset + size
		} else {
			if end-start > 0 {
				// build token
				rv = append(rv, &analysis.Token{
					Term:     input[0:end],
					Start:    0,
					End:      end,
					Position: count + 1,
					Type:     analysis.AlphaNumeric,
				})
				count++
			}
			start = offset + size
			end = start
			phase++
		}
		offset += size
	}
	// if we ended in the middle of a token, finish it
	if end-start > 0 {
		// build token
		rv = append(rv, &analysis.Token{
			Term:     input[0:end],
			Start:    0,
			End:      end,
			Position: count + 1,
			Type:     analysis.AlphaNumeric,
		})
	}

	return rv
}

func notFullRefToken(r rune, phase int) bool {
	switch phase {
	case 0:
		return r != []rune(":")[0]
	case 1:
		return r != []rune("/")[0]
	case 2:
		return r != []rune(":")[0]
	}
	return true
}

func PathHierarchyTokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	return &PathHierarchyTokenizer{}, nil
}

func FullRefTokenizerConstructor(config map[string]interface{}, cache *registry.Cache) (analysis.Tokenizer, error) {
	return &FullRefTokenizer{}, nil
}

func init() {
	registry.RegisterTokenizer("path_hierarchy", PathHierarchyTokenizerConstructor)
	registry.RegisterTokenizer("full_ref", FullRefTokenizerConstructor)
}
