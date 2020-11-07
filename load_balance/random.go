package load_balance

import (
	"errors"
	"math/rand"
)

/* 随机轮训 */

type RandomBalance struct {
	curIndex int
	rss      []string
	conf     LoadBalanceConf
}

func (r *RandomBalance) Add(parmas ...string) error {
	if len(parmas) == 0 {
		return errors.New("param len 0")
	}
	addr := parmas[0]
	r.rss = append(r.rss, addr)
	return nil
}

func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}

func (r *RandomBalance) Get(key string) (string, error) {
	return r.Next(), nil
}

func (r *RandomBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}
