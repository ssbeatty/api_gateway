package ipallowlist

import (
	"api_gateway/pkg/ip"
	"api_gateway/pkg/middlewares/logs"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

const (
	TypeName = "IPAllowLister"
)

// ipAllowLister is a middleware that provides Checks of the Requesting IP against a set of Allowlists.
type ipAllowLister struct {
	next        http.Handler
	allowLister *ip.Checker
	strategy    ip.Strategy
	name        string
}

func (b *IPAllowList) Schema() (string, error) {

	return "", nil
}

type IPAllowList struct {
	// SourceRange defines the set of allowed IPs (or ranges of allowed IPs by using CIDR notation).
	SourceRange []string `json:"sourceRange,omitempty"`
}

// New builds a new IPAllowLister given a list of CIDR-Strings to allow.
func New(ctx context.Context, next http.Handler, config *IPAllowList, name string) (http.Handler, error) {
	logger := logs.GetLogger(ctx, name, TypeName)
	logger.Debug().Msg("Creating middleware")

	if len(config.SourceRange) == 0 {
		return nil, errors.New("sourceRange is empty, IPAllowLister not created")
	}

	checker, err := ip.NewChecker(config.SourceRange)
	if err != nil {
		return nil, fmt.Errorf("cannot parse CIDRs %s: %w", config.SourceRange, err)
	}

	logger.Debug().Msgf("Setting up IPAllowLister with sourceRange: %s", config.SourceRange)

	return &ipAllowLister{
		allowLister: checker,
		next:        next,
		name:        name,
	}, nil
}

func (al *ipAllowLister) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	logger := logs.GetLogger(req.Context(), al.name, TypeName)
	ctx := logger.WithContext(req.Context())

	clientIP := al.strategy.GetIP(req)
	err := al.allowLister.IsAuthorized(clientIP)
	if err != nil {
		msg := fmt.Sprintf("Rejecting IP %s: %v", clientIP, err)
		logger.Debug().Msg(msg)
		reject(ctx, rw)
		return
	}
	logger.Debug().Msgf("Accepting IP %s", clientIP)

	al.next.ServeHTTP(rw, req)
}

func reject(ctx context.Context, rw http.ResponseWriter) {
	statusCode := http.StatusForbidden

	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(http.StatusText(statusCode)))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Send()
	}
}
