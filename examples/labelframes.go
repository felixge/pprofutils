//go:build ignore

package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	file, _ := os.Create("./labelframes.in.pprof")
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		labels := pprof.Labels("mylabel", fmt.Sprintf("val=%d", i))
		pprof.Do(ctx, labels, func(_ context.Context) {
			go cpuHog()
		})
	}
	time.Sleep(10 * time.Second)
}

func cpuHog() {
	for i := 0; ; i++ {
		fmt.Fprintf(io.Discard, "%d", i)
	}
}
