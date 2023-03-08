package loadbalancer

import (
	"api_gateway/pkg/safe"
	"context"
	"fmt"
	"net"
	"net/url"
	"reflect"
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
	observers    []Observer
	confIpWeight map[string]checkInfo
	activeList   []string
}

func (s *LoadBalanceCheckConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

func (s *LoadBalanceCheckConf) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

func (s *LoadBalanceCheckConf) GetConf() []string {
	var confList []string
	for _, ip := range s.activeList {
		weight := s.confIpWeight[ip]
		confList = append(confList, fmt.Sprintf("%s,%d", ip, weight.IpWeight))
	}
	return confList
}

func (s *LoadBalanceCheckConf) WatchConf(ctx context.Context) {
	confIpErrNum := map[string]int{}
	for {

		select {
		case <-ctx.Done():
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

func (s *LoadBalanceCheckConf) UpdateConf(conf []string) {
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

func NewLoadBalanceCheckConf(conf map[string]int, pool *safe.Pool) (*LoadBalanceCheckConf, error) {
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
	mConf := &LoadBalanceCheckConf{activeList: aList, confIpWeight: checkInfoCfg}

	pool.GoCtx(mConf.WatchConf)

	return mConf, nil
}
