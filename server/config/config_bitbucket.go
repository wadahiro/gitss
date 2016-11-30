package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"regexp"

	"github.com/pkg/errors"
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

func (b *BitbucketOrganizationSetting) JSON() ([]byte, error) {
	setting := &struct {
		Name string            `json:"name"`
		Scm  map[string]string `json:"scm,omitempty"`
	}{Name: b.Name, Scm: b.Scm}

	bytes, err := json.MarshalIndent(setting, "", "  ")
	if err != nil {
		return nil, err
	}
	return bytes, err
}

func regex(pattern string, isInclude bool) *regexp.Regexp {
	if pattern != "" {
		r, err := regexp.Compile(pattern)
		if err == nil {
			return r
		}
	}

	if isInclude {
		return MATCH_ALL
	} else {
		return nil
	}
}

func (b *BitbucketOrganizationSetting) GetRefFilters(project string, repository string) (*regexp.Regexp, *regexp.Regexp, *regexp.Regexp, *regexp.Regexp) {
	includeBranches := regex(b.IncludeBranches, true)
	excludeBranches := regex(b.ExcludeBranches, false)
	includeTags := regex(b.IncludeTags, true)
	excludeTags := regex(b.ExcludeTags, false)

	return includeBranches, excludeBranches, includeTags, excludeTags
}

func (b *BitbucketOrganizationSetting) SyncSCM() error {

	projects := make(map[string]*ProjectSetting)
	start := 0

	log.Println("Fetching repositories from bitbucket server: ", b.Scm["url"])

	includeProjects := regex(b.Scm["includeProjects"], true)
	excludeProjects := regex(b.Scm["excludeProjects"], false)
	includeRepositories := regex(b.Scm["includeRepositories"], true)
	excludeRepositories := regex(b.Scm["excludeRepositories"], false)

	for {
		client := &http.Client{}
		req, err := http.NewRequest("GET", b.Scm["url"]+"/rest/api/1.0/repos?start="+strconv.Itoa(start), nil)
		req.SetBasicAuth(b.Scm["user"], b.Scm["password"])
		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrapf(err, "Failed to fetch repositories from bitbucket server: %s", b.Scm["url"])
		}
		bodyText, err := ioutil.ReadAll(resp.Body)

		var res BitBucketRepositories
		json.Unmarshal(bodyText, &res)

		for i := range res.Values {
			r := res.Values[i]
			s := strings.Split(r.CloneURL, "@")

			password, _ := b.Scm["password"]
			password = strings.Replace(password, "@", "%40", -1)
			cloneUrl := s[0] + ":" + password + "@" + s[1]

			// include/exclude project
			if !includeProjects.MatchString(r.Project.Key) ||
				(excludeProjects != nil && excludeProjects.MatchString(r.Project.Key)) {
				log.Printf("%s:%s is ignored.\n", b.GetName(), r.Project.Key)
				continue
			}

			// log.Println("project ok")

			// include/exclude repository
			rs := RepositorySetting{Url: cloneUrl}
			rn := rs.GetName()
			if !includeRepositories.MatchString(rn) ||
				(excludeRepositories != nil && excludeRepositories.MatchString(rn)) {
				log.Printf("%s:%s/%s is ignored.\n", b.GetName(), r.Project.Key, rn)
				continue
			}

			// log.Println("repository ok")

			p, ok := projects[r.Project.Key]
			if !ok {
				projects[r.Project.Key] = &ProjectSetting{Name: r.Project.Key, Repositories: []RepositorySetting{RepositorySetting{Url: cloneUrl}}}
			} else {
				p.Repositories = append(projects[r.Project.Key].Repositories, RepositorySetting{Url: cloneUrl})
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

	// log.Printf("Updated: %#v\n", b)

	return nil
}
