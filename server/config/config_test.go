package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLineScanner(t *testing.T) {
	os.MkdirAll("./tmp/conf", 0644)
	c := Config{ConfDir: "./tmp/conf"}

	testConfig := []byte(`{
	"name": "test",
	"projects": [{
		"name": "myproject",
		"repositories": [{
			"url": "http://mygit/myproject/samplerepo.git",
			"refs": [{
				"name": "master",
				"latest": "a42cc944dde55039ff8e6bac07d686b478b151cf"
			}, {
				"name": "develop",
				"latest": "14fdbd6efb075a2ab379ef07975cd8033900366a"
			}]
		}]
	}]
}`)
	err := ioutil.WriteFile("./tmp/conf/test.json", testConfig, os.ModePerm)
	if err != nil {
		t.Errorf("test.json create error, %+v", err)
	}

	defer os.RemoveAll("./tmp")

	s, _ := c.GetSettings()

	if len(s) != 1 {
		t.Errorf("organization length should be 1, %v", len(s))
	}

	if s[0].Name != "test" {
		t.Errorf("organization name should be test, %v", s[0].Name)
	}

	if len(s[0].Projects) != 1 {
		t.Errorf("project length should be 1, %v", s[0].Projects)
	}

	if s[0].Projects[0].Name != "myproject" {
		t.Errorf("project name should be myproject, %v", s[0].Projects[0].Name)
	}

	if len(s[0].Projects[0].Repositories) != 1 {
		t.Errorf("repository length should be 1, %v", len(s[0].Projects[0].Repositories))
	}

	if s[0].Projects[0].Repositories[0].GetName() != "samplerepo" {
		t.Errorf("repository name should be samplerepo, %v", s[0].Projects[0].Repositories[0].GetName())
	}

	if len(s[0].Projects[0].Repositories[0].Refs) != 2 {
		t.Errorf("repository length should be 2, %v", len(s[0].Projects[0].Repositories[0].Refs))
	}
}
