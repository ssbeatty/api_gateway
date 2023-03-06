package tcp

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var tcpFuncs = map[string]func(*matchersTree, ...string) error{
	"HostSNI":       expect1Parameter(hostSNI),
	"HostSNIRegexp": expect1Parameter(hostSNIRegexp),
}

func expect1Parameter(fn func(*matchersTree, ...string) error) func(*matchersTree, ...string) error {
	return func(route *matchersTree, s ...string) error {
		if len(s) != 1 {
			return fmt.Errorf("unexpected number of parameters; got %d, expected 1", len(s))
		}

		return fn(route, s...)
	}
}

var almostFQDN = regexp.MustCompile(`^[[:alnum:]\.-]+$`)

// hostSNI checks if the SNI Host of the connection match the matcher host.
func hostSNI(tree *matchersTree, hosts ...string) error {
	host := hosts[0]

	if host == "*" {
		// Since a HostSNI(`*`) rule has been provided as catchAll for non-TLS TCP,
		// it allows matching with an empty serverName.
		tree.matcher = func(meta ConnData) bool { return true }
		return nil
	}

	if !almostFQDN.MatchString(host) {
		return fmt.Errorf("invalid value for HostSNI matcher, %q is not a valid hostname", host)
	}

	tree.matcher = func(meta ConnData) bool {
		if meta.serverName == "" {
			return false
		}

		if host == meta.serverName {
			return true
		}

		// trim trailing period in case of FQDN
		host = strings.TrimSuffix(host, ".")

		return host == meta.serverName
	}

	return nil
}

// hostSNIRegexp checks if the SNI Host of the connection matches the matcher host regexp.
func hostSNIRegexp(tree *matchersTree, templates ...string) error {
	template := templates[0]

	if !isASCII(template) {
		return fmt.Errorf("invalid value for HostSNIRegexp matcher, %q is not a valid hostname", template)
	}

	re, err := regexp.Compile(template)
	if err != nil {
		return fmt.Errorf("compiling HostSNIRegexp matcher: %w", err)
	}

	tree.matcher = func(meta ConnData) bool {
		return re.MatchString(meta.serverName)
	}

	return nil
}

// isASCII checks if the given string contains only ASCII characters.
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}

	return true
}
