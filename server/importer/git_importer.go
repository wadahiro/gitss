package importer

import (
	"fmt"
	"log"
	"strings"

	pb "gopkg.in/cheggaaa/pb.v1"
	// "io"
	// "io/ioutil"
	// "os"
	"sync"

	"time"

	"github.com/pkg/errors"
	"github.com/wadahiro/gitss/server/config"
	"github.com/wadahiro/gitss/server/indexer"
	"github.com/wadahiro/gitss/server/repo"
	// "github.com/wadahiro/gitss/server/util"
	"bytes"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

type GitImporter struct {
	config  *config.Config
	indexer indexer.Indexer
	reader  *repo.GitRepoReader
	debug   bool
}

func NewGitImporter(config *config.Config, indexer indexer.Indexer) *GitImporter {
	r := repo.NewGitRepoReader(config)
	return &GitImporter{config: config, indexer: indexer, reader: r, debug: config.Debug}
}

func (g *GitImporter) Run(organization string, project string, url string) {
	log.Printf("Clone from %s %s %s\n", organization, project, url)

	repo, err := g.reader.CloneGitRepo(organization, project, url)
	if err != nil {
		log.Printf("Not found the repository: %s:%s/%s %+v\n", organization, project, url, err)
		return
	}

	repo.FetchAll()

	log.Printf("Fetched all. %s %s %s \n", organization, project, url)

	// branches and tags in the git repository (include/exclude filters are applied)
	branchMap, tagMap, err := repo.GetLatestCommitIdsMap()
	if err != nil {
		log.Printf("Failed to get latest commitIds of branch and tag. %+v\n", err)
		return
	}

	// branches in the config file
	indexed := g.config.GetIndexed(organization, project, repo.Repository)

	log.Printf("Start indexing for %s:%s/%s branches: %v -> %v, tags: %v -> %v\n", organization, project, repo.Repository, indexed.Branches, branchMap, indexed.Tags, tagMap)

	// get sizeLimit for this repository
	sizeLimit := g.config.GetSizeLimit(organization, project, repo.Repository)

	// progress bar
	bar := pb.StartNew(0)
	bar.ShowPercent = true

	start := time.Now()

	err = g.runIndexing(bar, repo, url, indexed, branchMap, tagMap, sizeLimit)
	if err != nil {
		log.Printf("Failed to index. %+v", err)
		return
	}

	// Remove index for removed branches
	removeBranches := []string{}
	for ref, _ := range indexed.Branches {
		found := false
		for branch := range branchMap {
			if ref == branch {
				found = true
				break
			}
		}
		if !found {
			removeBranches = append(removeBranches, ref)
		}
	}

	bar.Total = bar.Total + int64(len(removeBranches))

	// Remove index for removed tags
	removeTags := []string{}
	for ref, _ := range indexed.Tags {
		found := false
		for branch := range branchMap {
			if ref == branch {
				found = true
				break
			}
		}
		if !found {
			removeTags = append(removeTags, ref)
		}
	}

	bar.Total = bar.Total + int64(len(removeTags))

	if len(removeBranches) > 0 || len(removeTags) > 0 {
		log.Printf("Start index deleting for %s:%s/%s (%v) (%v)\n", organization, project, repo.Repository, removeBranches, removeTags)
		g.indexer.DeleteIndexByRefs(organization, project, repo.Repository, removeBranches, removeTags)

		bar.Add(len(removeBranches) + len(removeTags))

		// Save config after deleting index completed
		g.config.DeleteIndexed(organization, project, repo.Repository, removeBranches, removeTags)
	}

	end := time.Now()
	time := (end.Sub(start)).Seconds()

	bar.FinishPrint(fmt.Sprintf("Indexing Complete! [%f seconds] for %s:%s/%s\n", time, organization, project, repo.Repository))
}

func (g *GitImporter) runIndexing(bar *pb.ProgressBar, repo *repo.GitRepo, url string, indexed config.Indexed, branchMap config.BrancheIndexedMap, tagMap config.TagIndexedMap, sizeLimit int64) error {
	// collect create file entries
	createBranches := make(map[string]string)
	updateBranches := make(map[string][2]string)
	for branch, latestCommitID := range branchMap {
		found := false
		for indexedBranch, prevCommitID := range indexed.Branches {
			if branch == indexedBranch {
				found = true
				if latestCommitID == prevCommitID {
					log.Printf("Already up-to-date. %s", getLoggingTag(repo, branch, latestCommitID))
				} else {
					updateBranches[branch] = [2]string{prevCommitID, latestCommitID}
				}
				break
			}
		}
		if !found {
			createBranches[branch] = latestCommitID
		}
	}

	createTags := make(map[string]string)
	updateTags := make(map[string][2]string)
	for tag, latestCommitID := range tagMap {
		found := false
		for indexedTag, prevCommitID := range indexed.Tags {
			if tag == indexedTag {
				found = true
				if latestCommitID == prevCommitID {
					log.Printf("Already up-to-date. %s", getLoggingTag(repo, tag, latestCommitID))
				} else {
					updateTags[tag] = [2]string{prevCommitID, latestCommitID}
				}
				break
			}
		}
		if !found {
			createTags[tag] = latestCommitID
		}
	}

	queue := make(chan indexer.FileIndexOperation, 100)

	// process
	g.UpsertIndex(queue, bar, repo, createBranches, createTags, updateBranches, updateTags, sizeLimit)

	callBatch := func(operations []indexer.FileIndexOperation) {
		err := g.indexer.BatchFileIndex(operations)
		if err != nil {
			errors.Errorf("Batch indexed error: %+v", err)
		} else {
			// fmt.Printf("Batch indexed %d files.\n", len(operations))
		}
		bar.Add(len(operations))
	}

	// batch
	operations := []indexer.FileIndexOperation{}
	var opsSize int64 = 0
	var batchLimitSize int64 = 1024 * 512 // 512KB

	// fmt.Println("start queue reading")

	for op := range queue {
		operations = append(operations, op)
		opsSize += op.FileIndex.Size

		// show progress
		// if len(operations)%80 == 0 {
		// 	fmt.Printf("\n")
		// }
		// fmt.Printf(".")

		if opsSize >= batchLimitSize {
			// fmt.Printf("\n")

			callBatch(operations)

			// reset
			operations = nil
			opsSize = 0
		}
	}

	// remains
	if len(operations) > 0 {
		// fmt.Printf("\n")
		callBatch(operations)
	}

	// Save config after index completed
	err := g.config.UpdateIndexed(config.Indexed{Organization: repo.Organization, Project: repo.Project, Repository: repo.Repository, Branches: branchMap, Tags: tagMap})

	if err != nil {
		return errors.Wrapf(err, "Faild to update indexed.")
	}

	return nil
}

func (g *GitImporter) UpsertIndex(queue chan indexer.FileIndexOperation, bar *pb.ProgressBar, r *repo.GitRepo, branchMap map[string]string, tagMap map[string]string, updateBranchMap map[string][2]string, updateTagMap map[string][2]string, sizeLimit int64) error {
	addFiles, err := r.GetFileEntriesMap(branchMap, tagMap)
	if err != nil {
		return errors.Wrapf(err, "Failed to get file entries. branches: %v tags: %v", branchMap, tagMap)
	}

	updateAddFiles, delFiles, err := r.GetDiffFileEntriesMap(updateBranchMap, updateTagMap)
	if err != nil {
		return errors.Wrapf(err, "Failed to get diff. branches: %v tags: %v", branchMap, tagMap)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		g.handleAddFiles(queue, bar, r, addFiles, sizeLimit)
		g.handleAddFiles(queue, bar, r, updateAddFiles, sizeLimit)
		g.handleDelFiles(queue, bar, r, delFiles)
	}()

	go func() {
		wg.Wait()
		close(queue)
	}()

	return nil
}

func (g *GitImporter) handleAddFiles(queue chan indexer.FileIndexOperation, bar *pb.ProgressBar, r *repo.GitRepo, addFiles map[string]repo.GitFile, sizeLimit int64) {
	if len(addFiles) == 0 {
		return
	}

	var wg sync.WaitGroup
	scanQueue := make(chan ScannedFile, 5)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go scanFiles(&wg, scanQueue, queue, g, r, bar)
	}

	for blob, file := range addFiles {
		// check size
		if sizeLimit > 0 && file.Size > sizeLimit {
			continue
		}

		scanQueue <- ScannedFile{Blob: blob, GitFile: file}
	}

	close(scanQueue)

	wg.Wait()
}

