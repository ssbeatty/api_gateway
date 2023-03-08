package ipallowlist

import (
	"api_gateway/pkg/ip"
	"api_gateway/pkg/middlewares/logs"
	"api_gateway/pkg/tcp"
	"context"
	"errors"
	"fmt"
)

const (
	TypeName = "IPAllowListerTCP"
)

// ipAllowLister is a middleware that provides Checks of the Requesting IP against a set of Allowlists.
type ipAllowLister struct {
	next        tcp.Handler
	allowLister *ip.Checker
	name        string
}

type TCPIPAllowList struct {
	// SourceRange defines the allowed IPs (or ranges of allowed IPs by using CIDR notation).
	SourceRange []string `json:"sourceRange,omitempty"`
}

func (b *TCPIPAllowList) Validator() error {

	return nil
}

// New builds a new TCP IPAllowLister given a list of CIDR-Strings to allow.
func New(ctx context.Context, next tcp.Handler, config *TCPIPAllowList, name string) (tcp.Handler, error) {
	logger := logs.GetLogger(ctx, name, TypeName)
	logger.Debug().Msg("Creating middleware IPAllowListerTCP")

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

func (al *ipAllowLister) ServeTCP(conn tcp.WriteCloser) {
	logger := logs.GetLogger(context.Background(), al.name, TypeName)

	addr := conn.RemoteAddr().String()

	err := al.allowLister.IsAuthorized(addr)
	if err != nil {
		logger.Error().Err(err).Msgf("Connection from %s rejected", addr)
		conn.Close()
		return
	}

	logger.Debug().Msgf("Connection from %s accepted", addr)

	al.next.ServeTCP(conn)
}
