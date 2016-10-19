package indexer

import (
	// "bytes"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
	// "fmt"
	"log"
	"path"
	// "strings"
	"github.com/wadahiro/gits/server/repo"

	// "fmt"

	"gopkg.in/olivere/elastic.v3"
)

type ESIndexer struct {
	client *elastic.Client
	reader *repo.GitRepoReader
	debug bool
}

// var LINE_TAG = regexp.MustCompile(`^\[([0-9]+)\]\s(.*)`)

func NewESIndexer(reader *repo.GitRepoReader, debugMode bool) Indexer {
	client, err := elastic.NewClient(elastic.SetURL())
	if err != nil {
		panic(err)
	}
	i := &ESIndexer{client: client, reader: reader, debug: debugMode}
	i.Init()
	return i
}

const PRE_TAG = "\u0001"
const POST_TAG = "\u0001"

var ES_HIT_TAG = regexp.MustCompile(`\x{0001}(.*)\x{0001}`)

var CRLF_PATTERN = regexp.MustCompile(`\r?\n|\r`)

func (esi *ESIndexer) Init() {

	// esi.client.DeleteIndex("gosource").Do()
	_, err := esi.client.CreateIndex("gosource").BodyString(`
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
			char_filter: {
				remove_tags: {
					type: "pattern_replace",
					pattern: "^\\[[0-9]+\\]\\u0020",
					flags: "MULTILINE",
					replacement: ""
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
					char_filter: ["remove_tags"],
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
						organization: {
							type: "multi_field",
							fields: {
								organization: {
									type: "string",
									index: "analyzed"
								},
								full: {
									type: "string",
									index: "not_analyzed"
								}
							}
						},
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
				},
				content: {
					type: "string",
					index_options: "offsets",
					analyzer: "kuromoji_analyzer"
				}
			}
		}
	}
}
		`).Do()

	if err != nil {
		log.Println(err)
	}
}

func (esi *ESIndexer) CreateFileIndex(organization string, project string, repo string, branch string, filePath string, blob string, content string) error {

	ext := path.Ext(filePath)

	fileIndex := FileIndex{Blob: blob, Metadata: []Metadata{Metadata{Organization: organization, Project: project, Repository: repo, Refs: branch, Path: filePath, Ext: ext}}, Content: content}

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

func (esi *ESIndexer) UpsertFileIndex(organization string, project string, repo string, branch string, filePath string, blob string, content string) error {

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

		mergeFileIndex(&fileIndex, organization, project, repo, branch, filePath, ext)

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
		lines := strings.Split(content, "\n")
		newLines := []string{}
		for i, l := range lines {
			newLines = append(newLines, "["+strconv.Itoa(i+1)+"] "+l)
		}

		fileIndex := FileIndex{Blob: blob, Metadata: []Metadata{Metadata{Organization: organization, Project: project, Repository: repo, Refs: branch, Path: filePath, Ext: ext}}, Content: strings.Join(newLines, "\n")}

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

func (esi *ESIndexer) SearchQuery(query string) SearchResult {
	start := time.Now()
	result := esi.search(query)
	end := time.Now()

	result.Time = (end.Sub(start)).Seconds()
	return result
}

func (esi *ESIndexer) search(query string) SearchResult {
	// termQuery := elastic.NewTermsQuery("content", strings.Split(query, " "))
	q := elastic.NewQueryStringQuery(query).DefaultField("content").DefaultOperator("AND")
	searchResult, err := esi.client.Search().
		Index("gosource"). // search in index "twitter"
		FetchSourceContext(elastic.NewFetchSourceContext(true).Include("blob", "metadata")).
		Query(q). // specify the query
		Highlight(elastic.NewHighlight().Field("content").PreTags(PRE_TAG).PostTags(POST_TAG)).
		Sort("metadata.path", true). // sort by "user" field, ascending
		From(0).Size(10).            // take documents 0-9
		Pretty(true).                // pretty print request and response JSON
		Do()                         // execute

	if err != nil {
		log.Println("error", err)
		return SearchResult{}
	}

	list := []Hit{}
	hitWordsSet := make(map[string]struct{})

	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var s Source
			json.Unmarshal(*hit.Source, &s)

			// find highlighted words
			hitWordsSet = mergeSet(hitWordsSet, getHitWords(ES_HIT_TAG, hit.Highlight["content"]))

			log.Println("hitWords", hitWordsSet)

			// get the file text
			gitRepo := getGitRepo(esi.reader, &s)

			// make preview
			preview := gitRepo.FilterBlob(s.Blob, func(line string) bool {
				for k, _ := range hitWordsSet {
					if strings.Contains(line, k) {
						log.Println(k)
						return true
					}
				}
				return false
			}, 3, 3)

			log.Println(preview)

			// hsList := []HighlightSource{}

			// for _, hc := range hit.Highlight["content"] {
			// 	groups := HIT_TAG.FindAllStringSubmatch(l, -1)
			// 	for _, group := range groups {
			// 		hitWordsSet[group[1]] = struct{}{}
			// 	}
			// 	// hs := HighlightSource{Offset: first, Content: strings.Join(list, "\n")}
			// 	// hsList = append(hsList, hs)
			// }

			h := Hit{Source: s, Preview: preview}
			list = append(list, h)
		}
	}

	return SearchResult{Hits: list, Size: searchResult.Hits.TotalHits}
}

