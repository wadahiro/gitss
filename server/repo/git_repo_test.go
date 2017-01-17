package repo

import (
	// "fmt"
	"testing"

	"github.com/wadahiro/gitss/server/config"
)

func TestGetFileEntries(t *testing.T) {
	r, _ := NewGitRepo("o", "p", "r", "../../", &config.Config{})

	entries, err := r.GetFileEntries("HEAD")

	if err != nil {
		t.Errorf("Unexpected returned err %+v", err)
	}

	if entries == nil {
		t.Errorf("Unexpected returned nil entries")
	}
}

func TestGetFileEntriesMap(t *testing.T) {
	r, _ := NewGitRepo("o", "p", "r", "../../", &config.Config{})

	entries, err := r.GetFileEntriesMap(map[string]string{}, map[string]string{})

	if err != nil {
		t.Errorf("Unexpected returned err %+v", err)
	}

	if entries == nil {
		t.Errorf("Unexpected returned nil entries")
	}

	if len(entries) > 0 {
		t.Errorf("Unexpected returned entries %v", entries)
	}

	branchMap := map[string]string{}
	branchMap["master"] = "625cd386575760dd4804ddbfd37e30c5f2e08fce"
	entries, err = r.GetFileEntriesMap(branchMap, map[string]string{})

	if err != nil {
		t.Errorf("Unexpected returned err %+v", err)
	}

	if entries == nil {
		t.Errorf("Unexpected returned nil entries")
	}

	if len(entries) != 91 {
		t.Errorf("Unexpected len(entries). expected: 91, actual: %v", len(entries))
	}

	entry, ok := entries["446396634db48e7999e7936f4ff1cf432ae68fd0"]
	if !ok {
		t.Errorf("Not found entry: 446396634db48e7999e7936f4ff1cf432ae68fd0")
	}

	if entry.Size != 1214 {
		t.Errorf("Unexpected git file size. expected: 1214, actual: %v", entry.Size)
	}

	location, ok := entry.Locations["client/app/components/Panel.tsx"]
	if !ok {
		t.Errorf("Not found location: client/app/components/Panel.tsx")
	}

	if len(location.Branches) != 1 {
		t.Errorf("Unexpected branch list. expected: master, actual: %v", location.Branches)
	}

	if location.Branches[0] != "master" {
		t.Errorf("Unexpected branch. expected: master, actual: %v", location.Branches[0])
	}
}
