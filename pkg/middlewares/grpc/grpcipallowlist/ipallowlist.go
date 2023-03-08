package grpcipallowlist

import (
	"api_gateway/pkg/ip"
	"api_gateway/pkg/logs"
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"strings"
)

const (
	TypeName = "IPAllowListerGRPC"
)

func (b *IPAllowList) Schema() (string, error) {

	return "", nil
}

type IPAllowList struct {
	// SourceRange defines the set of allowed IPs (or ranges of allowed IPs by using CIDR notation).
	SourceRange []string `json:"sourceRange,omitempty"`
}

func New(ctx context.Context, config *IPAllowList, name string) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	checker, err := ip.NewChecker(config.SourceRange)
	if err != nil {
		return nil
	}
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}
		peerAddr := peerCtx.Addr.String()
		addrPos := strings.LastIndex(peerAddr, ":")
		clientIP := peerAddr[0:addrPos]
		err = checker.IsAuthorized(clientIP)
		if err := handler(srv, ss); err != nil {
			log.Error().Str(logs.MiddlewareName, name).Err(err).Send()
			return err
		}
		return nil
	}
}
