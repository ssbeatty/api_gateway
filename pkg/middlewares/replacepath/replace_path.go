package replacepath

import (
	"api_gateway/pkg/middlewares/logs"
	"context"
	"net/http"
	"net/url"
	"strings"
)

const (
	// ReplacedPathHeader is the default header to set the old path to.
	ReplacedPathHeader = "X-Replaced-Path"
	TypeName           = "ReplacePath"
)

// ReplacePath is a middleware used to replace the path of a URL request.
type replacePath struct {
	next    http.Handler
	path    string
	name    string
	replace string
}

type ReplacePath struct {
	// Path defines the path to use as replacement in the request URL.
	Path    string `json:"path,omitempty"`
	Replace string `json:"replace,omitempty"`
}

func (b *ReplacePath) Schema() (string, error) {

	return "", nil
}

// New creates a new replace path middleware.
func New(ctx context.Context, next http.Handler, config *ReplacePath, name string) (http.Handler, error) {
	logs.GetLogger(ctx, name, TypeName).Debug().Msg("Creating middleware")

	return &replacePath{
		next:    next,
		path:    config.Path,
		name:    name,
		replace: config.Replace,
	}, nil
}

func (r *replacePath) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	currentPath := req.URL.RawPath
	if currentPath == "" {
		currentPath = req.URL.EscapedPath()
	}

	req.Header.Add(ReplacedPathHeader, currentPath)

	if r.path == "" {
		req.URL.RawPath = r.path
	} else {
		req.URL.RawPath = strings.Replace(req.URL.RawPath, r.path, r.replace, -1)
	}

	var err error
	req.URL.Path, err = url.PathUnescape(req.URL.RawPath)
	if err != nil {
		logs.GetLogger(context.Background(), r.name, TypeName).Error().Err(err).Send()
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	req.RequestURI = req.URL.RequestURI()

	r.next.ServeHTTP(rw, req)
}
