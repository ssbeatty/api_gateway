package loadbalancer

import (
	"errors"
	"strings"
)

type RoundRobinBalance struct {
	curIndex int
	rss      []string
	conf     LoadBalanceConf
}

func (r *RoundRobinBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

func (r *RoundRobinBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	lens := len(r.rss)
	if r.curIndex >= lens {
		r.curIndex = 0
	}
	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr
}

func (r *RoundRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RoundRobinBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

func (r *RoundRobinBalance) Update() {
	if conf, ok := r.conf.(*LoadBalanceCheckConf); ok {
		r.rss = nil
		for _, ip := range conf.GetConf() {
			err := r.Add(strings.Split(ip, ",")...)
			if err != nil {
				return
			}
		}
	}
}
