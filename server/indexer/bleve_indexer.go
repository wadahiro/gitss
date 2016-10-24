package indexer

import (
	"encoding/json"
	"log"
	"path"
	"time"
	// "strconv"
	"regexp"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/mapping"
	"github.com/wadahiro/gitss/server/repo"
)

var MAPPING = []byte(`{
	"types": {
		"latest": {
			"enabled": true,
			"dynamic": true,
			"properties": {
				"organization": {
					"enabled": true,
					"dynamic": true,
					"fields": [{
						"type": "text",
						"analyzer": "en",
						"store": true,
						"index": true,
						"include_term_vectors": true,
						"include_in_all": false
					}],
					"default_analyzer": ""
				},
				"project": {
					"enabled": true,
					"dynamic": true,
					"fields": [{
						"type": "text",
						"analyzer": "en",
						"store": true,
						"index": true,
						"include_term_vectors": true,
						"include_in_all": false
					}],
					"default_analyzer": ""
				},
				"repository": {
					"enabled": true,
					"dynamic": true,
					"fields": [{
						"type": "text",
						"analyzer": "en",
						"store": true,
						"index": true,
						"include_term_vectors": true,
						"include_in_all": false
					}],
					"default_analyzer": ""
				},
				"ref": {
					"enabled": true,
					"dynamic": true,
					"fields": [{
						"type": "text",
						"analyzer": "en",
						"store": true,
						"index": true,
						"include_term_vectors": true,
						"include_in_all": false
					}],
					"default_analyzer": ""
				}
			},
			"default_analyzer": ""
		},
		"file": {
			"enabled": true,
			"dynamic": true,
			"properties": {
				"blob": {
					"enabled": true,
					"dynamic": true,
					"fields": [{
						"type": "text",
						"analyzer": "en",
						"store": true,
						"index": true,
						"include_term_vectors": true,
						"include_in_all": false
					}],
					"default_analyzer": ""
				},
				"content": {
					"enabled": true,
					"dynamic": true,
					"fields": [{
						"type": "text",
						"analyzer": "en",
						"store": true,
						"index": true,
						"include_term_vectors": true,
						"include_in_all": true
					}],
					"default_analyzer": ""
				},
				"metadata": {
					"enabled": true,
					"dynamic": true,
					"properties": {
						"path": {
							"enabled": true,
							"dynamic": true,
							"fields": [{
								"type": "text",
								"analyzer": "en",
								"store": true,
								"index": true,
								"include_term_vectors": true,
								"include_in_all": false
							}],
							"default_analyzer": ""
						},
						"ext": {
							"enabled": true,
							"dynamic": true,
							"fields": [{
								"type": "text",
								"analyzer": "en",
								"store": true,
								"index": true,
								"include_term_vectors": true,
								"include_in_all": false
							}],
							"default_analyzer": ""
						},
						"organization": {
							"enabled": true,
							"dynamic": true,
							"fields": [{
								"type": "text",
								"analyzer": "en",
								"store": true,
								"index": true,
								"include_term_vectors": true,
								"include_in_all": false
							}],
							"default_analyzer": ""
						},
						"project": {
							"enabled": true,
							"dynamic": true,
							"fields": [{
								"type": "text",
								"analyzer": "en",
								"store": true,
								"index": true,
								"include_term_vectors": true,
								"include_in_all": false
							}],
							"default_analyzer": ""
						},
						"repository": {
							"enabled": true,
							"dynamic": true,
							"fields": [{
								"type": "text",
								"analyzer": "en",
								"store": true,
								"index": true,
								"include_term_vectors": true,
								"include_in_all": false
							}],
							"default_analyzer": ""
						},
						"ref": {
							"enabled": true,
							"dynamic": true,
							"fields": [{
								"type": "text",
								"analyzer": "en",
								"store": true,
								"index": true,
								"include_term_vectors": true,
								"include_in_all": false
							}],
							"default_analyzer": ""
						}
					},
					"default_analyzer": ""
				}
			},
			"default_analyzer": ""
		}
	},
	"default_mapping": {
		"enabled": true,
		"dynamic": true,
		"default_analyzer": ""
	},
	"type_field": "_type",
	"default_type": "gits",
	"default_analyzer": "standard",
	"default_datetime_parser": "dateTimeOptional",
	"default_field": "_all",
	"store_dynamic": true,
	"index_dynamic": true,
	"analysis": {}
}`)

var BLEVE_HIT_TAG = regexp.MustCompile(`<mark>(.*)</mark>`)

type BleveIndexer struct {
	client bleve.Index
	reader *repo.GitRepoReader
	debug  bool
}

func NewBleveIndexer(reader *repo.GitRepoReader, indexPath string, debugMode bool) Indexer {
	index, err := bleve.Open(indexPath)

	if err == bleve.ErrorIndexPathDoesNotExist {
		var mapping mapping.IndexMappingImpl
		err = json.Unmarshal(MAPPING, &mapping)

		if err != nil {
			log.Println(err)
			panic("error unmarshalling mapping")
		}

		index, err = bleve.New(indexPath, &mapping)

		if err != nil {
			log.Println(err)
			panic("error init index")
		}
	}

	i := &BleveIndexer{client: index, reader: reader, debug: debugMode}

	return i
}

