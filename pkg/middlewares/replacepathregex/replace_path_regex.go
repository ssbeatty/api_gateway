package replacepathregex

import (
	"api_gateway/pkg/middlewares/logs"
	"api_gateway/pkg/middlewares/replacepath"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const TypeName = "ReplacePathRegex"

// ReplacePathRegex is a middleware used to replace the path of a URL request with a regular expression.
type replacePathRegex struct {
	next        http.Handler
	regexp      *regexp.Regexp
	replacement string
	name        string
}

type ReplacePathRegex struct {
	// Regex defines the regular expression used to match and capture the path from the request URL.
	Regex string `json:"regex,omitempty"`
	// Replacement defines the replacement path format, which can include captured variables.
	Replacement string `json:"replacement,omitempty"`
}

func (b *ReplacePathRegex) Schema() (string, error) {

	return "", nil
}

// New creates a new replace path regex middleware.
func New(ctx context.Context, next http.Handler, config *ReplacePathRegex, name string) (http.Handler, error) {
	logs.GetLogger(ctx, name, TypeName).Debug().Msg("Creating middleware")

	exp, err := regexp.Compile(strings.TrimSpace(config.Regex))
	if err != nil {
		return nil, fmt.Errorf("error compiling regular expression %s: %w", config.Regex, err)
	}

	return &replacePathRegex{
		regexp:      exp,
		replacement: strings.TrimSpace(config.Replacement),
		next:        next,
		name:        name,
	}, nil
}

func (rp *replacePathRegex) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	currentPath := req.URL.RawPath
	if currentPath == "" {
		currentPath = req.URL.EscapedPath()
	}

	if rp.regexp != nil && len(rp.replacement) > 0 && rp.regexp.MatchString(currentPath) {
		req.Header.Add(replacepath.ReplacedPathHeader, currentPath)
		req.URL.RawPath = rp.regexp.ReplaceAllString(currentPath, rp.replacement)

		// as replacement can introduce escaped characters
		// Path must remain an unescaped version of RawPath
		// Doesn't handle multiple times encoded replacement (`/` => `%2F` => `%252F` => ...)
		var err error
		req.URL.Path, err = url.PathUnescape(req.URL.RawPath)
		if err != nil {
			logs.GetLogger(context.Background(), rp.name, TypeName).Error().Err(err).Send()
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		req.RequestURI = req.URL.RequestURI()
	}

	rp.next.ServeHTTP(rw, req)
}
