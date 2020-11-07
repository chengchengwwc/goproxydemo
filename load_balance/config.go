package load_balance

type LoadBalanceConf interface {
	Attach()
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}
