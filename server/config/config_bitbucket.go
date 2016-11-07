package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type BitBucketResponse struct {
	Size          int  `json:"size"`
	Limit         int  `json:"limit"`
	IsLastPage    bool `json:"isLastPage"`
	Start         int  `json:"start"`
	NextPageStart int  `json:"nextPageStart"`
}

type BitBucketRepositories struct {
	BitBucketResponse
	Values []BitBucketRepository `json:"values"`
}

type BitBucketProject struct {
	Key         string        `json:"key"`
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Public      bool          `json:"public"`
	Type        string        `json:"type"`
	Link        BitBucketLink `json:"link"`
	Links       struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

type BitBucketLink struct {
	URL string `json:"url"`
	Rel string `json:"rel"`
}

type BitBucketClone struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

type BitBucketLinks struct {
	Clone []BitBucketClone `json:"clone"`
	Self  []BitBucketLink  `json:"self"`
}

type BitBucketRepository struct {
	Slug          string           `json:"slug"`
	ID            int              `json:"id"`
	Name          string           `json:"name"`
	ScmID         string           `json:"scmId"`
	State         string           `json:"state"`
	StatusMessage string           `json:"statusMessage"`
	Forkable      bool             `json:"forkable"`
	Project       BitBucketProject `json:"project"`
	Public        bool             `json:"public"`
	CloneURL      string           `json:"cloneUrl"`
	Link          BitBucketLink    `json:"link"`
	Links         BitBucketLinks   `json:"links"`
}

type BitbucketOrganizationSetting struct {
	OrganizationSetting
}

func NewBitbucketOrganizationSetting(o OrganizationSetting) SyncSetting {
	return &BitbucketOrganizationSetting{o}
}

func (b *BitbucketOrganizationSetting) AddRepository(project string, repositoryUrl string) error {
	return b.AddRepository(project, repositoryUrl)
}

func (b *BitbucketOrganizationSetting) SyncSCM() error {

	projects := make(map[string]*ProjectSetting)
	start := 0

	for {
		client := &http.Client{}
		req, err := http.NewRequest("GET", b.Scm["url"]+"/rest/api/1.0/repos?start="+strconv.Itoa(start), nil)
		req.SetBasicAuth(b.Scm["user"], b.Scm["password"])
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}
		bodyText, err := ioutil.ReadAll(resp.Body)
		var res BitBucketRepositories
		json.Unmarshal(bodyText, &res)

		for i := range res.Values {
			r := res.Values[i]

			p, ok := projects[r.Project.Name]
			if !ok {
				projects[r.Project.Name] = &ProjectSetting{Name: r.Project.Name, Repositories: []RepositorySetting{RepositorySetting{Url: r.CloneURL}}}
			} else {
				p.Repositories = append(projects[r.Project.Name].Repositories, RepositorySetting{Url: r.CloneURL})
			}
		}

		if res.IsLastPage {
			break
		}

		start = res.NextPageStart
	}

	// clear
	b.Projects = nil

	for _, v := range projects {
		b.Projects = append(b.Projects, *v)
	}

	return nil
}