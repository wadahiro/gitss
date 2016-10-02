package indexer

import (
	"encoding/json"
	// "fmt"
	"log"
	"path"
	// "strings"

	"gopkg.in/olivere/elastic.v3"
)

type ESIndexer struct {
	client *elastic.Client
}

type FileIndex struct {
	Blob     string     `json:"blob"`
	Metadata []Metadata `json:"metadata"`
	Content  string     `json:"content"`
}

type Metadata struct {
	Project string `json:"project"`
	Repo    string `json:"repo"`
	Refs    string `json:"refs"`
	Path    string `json:"path"`
	Ext     string `json:"ext"`
}

type Hit struct {
	Source Source `json:"_source"`
	Highlight map[string][]string `json:"highlight"`
}
type Source struct {
	Blob string `json:"blob"`
	Metadata []Metadata `json:"metadata"`
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
				blob: {
					type: "string",
					index: "not_analyzed"
				},
				metadata: {
					type: "nested",
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
						path: {
							type: "string",
							analyzer: "path_analyzer"
						},
						ext: {
							type: "string",
							index: "not_analyzed"
						}
					}
				}
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

func (esi *ESIndexer) CreateFileIndex(project string, repo string, branch string, filePath string, blob string, content string) error {

	ext := path.Ext(filePath)

	fileIndex := FileIndex{Blob: blob, Metadata: []Metadata{Metadata{Project: project, Repo: repo, Refs: branch, Path: filePath, Ext: ext}}, Content: content}

	_, err := esi.client.Index().
		Index("gosource").
		Type("file").
		Id(blob).
		BodyJson(fileIndex).
		Refresh(true).
		Do()

	if err != nil {
		return err
	}
	return nil
}

func (esi *ESIndexer) UpsertFileIndex(project string, repo string, branch string, filePath string, blob string, content string) error {

	ext := path.Ext(filePath)

	get, err := esi.client.Get().
		Index("gosource").
		Type("file").
		Id(blob).
		Do()

	if err == nil && get.Found {
		var fileIndex FileIndex
		if err := json.Unmarshal(*get.Source, &fileIndex); err != nil {
			return err
		}
		f := func(x Metadata, i int) bool {
			return x.Project == project &&
				x.Repo == repo &&
				x.Refs == branch &&
				x.Path == filePath
		}
		found := find(f, fileIndex.Metadata)
		if found != nil {
			fileIndex.Metadata = append(fileIndex.Metadata, Metadata{Project: project, Repo: repo, Refs: branch, Path: filePath, Ext: ext})
		}

		_, err := esi.client.Update().
			Index("gosource").
			Type("file").
			Id(blob).
			Doc(fileIndex).
			Do()

		if err != nil {
			log.Println("Upsert Doc error", err)
			return err
		}

	} else {
		fileIndex := FileIndex{Blob: blob, Metadata: []Metadata{Metadata{Project: project, Repo: repo, Refs: branch, Path: filePath, Ext: ext}}, Content: content}

		_, err := esi.client.Index().
			Index("gosource").
			Type("file").
			Id(blob).
			BodyJson(fileIndex).
			Refresh(true).
			Do()

		if err != nil {
			log.Println("Add Doc error", err)
			return err
		}
	}

	log.Println("Indexed!")
	return nil
}

func (esi *ESIndexer) SearchQuery(query string) []Hit {
	termQuery := elastic.NewTermsQuery("content", query)
	searchResult, _ := esi.client.Search().
		Index("gosource").                       // search in index "twitter"
		Query(termQuery). // specify the query
		Highlight(elastic.NewHighlight().Field("content")).
		Sort("metadata.path", true). // sort by "user" field, ascending
		From(0).Size(10).            // take documents 0-9
		Pretty(true).                // pretty print request and response JSON
		Do()                         // execute

	list := []Hit{}
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var s Source
			json.Unmarshal(*hit.Source, &s)
			h := &Hit{Source: s, Highlight: hit.Highlight}
			// log.Println("highlight", h)
			
			list = append(list, *h)
		}
	}
	return list
}

func find(f func(s Metadata, i int) bool, s []Metadata) *Metadata {
	for index, x := range s {
		if f(x, index) == true {
			return &x
		}
	}
	return nil
}

func filter(f func(s Metadata, i int) bool, s []Metadata) []Metadata {
	ans := make([]Metadata, 0)
	for index, x := range s {
		if f(x, index) == true {
			ans = append(ans, x)
		}
	}
	return ans
}
