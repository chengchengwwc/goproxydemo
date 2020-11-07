package load_balance

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

/* 哈希一致性轮询 */

type Hash func(data []byte) uint32

type UInt32Slice []uint32

func (s UInt32Slice) Len() int {
	return len(s)
}

func (s UInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s UInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ConsistentHashBalance struct {
	mux      sync.RWMutex
	hash     Hash
	replicas int               //复制因子
	keys     UInt32Slice       //排序的hash切片
	hashMap  map[uint32]string //节点哈希和key的map
	conf     LoadBalanceConf
}

func NewConsistentHashBalance(replicas int, fn Hash) *ConsistentHashBalance {
	m := &ConsistentHashBalance{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[uint32]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 验证是否为空
func (c *ConsistentHashBalance) IsEmpty() bool {
	return len(c.keys) == 0
}

func (c *ConsistentHashBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	c.mux.Lock()
	defer c.mux.Unlock()
	for i := 0; i < c.replicas; i++ {
		hash := c.hash([]byte(strconv.Itoa(i) + addr))
		c.keys = append(c.keys, hash)
		c.hashMap[hash] = addr
	}
	sort.Sort(c.keys)
	return nil
}

/*  根据绑定的对象获取靠近的节点 */
func (c *ConsistentHashBalance) Get(key string) (string, error) {
	if c.IsEmpty() {
		return "", errors.New("bad")
	}
	hash := c.hash([]byte(key))
	idx := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hash
	})

	if idx == len(c.keys) {
		idx = 0
	}
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.hashMap[c.keys[idx]], nil

}
