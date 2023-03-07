package loadbalancer

import (
	"errors"
	"strconv"
	"strings"
)

type WeightRoundRobinBalance struct {
	rss  []*WeightNode
	conf LoadBalanceConf
}

type WeightNode struct {
	addr            string
	weight          int
	currentWeight   int
	effectiveWeight int
}

func (r *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("param len need 2")
	}
	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}
	node := &WeightNode{addr: params[0], weight: int(parInt)}
	node.effectiveWeight = node.weight
	r.rss = append(r.rss, node)
	return nil
}

func (r *WeightRoundRobinBalance) Next() string {
	total := 0
	var best *WeightNode
	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]
		total += w.effectiveWeight

		w.currentWeight += w.effectiveWeight

		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}
		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}
	if best == nil {
		return ""
	}
	best.currentWeight -= total
	return best.addr
}

func (r *WeightRoundRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *WeightRoundRobinBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

func (r *WeightRoundRobinBalance) Update() {
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
