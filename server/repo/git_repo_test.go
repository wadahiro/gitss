package repo

import (
	"fmt"
	"testing"

	"github.com/wadahiro/gitss/server/config"
	"github.com/wadahiro/gitss/server/repo"
)

func TestGetFileEntries(t *testing.T) {
	r, _ := repo.NewGitRepo("o", "p", "r", "../../", &config.Config{})

	entries, err := r.GetFileEntries("HEAD")

	fmt.Println(len(entries))

	if err != nil {
		t.Errorf("Unexpected returned err %+v", err)
	}

	if entries == nil {
		t.Errorf("Unexpected returned nil entries")
	}
}

func TestGetFileEntriesMap(t *testing.T) {
	r, _ := repo.NewGitRepo("o", "p", "r", "../../", &config.Config{})

	entries, err := r.GetFileEntriesMap([]string{}, []string{})

	fmt.Println("len(entries): ", len(entries))

	if err != nil {
		t.Errorf("Unexpected returned err %+v", err)
	}

	if entries == nil {
		t.Errorf("Unexpected returned nil entries")
	}

	fmt.Println("dump: %v", entries)
}