func (b *BleveIndexer) CreateFileIndex(organization string, project string, repo string, branch string, filePath string, blob string, content string) error {

	ext := path.Ext(filePath)

	fileIndex := FileIndex{Blob: blob, Metadata: []Metadata{Metadata{Organization: organization, Project: project, Repository: repo, Ref: branch, Path: filePath, Ext: ext}}, Content: content}

	err := b.index(blob, &fileIndex)

	if err != nil {
		return err
	}
	return nil
}

func (b *BleveIndexer) BatchFileIndex(fileIndex *[]FileIndex) error {
	batch := b.client.NewBatch()
	for i := range *fileIndex {
		f := (*fileIndex)[i]
		batch.Index(f.Blob, f)
	}
	b.client.Batch(batch)
	return nil
}

func (b *BleveIndexer) UpsertFileIndex(organization string, project string, repo string, ref string, filePath string, blob string, content string) error {

	ext := path.Ext(filePath)

	doc, _ := b.client.Document(blob)

	if doc != nil {
		// Update

		// Restore fileIndex from index
		fileIndex := docToFileIndex(doc)

		// Merge metadata
		same := mergeFileIndex(fileIndex, organization, project, repo, ref, filePath, ext)

		if same {
			if b.debug {
				log.Println("Skipped index")
			}
			return nil
		}

		err := b.index(blob, fileIndex)

		if err != nil {
			log.Println("Upsert Doc error", err)
			return err
		}
		if b.debug {
			log.Println("Updated index")
		}

	} else {
		fileIndex := FileIndex{Blob: blob, Metadata: []Metadata{Metadata{Organization: organization, Project: project, Repository: repo, Ref: ref, Path: filePath, Ext: ext}}, Content: content}

		err := b.index(blob, &fileIndex)

		if err != nil {
			log.Println("Add Doc error", err)
			return err
		}
		if b.debug {
			log.Println("Added index")
		}
	}

	return nil
}

func (b *BleveIndexer) index(blob string, f *FileIndex) error {
	var fileDoc interface{}
	j, _ := json.Marshal(f)

	// log.Println(string(j))

	json.Unmarshal(j, &fileDoc)

	// log.Println(f.Content)

	return b.client.Index(blob, fileDoc)
}

func (b *BleveIndexer) SearchQuery(query string) SearchResult {
	start := time.Now()
	result := b.search(query)
	end := time.Now()

	result.Time = (end.Sub(start)).Seconds()
	return result
}

func (b *BleveIndexer) search(query string) SearchResult {

	q := bleve.NewWildcardQuery(query)
	s := bleve.NewSearchRequest(q)

	s.Fields = []string{"blob", "content", "metadata.organization", "metadata.project", "metadata.repository", "metadata.ref", "metadata.path", "metadata.ext"}
	s.Highlight = bleve.NewHighlight()
	searchResults, err := b.client.Search(s)

	if err != nil {
		log.Println(err)
	}

	list := []Hit{}
	hitWordsSet := make(map[string]struct{})

	// log.Println(searchResults)

	for _, hit := range searchResults.Hits {
		doc, err := b.client.Document(hit.ID)
		if err != nil {
			log.Println("Already deleted? blob:" + hit.ID)
			continue
		}

		fileIndex := docToFileIndex(doc)

		s := Source{Blob: hit.ID, Metadata: fileIndex.Metadata}

		// find highlighted words
		hitWordsSet = mergeSet(hitWordsSet, getHitWords(BLEVE_HIT_TAG, hit.Fragments["content"]))

		// log.Println("hitWords", hitWordsSet)?

		// get the file text
		gitRepo := getGitRepo(b.reader, &s)

		// make preview
		preview := gitRepo.FilterBlob(s.Blob, func(line string) bool {
			for k, _ := range hitWordsSet {
				if strings.Contains(line, k) {
					return true
				}
			}
			return false
		}, 3, 3)

		// log.Println(preview)

		h := Hit{Source: s, Preview: preview}
		list = append(list, h)
	}
	// log.Println(searchResults.Total)
	return SearchResult{Hits: list, Size: int64(searchResults.Total)}
}

func docToFileIndex(doc *document.Document) *FileIndex {
	var fileIndex FileIndex
	metadataMap := map[uint64]*Metadata{}

	for i := range doc.Fields {
		f := doc.Fields[i]
		name := strings.Split(f.Name(), ".")
		value := string(f.Value())

		switch name[0] {
		case "blob":
			fileIndex.Blob = value

		case "content":
			fileIndex.Content = value

		case "metadata":
			pos := f.ArrayPositions()[0]
			_, ok := metadataMap[pos]
			if !ok {
				metadataMap[pos] = &Metadata{}
			}
			m := metadataMap[pos]
			switch name[1] {
			case "organization":
				m.Organization = value
			case "project":
				m.Project = value
			case "repository":
				m.Repository = value
			case "ref":
				m.Ref = value
			case "path":
				m.Path = value
			case "ext":
				m.Ext = value
			}
		}
	}

	metadatas := make([]Metadata, len(metadataMap))
	for k, v := range metadataMap {
		metadatas[k] = *v
	}
	// Restored!
	fileIndex.Metadata = metadatas

	return &fileIndex
}
