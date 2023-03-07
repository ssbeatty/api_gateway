package upstream

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/gateway/manager/upstream/loadbalancer"
	"api_gateway/pkg/safe"
	"net/http"
)

type Factory struct {
	staticConfiguration config.Gateway
	routinesPool        *safe.Pool
}

func NewFactory(staticConfiguration config.Gateway, routinesPool *safe.Pool) *Factory {

	return &Factory{
		staticConfiguration: staticConfiguration,
		routinesPool:        routinesPool,
	}
}

func (f *Factory) buildUpstreamLoadBalancer(upstreamConfig *config.Upstream) (loadbalancer.LoadBalance, error) {
	ipConf := map[string]int{}
	for ipIndex, ipItem := range upstreamConfig.Paths {
		if upstreamConfig.Weights == nil {
			ipConf[ipItem] = 50
			continue
		}
		ipConf[ipItem] = upstreamConfig.Weights[ipIndex]
	}

	mConf, err := loadbalancer.NewLoadBalanceCheckConf(ipConf, f.routinesPool)
	if err != nil {
		return nil, err
	}
	lb := loadbalancer.LoadBalanceFactorWithConf(upstreamConfig.LoadBalancerType, mConf)

	return lb, nil
}

func (f *Factory) BuildHttpUpstreamHandler(upstreamConfig *config.Upstream) (http.Handler, error) {
	lb, err := f.buildUpstreamLoadBalancer(upstreamConfig)

	if err != nil {
		return nil, err
	}

	return f.NewLoadBalanceReverseProxy(lb), nil
}
