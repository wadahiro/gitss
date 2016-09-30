package indexer

import (
	"strings"

	"gopkg.in/olivere/elastic.v3"
)

type ESIndexer struct {
	client *elastic.Client
}

type FileIndex struct {
	Project    string `json:"project"`
	Repository string `json:"repository"`
	Refs       string `json:"refs"`
	Blob       string `json:"blob"`
	Path       string `json:"path"`
	Ext        string `json:"ext"`
	Content    string `json:"content"`
}

func NewESIndexer() Indexer {
	client, err := elastic.NewClient(elastic.SetURL())
	if err != nil {
		panic(err)
	}
	i := &ESIndexer{client: client}
	i.Init()
	return i
}

func (esi *ESIndexer) Init() {
	esi.client.CreateIndex("gosource").BodyString(`
{
	settings: {
		analysis: {
			filter: {
				pos_filter: {
					type: "kuromoji_part_of_speech",
					stoptags: ["助詞-格助詞-一般", "助詞-終助詞"]
				},
				greek_lowercase_filter: {
					type: "lowercase",
					language: "greek"
				}
			},
			analyzer: {
				path_analyzer: {
					type: "custom",
					tokenizer: "path_tokenizer"
				},
				kuromoji_analyzer: {
					type: "custom",
					tokenizer: "kuromoji_tokenizer",
					filter: ["kuromoji_baseform", "pos_filter", "greek_lowercase_filter", "cjk_width"]
				}
			},
			tokenizer: {
				path_tokenizer: {
					type: "path_hierarchy",
					reverse: true
				}
			}
		}
	},
	mappings: {
		file: {
			properties: {
				project: {
					type: "multi_field",
					fields: {
						project: {
							type: "string",
							index: "analyzed"
						},
						full: {
							type: "string",
							index: "not_analyzed"
						}
					}
				},
				repository: {
					type: "multi_field",
					fields: {
						repository: {
							type: "string",
							index: "analyzed"
						},
						full: {
							type: "string",
							index: "not_analyzed"
						}
					}
				},
				refs: {
					type: "multi_field",
					fields: {
						refs: {
							type: "string",
							index: "analyzed"
						},
						full: {
							type: "string",
							index: "not_analyzed"
						}
					}
				},
				blob: {
					type: "string",
					index: "not_analyzed"
				},
				path: {
					type: "string",
					analyzer: "path_analyzer"
				},
				ext: {
					type: "string",
					index: "not_analyzed"
				},
				contents: {
					type: "string",
					index_options: "offsets",
					analyzer: "kuromoji_analyzer"
				}
			}
		}
	}
}
		`).Do()
}

func (esi *ESIndexer) CreateFileIndex(project string, repo string, branch string, fileName string, blob string, content string) {

	exts := strings.Split(fileName, ".")
	ext := ""
	if len(exts) > 1 {
		ext = exts[len(exts)-1]
	}

	fileIndex := FileIndex{Project: project, Repository: repo, Refs: branch, Path: fileName, Ext: ext, Blob: blob, Content: content}

	_, err := esi.client.Index().
		Index("gosource").
		Type("file").
		Id(blob).
		BodyJson(fileIndex).
		Refresh(true).
		Do()

	if err != nil {
		panic(err)
	}
}
