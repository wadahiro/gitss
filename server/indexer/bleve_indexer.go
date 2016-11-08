package indexer

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	// "strconv"

	"sort"
	"strings"

	"github.com/bcampbell/qs"
	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/document"
	"github.com/blevesearch/bleve/mapping"
	"github.com/blevesearch/bleve/search"
	"github.com/blevesearch/bleve/search/query"
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
						"analyzer": "keyword",
						"store": true,
						"index": true,
						"include_term_vectors": true,
						"include_in_all": false
					}],
					"default_analyzer": ""
				},
				"fullRefs": {
					"enabled": true,
					"dynamic": true,
					"fields": [{
						"type": "text",
						"analyzer": "full_ref",
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
						"analyzer": "keyword",
						"store": true,
						"index": true,
						"include_term_vectors": false,
						"include_in_all": false
					}],
					"default_analyzer": ""
				},
				"project": {
					"enabled": true,
					"dynamic": true,
					"fields": [{
						"type": "text",
						"analyzer": "keyword",
						"store": true,
						"index": true,
						"include_term_vectors": false,
						"include_in_all": false
					}],
					"default_analyzer": ""
				},
				"repository": {
					"enabled": true,
					"dynamic": true,
					"fields": [{
						"type": "text",
						"analyzer": "keyword",
						"store": true,
						"index": true,
						"include_term_vectors": false,
						"include_in_all": false
					}],
					"default_analyzer": ""
				},
				"refs": {
					"enabled": true,
					"dynamic": true,
					"fields": [{
						"type": "text",
						"analyzer": "keyword",
						"store": true,
						"index": true,
						"include_term_vectors": false,
						"include_in_all": false
					}],
					"default_analyzer": ""
				},
				"path": {
					"enabled": true,
					"dynamic": true,
					"fields": [{
						"type": "text",
						"analyzer": "path_hierarchy",
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
						"analyzer": "keyword",
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
						"store": false,
						"index": true,
						"include_term_vectors": true,
						"include_in_all": true
					}],
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

type BleveIndexer struct {
	config config.Config
	client bleve.Index
	reader *repo.GitRepoReader
	debug  bool
}

func NewBleveIndexer(config *config.Config, reader *repo.GitRepoReader) Indexer {
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
	b.client.Close()
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
	fillFileIndex(&requestFileIndex)

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
	fillFileIndex(&requestFileIndex)

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
		fmt.Println(getDocId(f))
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

func (b *BleveIndexer) SearchQuery(query string, filterParams FilterParams, page int) SearchResult {
	start := time.Now()
	result := b.search(query, filterParams, page)
	end := time.Now()

	result.Time = (end.Sub(start)).Seconds()
	return result
}

func (b *BleveIndexer) searchByRefs(organization string, project string, repository string, refs []string, callback func(searchResult *bleve.SearchResult)) error {
	oq := bleve.NewQueryStringQuery("organization:" + organization)
	pq := bleve.NewQueryStringQuery("project:" + project)
	rq := bleve.NewQueryStringQuery("repository:" + repository)
	q1 := bleve.NewConjunctionQuery(oq, pq, rq)

	q2 := bleve.NewDisjunctionQuery()
	for _, ref := range refs {
		rq := bleve.NewQueryStringQuery("refs:" + ref)
		q2.AddQuery(rq)
	}
	s := bleve.NewSearchRequest(bleve.NewConjunctionQuery(q1, q2))
	s.From = 0
	s.Size = 100

	return b.handleSearch(s, callback)
}

func (b *BleveIndexer) searchByOrganization(organization string, callback func(searchResult *bleve.SearchResult)) error {
	q := bleve.NewQueryStringQuery("organization:" + organization)

	s := bleve.NewSearchRequest(q)
	s.From = 0
	s.Size = 100

	return b.handleSearch(s, callback)
}

func (b *BleveIndexer) searchByProject(organization string, project string, callback func(searchResult *bleve.SearchResult)) error {
	oq := bleve.NewQueryStringQuery("organization:" + organization)
	pq := bleve.NewQueryStringQuery("project:" + project)
	q := bleve.NewConjunctionQuery(oq, pq)

	s := bleve.NewSearchRequest(q)
	s.From = 0
	s.Size = 100

	return b.handleSearch(s, callback)
}

func (b *BleveIndexer) searchByRepository(organization string, project string, repository string, callback func(searchResult *bleve.SearchResult)) error {
	oq := bleve.NewQueryStringQuery("organization:" + organization)
	pq := bleve.NewQueryStringQuery("project:" + project)
	rq := bleve.NewQueryStringQuery("repository:" + repository)
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

func (b *BleveIndexer) search(queryString string, filterParams FilterParams, page int) SearchResult {
	p := qs.Parser{DefaultOp: qs.AND}
	q, err := p.Parse(queryString)

	if err != nil {
		log.Printf("Query parse error. %+v", err)
		return SearchResult{
			Query:         queryString,
			FilterParams:  filterParams,
			Hits:          []Hit{},
			Size:          0,
			Current:       0,
			Limit:         10,
			Facets:        nil,
			FullRefsFacet: nil,
		}
	}

	if b.debug {
		log.Printf("ParsedQuery: %v\n", q)
	}

	q = appendFilters(q, filterParams.Exts, "ext", true)
	q = appendFilters(q, filterParams.Organizations, "organization", false)
	q = appendFilters(q, filterParams.Projects, "project", false)
	q = appendFilters(q, filterParams.Repositories, "repository", false)
	q = appendFilters(q, filterParams.Refs, "refs", false)

	s := bleve.NewSearchRequest(q)

	//
	// organizationFacet := bleve.NewFacetRequest("organization", 5)
	// s.AddFacet("organization", organizationFacet)
	fullRefsFacet := bleve.NewFacetRequest("fullRefs", 100)
	extFacet := bleve.NewFacetRequest("ext", 100)
	organizationFacet := bleve.NewFacetRequest("organization", 100)
	projectFacet := bleve.NewFacetRequest("project", 100)
	repositoryFacet := bleve.NewFacetRequest("repository", 100)
	refsFacet := bleve.NewFacetRequest("refs", 100)

	s.AddFacet("fullRefs", fullRefsFacet)
	s.AddFacet("ext", extFacet)
	s.AddFacet("organization", organizationFacet)
	s.AddFacet("project", projectFacet)
	s.AddFacet("repository", repositoryFacet)
	s.AddFacet("refs", refsFacet)

	s.Fields = []string{"blob", "fullRefs", "organization", "project", "repository", "refs", "path", "ext"}
	s.Highlight = bleve.NewHighlight()

	s.From = page * 10
	s.Size = 10 // @TODO

	searchResults, err := b.client.Search(s)

	if err != nil {
		log.Printf("Query error. %+v", err)
		return SearchResult{
			Query:         queryString,
			FilterParams:  filterParams,
			Hits:          []Hit{},
			Size:          0,
			Current:       0,
			Limit:         10,
			Facets:        nil,
			FullRefsFacet: nil,
		}
	}

	list := []Hit{}

	// log.Println(searchResults)
	// // f := searchResults.Facets
	// j, _ := json.MarshalIndent(searchResults, "", "  ")
	// fmt.Printf("facets: %s\n", string(j))

	for _, hit := range searchResults.Hits {
		doc, err := b.client.Document(hit.ID)
		if err != nil {
			log.Println("Already deleted from index? ID:" + hit.ID)
			continue
		}

		fileIndex := docToFileIndex(doc)

		// find highlighted words
		hitWordSet := make(map[string]struct{})
		for hitWord, _ := range hit.Locations["content"] {
			hitWordSet[hitWord] = struct{}{}
		}

		// get the file text
		gitRepo, err := getGitRepo(b.reader, fileIndex)
		if err != nil {
			log.Println("Already deleted from git repository? ID:" + hit.ID)
			continue
		}

		// make preview
		preview := gitRepo.FilterBlob(fileIndex.Blob, func(line string) bool {
			for k, _ := range hitWordSet {
				if strings.Contains(strings.ToLower(line), strings.ToLower(k)) {
					return true
				}
			}
			return false
		}, 3, 3)

		// // wrap hit words with \u0000
		// for i := range preview {
		// 	for k, _ := range hitWordSet {
		// 		preview[i].Preview = strings.Replace(preview[i].Preview, k, "\u0000"+k+"\u0000", -1)
		// 	}
		// }
		keyword := []string{}
		for k, _ := range hitWordSet {
			keyword = append(keyword, k)
		}
		// log.Println(preview)

		h := Hit{Metadata: fileIndex.Metadata, Preview: preview, Keyword: keyword}
		list = append(list, h)
	}

	facets := FacetResults{}

	for k, v := range searchResults.Facets {
		sort.Sort(&v.Terms)

		tf := TermFacets{}
		for _, term := range v.Terms {
			tf = append(tf, TermFacet{Term: term.Term, Count: term.Count})
		}
		facets[k] = FacetResult{
			Field:   v.Field,
			Missing: v.Missing,
			Other:   v.Other,
			Terms:   tf,
			Total:   v.Total,
		}
	}

	// fullRefs
	fullRefsFacetResult := facetResultToFullRefsFacet(searchResults.Facets["fullRefs"])

	// log.Println(searchResults.Total)
	return SearchResult{
		Query:         queryString,
		FilterParams:  filterParams,
		Hits:          list,
		Size:          int64(searchResults.Total),
		Limit:         10,
		Current:       page,
		Facets:        facets,
		FullRefsFacet: fullRefsFacetResult,
	}
}

func appendFilters(q query.Query, list []string, key string, shouldWrap bool) query.Query {
	filters := []query.Query{}
	var wrap string
	if shouldWrap {
		wrap = `"`
	}
	for i := range list {
		val := list[i]
		if val != "" {
			filter := bleve.NewQueryStringQuery(key + ":" + wrap + val + wrap)
			filters = append(filters, filter)
		}
	}
	if len(filters) > 0 {
		return bleve.NewConjunctionQuery(q, bleve.NewDisjunctionQuery(filters...))
	}
	return q
}

func facetResultToFullRefsFacet(facetResult *search.FacetResult) []OrganizationFacet {
	organizationsMap := make(map[string]*OrganizationFacet)
	projectsMap := make(map[string]*ProjectFacet)
	repositoriesMap := make(map[string]*RepositoryFacet)
	refsMap := make(map[string]*RefFacet)

	for i := range facetResult.Terms {
		termFacet := facetResult.Terms[i]

		if ok, organization := isOrganization(termFacet.Term); ok {
			organizationsMap[termFacet.Term] = &OrganizationFacet{Term: organization, Count: termFacet.Count}
		}
		if ok, project := isProject(termFacet.Term); ok {
			projectsMap[termFacet.Term] = &ProjectFacet{Term: project, Count: termFacet.Count}
		}
		if ok, repository := isRepository(termFacet.Term); ok {
			repositoriesMap[termFacet.Term] = &RepositoryFacet{Term: repository, Count: termFacet.Count}
		}
		if ok, ref := isRef(termFacet.Term); ok {
			refsMap[termFacet.Term] = &RefFacet{Term: ref, Count: termFacet.Count}
		}
	}

	for k, ref := range refsMap {
		parent := repositoriesMap[k[0:strings.LastIndex(k, ":")]]
		parent.Refs = append(parent.Refs, *ref)
	}

	for k, repository := range repositoriesMap {
		parent := projectsMap[strings.Split(k, "/")[0]]
		parent.Repositories = append(parent.Repositories, *repository)
	}

	for k, project := range projectsMap {
		parent := organizationsMap[strings.Split(k, ":")[0]]
		parent.Projects = append(parent.Projects, *project)
	}

	organizations := []OrganizationFacet{}

	for _, organization := range organizationsMap {
		organizations = append(organizations, *organization)
	}

	return organizations
}

func isOrganization(path string) (bool, string) {
	if !strings.Contains(path, ":") {
		return true, path
	} else {
		return false, ""
	}
}

func isProject(path string) (bool, string) {
	if strings.Contains(path, ":") && !strings.Contains(path, "/") {
		return true, strings.Split(path, ":")[1]
	} else {
		return false, ""
	}
}

func isRepository(path string) (bool, string) {
	if strings.Count(path, ":") == 1 && strings.Contains(path, "/") {
		return true, strings.Split(path, "/")[1]
	} else {
		return false, ""
	}
}

func isRef(path string) (bool, string) {
	if strings.Count(path, ":") == 2 && strings.Contains(path, "/") {
		return true, path[strings.LastIndex(path, ":")+1:]
	} else {
		return false, ""
	}
}

func docToFileIndex(doc *document.Document) *FileIndex {
	var fileIndex FileIndex
	fullRefsMap := map[uint64]string{}
	refsMap := map[uint64]string{}

	for i := range doc.Fields {
		f := doc.Fields[i]
		name := f.Name()
		value := string(f.Value())

		switch name {
		case "blob":
			fileIndex.Blob = value

		case "fullRefs":
			pos := f.ArrayPositions()[0]
			_, ok := fullRefsMap[pos]
			if !ok {
				fullRefsMap[pos] = value
			}

		case "content":
			fileIndex.Content = value

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

	fullRefs := make([]string, len(fullRefsMap))
	for k, v := range fullRefsMap {
		fullRefs[k] = v
	}
	// Restored!
	fileIndex.FullRefs = fullRefs

	refs := make([]string, len(refsMap))
	for k, v := range refsMap {
		refs[k] = v
	}
	// Restored!
	fileIndex.Metadata.Refs = refs

	return &fileIndex
}