type ScannedFile struct {
	Blob    string
	GitFile repo.GitFile
}

func scanFiles(wg *sync.WaitGroup, scanQueue chan ScannedFile, queue chan indexer.FileIndexOperation, g *GitImporter, r *repo.GitRepo, bar *pb.ProgressBar) {
	defer wg.Done()

	for {
		scannedFile, ok := <-scanQueue
		if !ok {
			return
		}

		blob := scannedFile.Blob
		file := scannedFile.GitFile

		for path, loc := range file.Locations {
			// check contentType and retrive the file content
			// !! this will be heavy process !!
			contentType, content, err := g.parseContent(r, blob)
			if err != nil {
				log.Printf("Failed to parse file. [%s] - %s %+v\n", blob, path, err)
				continue
				// return errors.Wrapf(err, "Failed to parse file. [%s] - %s\n", blob, path)
			}

			// @TODO Extract text from binary in the future?
			if !strings.HasPrefix(contentType, "text/") && contentType != "application/octet-stream" {
				continue
			}

			text, encoding, err := readText(content)
			if err != nil {
				text = string(content)
				encoding = "utf8"
			}

			fileIndex := indexer.FileIndex{
				Metadata: indexer.Metadata{
					Blob:         blob,
					Organization: r.Organization,
					Project:      r.Project,
					Repository:   r.Repository,
					Branches:     loc.Branches,
					Tags:         loc.Tags,
					Path:         path,
					Ext:          indexer.GetExt(path),
					Encoding:     encoding,
					Size:         file.Size,
				},
				Content: text,
			}

			bar.Total = bar.Total + 1

			queue <- indexer.FileIndexOperation{Method: indexer.ADD, FileIndex: fileIndex}
		}
	}
}

