package loadbalancer

import (
	"api_gateway/pkg/safe"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"net/url"
	"reflect"
	"runtime"
	"sort"
	"time"
)

const (
	DefaultCheckMethod    = 0
	DefaultCheckTimeout   = 5
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
)

type checkInfo struct {
	IpWeight int
	Host     string
}

type LoadBalanceCheckConf struct {
	s *loadBalanceCheckConf
}

type loadBalanceCheckConf struct {
	observers    []Observer
	confIpWeight map[string]checkInfo
	activeList   []string
	closer       chan struct{}
}

func (s *LoadBalanceCheckConf) Attach(o Observer) {
	s.s.observers = append(s.s.observers, o)
}

func (s *LoadBalanceCheckConf) GetConf() []string {
	var confList []string
	for _, ip := range s.s.activeList {
		weight := s.s.confIpWeight[ip]
		confList = append(confList, fmt.Sprintf("%s,%d", ip, weight.IpWeight))
	}
	return confList
}

func (s *LoadBalanceCheckConf) WatchConf(ctx context.Context) {
	s.s.WatchConf(ctx)
}

func (s *loadBalanceCheckConf) WatchConf(ctx context.Context) {
	defer log.Debug().Msg("load balance check exit.")
	confIpErrNum := map[string]int{}
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.closer:
			return
		default:
			var changedList []string
			for item, checkInfo := range s.confIpWeight {
				conn, err := net.DialTimeout("tcp", checkInfo.Host, time.Duration(DefaultCheckTimeout)*time.Second)
				if err == nil {
					conn.Close()
					confIpErrNum[item] = 0
				}
				if err != nil {
					confIpErrNum[item] += 1
				}
				if confIpErrNum[item] < DefaultCheckMaxErrNum {
					changedList = append(changedList, item)
				}
			}
			sort.Strings(changedList)
			sort.Strings(s.activeList)
			if !reflect.DeepEqual(changedList, s.activeList) {
				s.UpdateConf(changedList)
			}
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}

	}
}

func (s *loadBalanceCheckConf) UpdateConf(conf []string) {
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

func (s *LoadBalanceCheckConf) UpdateConf(conf []string) {
	s.s.UpdateConf(conf)
}

func NewLoadBalanceCheckConf(conf map[string]int, pool *safe.Pool) (*LoadBalanceCheckConf, error) {

	lc, err := newLoadBalanceCheckConf(conf)

	if err != nil {
		return nil, err
	}

	mConf := &LoadBalanceCheckConf{
		s: lc,
	}

	pool.GoCtx(mConf.WatchConf)

	runtime.SetFinalizer(mConf, func(mConf *LoadBalanceCheckConf) {
		mConf.s.closer <- struct{}{}
	})

	return mConf, nil
}

func newLoadBalanceCheckConf(conf map[string]int) (*loadBalanceCheckConf, error) {
	var (
		aList        []string
		checkInfoCfg = make(map[string]checkInfo)
	)
	for item, w := range conf {
		aList = append(aList, item)
		if uri, err := url.Parse(item); err == nil && uri.Hostname() != "" {
			var port string
			if uri.Port() == "" && uri.Scheme == "https" {
				port = "443"
			} else if uri.Port() == "" && uri.Scheme == "http" {
				port = "80"
			} else {
				port = uri.Port()
			}
			host := fmt.Sprintf("%s:%s", uri.Hostname(), port)
			checkInfoCfg[item] = checkInfo{
				Host:     host,
				IpWeight: w,
			}
		} else {
			host, port, err := net.SplitHostPort(item)
			if err != nil {
				return nil, err
			}
			checkInfoCfg[item] = checkInfo{
				Host:     fmt.Sprintf("%s:%s", host, port),
				IpWeight: w,
			}
		}
	}
	mConf := &loadBalanceCheckConf{
		activeList:   aList,
		confIpWeight: checkInfoCfg,
		closer:       make(chan struct{}, 1),
	}

	return mConf, nil
}
