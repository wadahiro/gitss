package indexer

import (
	"github.com/blevesearch/bleve/analysis"
	"github.com/blevesearch/bleve/registry"
)

func PathHierarchyAnalyzer(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	tokenizer, err := cache.TokenizerNamed("path_hierarchy")
	if err != nil {
		return nil, err
	}
	rv := analysis.Analyzer{
		Tokenizer: tokenizer,
	}
	return &rv, nil
}

func FullRefAnalyzer(config map[string]interface{}, cache *registry.Cache) (*analysis.Analyzer, error) {
	tokenizer, err := cache.TokenizerNamed("full_ref")
	if err != nil {
		return nil, err
	}
	rv := analysis.Analyzer{
		Tokenizer: tokenizer,
	}
	return &rv, nil
}

func init() {
	registry.RegisterAnalyzer("path_hierarchy", PathHierarchyAnalyzer)
	registry.RegisterAnalyzer("full_ref", FullRefAnalyzer)
}