// How to detect encoding
// http://qiita.com/nobuhito/items/ff782f64e32f7ed95e43
func readText(body []byte) (string, string, error) {
	var f []byte
	encodings := []string{"shift_jis", "utf8"}
	var enc string
	for i := range encodings {
		enc = encodings[i]
		if enc != "" {
			ee, _ := charset.Lookup(enc)
			if ee == nil {
				continue
			}
			var buf bytes.Buffer
			ic := transform.NewWriter(&buf, ee.NewDecoder())
			_, err := ic.Write(body)
			if err != nil {
				continue
			}
			err = ic.Close()
			if err != nil {
				continue
			}
			f = buf.Bytes()
			break
		}
	}
	return string(f), enc, nil
}

func (g *GitImporter) handleDelFiles(queue chan indexer.FileIndexOperation, bar *pb.ProgressBar, r *repo.GitRepo, delFiles map[string]repo.GitFile) {
	for blob, file := range delFiles {
		for path, loc := range file.Locations {
			fileIndex := indexer.FileIndex{
				Metadata: indexer.Metadata{
					Blob:         blob,
					Organization: r.Organization,
					Project:      r.Project,
					Repository:   r.Repository,
					Branches:     loc.Branches,
					Tags:         loc.Tags,
					Path:         path,
				},
			}

			bar.Total = bar.Total + 1

			// Delete index
			queue <- indexer.FileIndexOperation{Method: indexer.DELETE, FileIndex: fileIndex}
		}
	}
}

func (g *GitImporter) parseContent(repo *repo.GitRepo, blob string) (string, []byte, error) {
	contentType, content, err := repo.DetectBlobContentType(blob)
	if err != nil {
		return "", nil, errors.Wrapf(err, "Failed to read contentType. %s", blob)
	}
	return contentType, content, nil
}

func getLoggingTag(repo *repo.GitRepo, ref string, commitId string) string {
	tag := fmt.Sprintf("%s:%s/%s (%s) [%s]", repo.Organization, repo.Project, repo.Repository, ref, commitId)
	return tag
}
