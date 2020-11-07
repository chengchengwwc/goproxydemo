package load_balance

import (
	"errors"
	"strconv"
)

type WeightRoundRobinBalance struct {
	curIndex int
	rss      []*WeightNode
	rsw      []int
}

type WeightNode struct {
	addr            string
	weight          int // 当前权重
	currentWeight   int // 节点当前权重
	effectiveWeight int // 有效权重
}

func (r *WeightRoundRobinBalance) Add(parmas ...string) error {
	if len(parmas) == 0 {
		return errors.New("bad")
	}
	parInt, err := strconv.ParseInt(parmas[1], 10, 64)
	if err != nil {
		return err
	}
	node := &WeightNode{addr: parmas[0], weight: int(parInt)}
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
