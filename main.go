package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	valueCtx := context.WithValue(ctx, "ddd", "add value")
	go watch(valueCtx)
	time.Sleep(20 * time.Second)
	cancel()
	time.Sleep(2 * time.Second)

}

func watch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Value("ddd"), "is cancel")
		default:
			fmt.Println(ctx.Value("ddd"), "int goroutine")
			time.Sleep(2 * time.Second)
		}
	}
}
