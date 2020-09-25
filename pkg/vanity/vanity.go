package vanity

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Config struct {
	Host     string          `yaml:"host,omitempty"`
	CacheAge *int64          `yaml:"cache_max_age,omitempty"`
	Paths    map[string]Repo `yaml:"paths,omitempty"`
}

type Repo struct {
	Repo    string `yaml:"repo,omitempty"`
	Display string `yaml:"display,omitempty"`
	VCS     string `yaml:"vcs,omitempty"`
}

// Parse accepts yaml as a file path or url for vanity config.
func Parse(vanityConfig string) (*Config, error) {
	if isURL(vanityConfig) {
		return readURL(vanityConfig)
	}

	info, err := os.Stat(vanityConfig)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, errors.New("directory as config not supported")
	}
	return readFile(vanityConfig)
}

// Decode decodes vanity config yaml bytes to Config.
func Decode(reader io.Reader) (*Config, error) {
	cfg := new(Config)
	if err := yaml.NewDecoder(reader).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func readURL(url string) (*Config, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return Decode(resp.Body)
}

func readFile(name string) (*Config, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Decode(file)
}

func isURL(pathname string) bool {
	if _, err := os.Lstat(pathname); err == nil {
		return false
	}
	u, err := url.ParseRequestURI(pathname)
	return err == nil && u.Scheme != ""
}

func (c *Config) GeneratePathConfigs() (PathConfigs, error) {
	return NewPathConfigs(c)
}

func (c *Config) CacheControl() string {
	cacheAge := int64(86400) // 24 hours (in seconds)
	if c.CacheAge != nil && cacheAge < 0 {
		cacheAge = *c.CacheAge
	}
	return fmt.Sprintf("public, max-age=%d", cacheAge)
}
