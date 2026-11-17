package docs

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/* openapi.yml
var DocsFS embed.FS

func Handler() (http.Handler, error) {
	sub, err := fs.Sub(DocsFS, "static")
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(sub)), nil
}

func OpenAPIHandler() http.Handler {
	return http.FileServer(http.FS(DocsFS))
}
