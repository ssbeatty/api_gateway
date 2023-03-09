package stripprefix

import (
	"api_gateway/pkg/middlewares/logs"
	"context"
	"net/http"
	"strings"
)

const (
	// ForwardedPrefixHeader is the default header to set prefix.
	ForwardedPrefixHeader = "X-Forwarded-Prefix"
	TypeName              = "StripPrefix"
	prefixField           = "prefixPath"
)

// stripPrefix is a middleware used to strip prefix from an URL request.
type stripPrefix struct {
	next     http.Handler
	prefixes []string
	name     string
}

type StripPrefix struct {
	// Prefixes defines the prefixes to strip from the request URL.
	Prefixes []string `json:"prefixes,omitempty" toml:"prefixes,omitempty" yaml:"prefixes,omitempty" export:"true"`
}

func (b *StripPrefix) Schema() (string, error) {

	return "", nil
}

// New creates a new strip prefix middleware.
func New(ctx context.Context, next http.Handler, config *StripPrefix, name string) (http.Handler, error) {
	logs.GetLogger(ctx, name, TypeName).Debug().Strs(prefixField, config.Prefixes).Msg("Creating middleware")
	return &stripPrefix{
		prefixes: config.Prefixes,
		next:     next,
		name:     name,
	}, nil
}

func (s *stripPrefix) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, prefix := range s.prefixes {
		if strings.HasPrefix(req.URL.Path, prefix) {
			req.URL.Path = s.getPrefixStripped(req.URL.Path, prefix)
			if req.URL.RawPath != "" {
				req.URL.RawPath = s.getPrefixStripped(req.URL.RawPath, prefix)
			}
			s.serveRequest(rw, req, strings.TrimSpace(prefix))
			return
		}
	}
	s.next.ServeHTTP(rw, req)
}

func (s *stripPrefix) serveRequest(rw http.ResponseWriter, req *http.Request, prefix string) {
	req.Header.Add(ForwardedPrefixHeader, prefix)
	req.RequestURI = req.URL.RequestURI()
	s.next.ServeHTTP(rw, req)
}

func (s *stripPrefix) getPrefixStripped(urlPath, prefix string) string {
	return ensureLeadingSlash(strings.TrimPrefix(urlPath, prefix))
}

func ensureLeadingSlash(str string) string {
	if str == "" {
		return str
	}

	if str[0] == '/' {
		return str
	}

	return "/" + str
}
