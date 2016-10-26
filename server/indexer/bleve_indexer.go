package indexer

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	// "strconv"
	"regexp"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/mapping"
	"github.com/wadahiro/gitss/server/config"
	"github.com/wadahiro/gitss/server/repo"
)

var MAPPING = []byte(`{
	"types": {
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
						"refs": {
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
	"default_type": "file",
	"default_analyzer": "standard",
	"default_datetime_parser": "dateTimeOptional",
	"default_field": "_all",
	"store_dynamic": true,
	"index_dynamic": true,
	"analysis": {}
}`)

var BLEVE_HIT_TAG = regexp.MustCompile(`<mark>(.*?)</mark>`)

type BleveIndexer struct {
	config config.Config
	client bleve.Index
	reader *repo.GitRepoReader
	debug  bool
}

func NewBleveIndexer(config config.Config, reader *repo.GitRepoReader) Indexer {
	indexPath := config.DataDir + "/bleve_index"
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

	i := &BleveIndexer{client: index, reader: reader, debug: config.Debug}

	return i
}

func (b *BleveIndexer) CreateFileIndex(requestFileIndex FileIndex) error {
	return b.create(requestFileIndex, nil)
}

func (b *BleveIndexer) UpsertFileIndex(requestFileIndex FileIndex) error {
	return b.upsert(requestFileIndex, nil)
}

func (b *BleveIndexer) BatchFileIndex(requestBatch []FileIndexOperation) error {
	batch := b.client.NewBatch()
	for i := range requestBatch {
		op := requestBatch[i]
		f := op.FileIndex

		switch op.Method {
		case ADD:
			b.upsert(f, batch)
		case DELETE:
			b.delete(f, batch)
			batch.Delete(f.Blob)
		}
	}
	b.client.Batch(batch)
	return nil
}

