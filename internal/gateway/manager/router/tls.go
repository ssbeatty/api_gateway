package router

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/pkg/types"
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg/errors"
	"sync"
)

func generateCertificate(cfg *config.TLS) (*tls.Certificate, error) {
	var (
		err error
		crt tls.Certificate
	)

	if cfg.Type == config.TLSTypePath {
		crt, err = tls.LoadX509KeyPair(cfg.CsrFile, cfg.KeyFile)
	} else {
		crt, err = tls.X509KeyPair([]byte(cfg.CsrFile), []byte(cfg.KeyFile))
	}
	if err != nil {
		return nil, err
	}

	return &crt, nil
}

func generateTLSConfig(cfg *config.TLS) (*tls.Config, error) {

	crt, err := generateCertificate(cfg)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		ClientAuth: cfg.ClientAuth,
		Certificates: []tls.Certificate{
			*crt,
		},
	}
	if len(cfg.CaFiles) > 0 {
		pool := x509.NewCertPool()
		for _, caFile := range cfg.CaFiles {
			ok := pool.AppendCertsFromPEM([]byte(caFile))
			if !ok {
				return nil, errors.New("invalid certificate(s) content")
			}
		}
		tlsConfig.ClientCAs = pool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return tlsConfig, nil
}

func (f *Factory) generateHTTPSConfig(rtConf *config.Endpoint) (*tls.Config, error) {
	var (
		firstUseful config.TLS
		once        sync.Once
	)
	tlsCfgs := make(map[string]config.TLS)

	for _, router := range getRouters(rtConf, true) {
		once.Do(func() {
			firstUseful = router.TLSConfig
		})
		tlsCfgs[router.Host] = router.TLSConfig
	}

	if len(tlsCfgs) == 0 {
		return nil, errors.New("Empty tls config")
	}
	// todo use first certs
	tlsConfig, err := generateTLSConfig(&firstUseful)
	if err != nil {
		return nil, err
	}

	tlsConfig.GetCertificate = func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		domainToCheck := types.CanonicalDomain(clientHello.ServerName)

		bestCertificate, ok := tlsCfgs[domainToCheck]
		if ok {
			return generateCertificate(&bestCertificate)
		}

		return nil, errors.New("Unknown cert")
	}

	return tlsConfig, nil
}
