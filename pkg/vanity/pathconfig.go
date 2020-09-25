package vanity

import (
	"fmt"
	"sort"
	"strings"
)

// This is based on https://github.com/GoogleCloudPlatform/govanityurls/blob/master/handler.go#L173-L225
// and https://github.com/GoogleCloudPlatform/govanityurls/blob/master/handler.go#L64-L90
type PathConfigs []PathConfig

type PathConfig struct {
	Path    string
	Repo    string
	Display string
	VCS     string
}

func NewPathConfigs(cfg *Config) (PathConfigs, error) {
	var paths PathConfigs
	for path, e := range cfg.Paths {
		pc := PathConfig{
			Path:    strings.TrimSuffix(path, "/"),
			Repo:    e.Repo,
			Display: e.Display,
			VCS:     e.VCS,
		}
		switch {
		case e.Display != "":
			// Already filled in.
		case strings.HasPrefix(e.Repo, "https://github.com/"):
			pc.Display = fmt.Sprintf("%v %v/tree/master{/dir} %v/blob/master{/dir}/{file}#L{line}", e.Repo, e.Repo, e.Repo)
		case strings.HasPrefix(e.Repo, "https://bitbucket.org"):
			pc.Display = fmt.Sprintf("%v %v/src/default{/dir} %v/src/default{/dir}/{file}#{file}-{line}", e.Repo, e.Repo, e.Repo)
		}
		switch {
		case e.VCS != "":
			// Already filled in.
			if e.VCS != "bzr" && e.VCS != "git" && e.VCS != "hg" && e.VCS != "svn" {
				return nil, fmt.Errorf("configuration for %v: unknown VCS %s", path, e.VCS)
			}
		case strings.HasPrefix(e.Repo, "https://github.com/"):
			pc.VCS = "git"
		default:
			return nil, fmt.Errorf("configuration for %v: cannot infer VCS from %s", path, e.Repo)
		}
		paths = append(paths, pc)
	}
	sort.Sort(paths)
	return paths, nil
}

func (ps PathConfigs) Len() int {
	return len(ps)
}

func (ps PathConfigs) Less(i, j int) bool {
	return ps[i].Path < ps[j].Path
}

func (ps PathConfigs) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps PathConfigs) Find(path string) (pc *PathConfig, subpath string) {
	// Fast path with binary search to retrieve exact matches
	// e.g. given pset ["/", "/abc", "/xyz"], path "/def" won't match.
	i := sort.Search(len(ps), func(i int) bool {
		return ps[i].Path >= path
	})
	if i < len(ps) && ps[i].Path == path {
		return &ps[i], ""
	}
	if i > 0 && strings.HasPrefix(path, ps[i-1].Path+"/") {
		return &ps[i-1], path[len(ps[i-1].Path)+1:]
	}

	// Slow path, now looking for the longest prefix/shortest subpath i.e.
	// e.g. given pset ["/", "/abc/", "/abc/def/", "/xyz"/]
	//  * query "/abc/foo" returns "/abc/" with a subpath of "foo"
	//  * query "/x" returns "/" with a subpath of "x"
	lenShortestSubpath := len(path)
	var bestMatchConfig *PathConfig

	// After binary search with the >= lexicographic comparison,
	// nothing greater than i will be a prefix of path.
	max := i
	for i := 0; i < max; i++ {
		p := ps[i]
		if len(p.Path) >= len(path) {
			// We previously didn't find the path by search, so any
			// route with equal or greater length is NOT a match.
			continue
		}
		sSubpath := strings.TrimPrefix(path, p.Path)
		if len(sSubpath) < lenShortestSubpath {
			subpath = sSubpath
			lenShortestSubpath = len(sSubpath)
			bestMatchConfig = &ps[i]
		}
	}
	return bestMatchConfig, subpath
}
