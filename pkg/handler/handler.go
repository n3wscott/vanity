package handler

import (
	"github.com/n3wscott/vanity/pkg/vanity"
	"html/template"
	"net/http"
)

func New(config *vanity.Config) http.Handler {
	return &handler{
		cfg: config,
	}
}

type handler struct {
	cfg   *vanity.Config
	paths vanity.PathConfigs
}

var indexTmpl = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
	<body>
		<h1>{{.Host}}</h1>
		<ul>
			{{range .Handlers}}<li><a href="https://pkg.go.dev/{{.}}">{{.}}</a></li>{{end}}
		</ul>
	</body>
</html>
`))

var vanityTmpl = template.Must(template.New("vanity").Parse(`<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
		<meta name="go-import" content="{{.Import}} {{.VCS}} {{.Repo}}">
		<meta name="go-source" content="{{.Import}} {{.Display}}">
		<meta http-equiv="refresh" content="0; url=https://pkg.go.dev/{{.Import}}/{{.Subpath}}">
	</head>
	<body>
		Nothing to see here; <a href="https://pkg.go.dev/{{.Import}}/{{.Subpath}}">see the package on pkg.go.dev</a>.
	</body>
</html>`))

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.paths == nil {
		paths, err := h.cfg.GeneratePathConfigs()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.paths = paths
	}

	current := r.URL.Path
	pc, subpath := h.paths.Find(current)
	if pc == nil && current == "/" {
		h.serveIndex(w, r)
		return
	}
	if pc == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Cache-Control", h.cfg.CacheControl())
	if err := vanityTmpl.Execute(w, struct {
		Import  string
		Subpath string
		Repo    string
		Display string
		VCS     string
	}{
		Import:  h.Host(r) + pc.Path,
		Subpath: subpath,
		Repo:    pc.Repo,
		Display: pc.Display,
		VCS:     pc.VCS,
	}); err != nil {
		http.Error(w, "cannot render the page", http.StatusInternalServerError)
	}
}

func (h *handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	host := h.Host(r)
	handlers := make([]string, len(h.paths))
	for i, p := range h.paths {
		handlers[i] = host + p.Path
	}
	if err := indexTmpl.Execute(w, struct {
		Host     string
		Handlers []string
	}{
		Host:     host,
		Handlers: handlers,
	}); err != nil {
		http.Error(w, "cannot render the page", http.StatusInternalServerError)
	}
}

func (h *handler) Host(r *http.Request) string {
	host := h.cfg.Host
	if host == "" {
		host = r.Host
	}
	return host
}
