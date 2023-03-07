package loadbalancer

import "context"

type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)

	Update()
}

type LoadBalanceConf interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf(ctx context.Context)
	UpdateConf(conf []string)
}

type Observer interface {
	Update()
}