func (b *BleveIndexer) DeleteIndexByRefs(organization string, project string, repository string, refs []string) error {
	b.searchByRefs(organization, project, repository, refs, func(searchResult *bleve.SearchResult) {
		batch := b.client.NewBatch()

		for i := range searchResult.Hits {
			hit := searchResult.Hits[i]
			doc, err := b.client.Document(hit.ID)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = b.deleteByDoc(doc, refs, batch)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		err := b.client.Batch(batch)
		if err != nil {
			fmt.Println(err)
		}
	})

	return nil
}

func (b *BleveIndexer) create(requestFileIndex FileIndex, batch *bleve.Batch) error {
	fillFileExt(&requestFileIndex)

	err := b._index(&requestFileIndex, batch)

	if err != nil {
		log.Println("Create Doc error", err)
		return err
	}
	if b.debug {
		log.Println("Created index")
	}
	return nil
}

func (b *BleveIndexer) upsert(requestFileIndex FileIndex, batch *bleve.Batch) error {
	fillFileExt(&requestFileIndex)

	doc, _ := b.client.Document(getDocId(&requestFileIndex))

	// Create case
	if doc == nil {
		return b.create(requestFileIndex, batch)
	}

	// Update case

	// Restore fileIndex from index
	fileIndex := docToFileIndex(doc)

	// Merge ref
	same := mergeRef(fileIndex, requestFileIndex.Metadata.Refs)

	if same {
		if b.debug {
			log.Println("Skipped index")
		}
		return nil
	}

	err := b._index(fileIndex, batch)

	if err != nil {
		log.Println("Update Doc error", err)
		return err
	}
	if b.debug {
		log.Println("Updated index")
	}

	return nil
}

func (b *BleveIndexer) delete(requestFileIndex FileIndex, batch *bleve.Batch) error {
	doc, err := b.client.Document(getDocId(&requestFileIndex))
	if err != nil {
		return err
	}
	return b.deleteByDoc(doc, requestFileIndex.Metadata.Refs, batch)
}

func (b *BleveIndexer) deleteByDoc(doc *document.Document, refs []string, batch *bleve.Batch) error {
	if doc != nil {
		// Restore fileIndex from index
		fileIndex := docToFileIndex(doc)

		// Remove ref
		allRemoved := removeRef(fileIndex, refs)

		if allRemoved {
			err := b._delete(doc.ID, batch)

			if err != nil {
				log.Println("Delete Doc error", err)
				return err
			}
			if b.debug {
				log.Println("Deleted index")
			}
		} else {
			err := b._index(fileIndex, batch)

			if err != nil {
				log.Println("Update(for delete) Doc error", err)
				return err
			}
			if b.debug {
				log.Println("Updated(for delete) index")
			}
		}
	}
	return nil
}

func (b *BleveIndexer) _index(f *FileIndex, batch *bleve.Batch) error {
	if batch == nil {
		return b.client.Index(getDocId(f), f)
	} else {
		return batch.Index(getDocId(f), f)
	}
}

func (b *BleveIndexer) _delete(docID string, batch *bleve.Batch) error {
	if batch == nil {
		return b.client.Delete(docID)
	}
	batch.Delete(docID)
	return nil
}

func (b *BleveIndexer) SearchQuery(query string) SearchResult {
	start := time.Now()
	result := b.search(query)
	end := time.Now()

	result.Time = (end.Sub(start)).Seconds()
	return result
}

func (b *BleveIndexer) searchByRefs(organization string, project string, repository string, refs []string, callback func(searchResult *bleve.SearchResult)) error {
	oq := bleve.NewQueryStringQuery("metadata.organization:" + organization)
	pq := bleve.NewQueryStringQuery("metadata.project:" + project)
	rq := bleve.NewQueryStringQuery("metadata.repository:" + repository)
	q1 := bleve.NewConjunctionQuery(oq, pq, rq)

	q2 := bleve.NewDisjunctionQuery()
	for _, ref := range refs {
		rq := bleve.NewQueryStringQuery("metadata.refs:" + ref)
		q2.AddQuery(rq)
	}
	s := bleve.NewSearchRequest(bleve.NewConjunctionQuery(q1, q2))
	s.From = 0
	s.Size = 100

	return b.handleSearch(s, callback)
}

func (b *BleveIndexer) searchByOrganization(organization string, callback func(searchResult *bleve.SearchResult)) error {
	q := bleve.NewQueryStringQuery("metadata.organization:" + organization)

	s := bleve.NewSearchRequest(q)
	s.From = 0
	s.Size = 100

	return b.handleSearch(s, callback)
}

func (b *BleveIndexer) searchByProject(organization string, project string, callback func(searchResult *bleve.SearchResult)) error {
	oq := bleve.NewQueryStringQuery("metadata.organization:" + organization)
	pq := bleve.NewQueryStringQuery("metadata.project:" + project)
	q := bleve.NewConjunctionQuery(oq, pq)

	s := bleve.NewSearchRequest(q)
	s.From = 0
	s.Size = 100

	return b.handleSearch(s, callback)
}

func (b *BleveIndexer) searchByRepository(organization string, project string, repository string, callback func(searchResult *bleve.SearchResult)) error {
	oq := bleve.NewQueryStringQuery("metadata.organization:" + organization)
	pq := bleve.NewQueryStringQuery("metadata.project:" + project)
	rq := bleve.NewQueryStringQuery("metadata.repository:" + repository)
	q := bleve.NewConjunctionQuery(oq, pq, rq)

	s := bleve.NewSearchRequest(q)
	s.From = 0
	s.Size = 100

	return b.handleSearch(s, callback)
}

func (b *BleveIndexer) handleSearch(searchRequest *bleve.SearchRequest, callback func(searchResult *bleve.SearchResult)) error {
	for {
		searchResult, err := b.client.Search(searchRequest)
		if err != nil {
			return err
		}

		if len(searchResult.Hits) == 0 {
			break
		}

		callback(searchResult)

		searchRequest.From = searchRequest.From + len(searchResult.Hits)
	}
	return nil
}

func (b *BleveIndexer) search(query string) SearchResult {

	q := bleve.NewWildcardQuery(query)
	s := bleve.NewSearchRequest(q)

	s.Fields = []string{"blob", "content", "metadata.organization", "metadata.project", "metadata.repository", "metadata.refs", "metadata.path", "metadata.ext"}
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
			log.Println("Already deleted from index? ID:" + hit.ID)
			continue
		}

		fileIndex := docToFileIndex(doc)

		s := Source{Blob: fileIndex.Blob, Metadata: fileIndex.Metadata}

		// find highlighted words
		hitWordsSet = mergeSet(hitWordsSet, getHitWords(BLEVE_HIT_TAG, hit.Fragments["content"]))

		// log.Println("hitWords", hitWordsSet)?

		// get the file text
		gitRepo, err := getGitRepo(b.reader, &s)
		if err != nil {
			log.Println("Already deleted from git repository? ID:" + hit.ID)
			continue
		}

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
	refsMap := map[uint64]string{}

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
			switch name[1] {
			case "organization":
				fileIndex.Metadata.Organization = value
			case "project":
				fileIndex.Metadata.Project = value
			case "repository":
				fileIndex.Metadata.Repository = value
			case "refs":
				pos := f.ArrayPositions()[0]
				_, ok := refsMap[pos]
				if !ok {
					refsMap[pos] = value
				}
			case "path":
				fileIndex.Metadata.Path = value
			case "ext":
				fileIndex.Metadata.Ext = value
			}
		}
	}

	refs := make([]string, len(refsMap))
	for k, v := range refsMap {
		refs[k] = v
	}
	// Restored!
	fileIndex.Metadata.Refs = refs

	return &fileIndex
}
