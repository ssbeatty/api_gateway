package stripprefixregex

import (
	"api_gateway/pkg/middlewares/logs"
	"api_gateway/pkg/middlewares/stripprefix"
	"context"
	"net/http"
	"regexp"
	"strings"
)

const (
	TypeName = "StripPrefixRegex"
)

// StripPrefixRegex is a middleware used to strip prefix from an URL request.
type stripPrefixRegex struct {
	next        http.Handler
	expressions []*regexp.Regexp
	name        string
}

type StripPrefixRegex struct {
	// Regex defines the regular expression to match the path prefix from the request URL.
	Regex []string `json:"regex,omitempty" toml:"regex,omitempty" yaml:"regex,omitempty" export:"true"`
}

func (b *StripPrefixRegex) Schema() (string, error) {

	return "", nil
}

// New builds a new StripPrefixRegex middleware.
func New(ctx context.Context, next http.Handler, config *StripPrefixRegex, name string) (http.Handler, error) {
	logs.GetLogger(ctx, name, TypeName).Debug().Msg("Creating middleware")

	stripPrefix := stripPrefixRegex{
		next: next,
		name: name,
	}

	for _, exp := range config.Regex {
		reg, err := regexp.Compile(strings.TrimSpace(exp))
		if err != nil {
			return nil, err
		}
		stripPrefix.expressions = append(stripPrefix.expressions, reg)
	}

	return &stripPrefix, nil
}

func (s *stripPrefixRegex) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, exp := range s.expressions {
		parts := exp.FindStringSubmatch(req.URL.Path)
		if len(parts) > 0 && len(parts[0]) > 0 {
			prefix := parts[0]
			if !strings.HasPrefix(req.URL.Path, prefix) {
				continue
			}

			req.Header.Add(stripprefix.ForwardedPrefixHeader, prefix)

			req.URL.Path = ensureLeadingSlash(strings.Replace(req.URL.Path, prefix, "", 1))
			if req.URL.RawPath != "" {
				req.URL.RawPath = ensureLeadingSlash(req.URL.RawPath[len(prefix):])
			}

			req.RequestURI = req.URL.RequestURI()
			s.next.ServeHTTP(rw, req)
			return
		}
	}

	s.next.ServeHTTP(rw, req)
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
