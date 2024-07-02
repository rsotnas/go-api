package main

import (
	"context"
	"fmt"
)

func otherProcess(ctx context.Context) {
	fmt.Printf("\nReceived in context: %s and %s\n",
		ctx.Value("parameter"), ctx.Value("otherParameter"))
}

func process(ctx context.Context) {
	fmt.Printf("\nReceived in context: %s\n", ctx.Value("parameter"))
	myCtx := context.WithValue(ctx, "otherParameter", "other value")
	otherProcess(myCtx)
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "parameter", "value")
	process(ctx)
}
