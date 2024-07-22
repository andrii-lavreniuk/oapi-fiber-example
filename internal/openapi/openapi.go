package openapi

import (
	"embed"
	"strings"

	"github.com/google/wire"
	"github.com/mvrilo/go-redoc"
)

var (
	ProviderSet = wire.NewSet(New)

	//go:embed openapi.yaml
	spec embed.FS
)

type Spec struct{}

// New returns a new instance of OpenAPI specification.
func New() *Spec {
	return &Spec{}
}

func (s *Spec) Redoc(path string) redoc.Redoc {
	path = strings.Trim(path, "/")

	return redoc.Redoc{
		Title:    "API V1",
		SpecFS:   &spec,
		SpecFile: "openapi.yaml",
		DocsPath: "/" + path,
		SpecPath: "/" + path + "/openapi.yaml",
	}
}
