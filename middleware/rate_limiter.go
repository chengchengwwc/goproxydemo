package middleware

import (
	"golang.org/x/time/rate"
)

func RateLimiter() func(c *SliceRouterContext) {
	l := rate.NewLimiter(1, 2)
	return func(c *SliceRouterContext) {
		if !l.Allow() {
			c.Abort()
			return
		}
		c.Next()
	}
}
