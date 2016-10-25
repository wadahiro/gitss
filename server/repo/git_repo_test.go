package repo

import (
	"fmt"
	"testing"

	"github.com/wadahiro/gitss/server/repo"
)

func TestGetFileEntries(t *testing.T) {
	r, _ := repo.NewGitRepo("o", "p", "r", "../../", true)

	entries, _ := r.GetFileEntries("HEAD")

	fmt.Println(len(entries))
}
