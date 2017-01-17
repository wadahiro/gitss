package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLineScanner(t *testing.T) {
	os.MkdirAll("./tmp/conf", 0755)
	c := Config{ConfDir: "./tmp/conf"}

	testConfig := []byte(`{
	"name": "test",
	"projects": [{
		"name": "myproject",
		"repositories": [{
			"url": "http://mygit/myproject/samplerepo.git"
		}]
	}]
}`)
	err := ioutil.WriteFile("./tmp/conf/test.json", testConfig, os.ModePerm)
	if err != nil {
		t.Errorf("test.json create error, %+v", err)
	}

	defer os.RemoveAll("./tmp")

	c.reloadSettings()
	s := c.GetSettings()

	if len(s) != 1 {
		t.Errorf("organization length should be 1, %v", len(s))
	}

	if s[0].GetName() != "test" {
		t.Errorf("organization name should be test, %v", s[0].GetName())
	}

	if len(s[0].GetProjects()) != 1 {
		t.Errorf("project length should be 1, %v", s[0].GetProjects())
	}

	if s[0].GetProjects()[0].Name != "myproject" {
		t.Errorf("project name should be myproject, %v", s[0].GetProjects()[0].Name)
	}

	if len(s[0].GetProjects()[0].Repositories) != 1 {
		t.Errorf("repository length should be 1, %v", len(s[0].GetProjects()[0].Repositories))
	}

	if s[0].GetProjects()[0].Repositories[0].GetName() != "samplerepo" {
		t.Errorf("repository name should be samplerepo, %v", s[0].GetProjects()[0].Repositories[0].GetName())
	}
}
