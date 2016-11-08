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
		t.Errorf("Unexpected returnd err %+v", err)
	}

	if entries == nil {
		t.Errorf("Unexpected returnd nil entries")
	}
}
