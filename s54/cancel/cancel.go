package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func raffleNumbers(ctx context.Context, ch chan int) {
	seed1 := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed1)
	number := r.Int()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Cancel")
			return

		case ch <- number:
			number = r.Int()
		}
	}
}

func invoke(ctx context.Context) <-chan int {
	ch := make(chan int)
	go raffleNumbers(ctx, ch)
	return ch
}

func main() {
	ctx := context.Background()
	newCtx, cancel := context.WithCancel(ctx)
	count := 0
	for n := range invoke(newCtx) {
		fmt.Println(n)
		count++
		if count > 2 {
			break
		}
	}
	cancel()
	time.Sleep(1 * time.Second)
	fmt.Println("Main finished")
}
