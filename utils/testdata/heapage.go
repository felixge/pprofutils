package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	go oneSec()
	go twoSec()
	time.Sleep(10 * time.Second)

	file, err := os.Create("./heapage.pprof")
	if err != nil {
		return err
	}
	defer file.Close()

	//runtime.GC()
	//runtime.GC()
	if err := pprof.Lookup("allocs").WriteTo(file, 0); err != nil {
		return err
	}
	return file.Close()
}

func oneSec() {
	allocLoop(10000, 10*1024, time.Second)
}

func twoSec() {
	allocLoop(10000, 10*1024, 2*time.Second)
}

func allocLoop(perSec int, size int, age time.Duration) {
	var alive []*Alloc
	for {
		now := time.Now()
		for i := 0; i < perSec/100; i++ {
			alive = append(alive, NewAlloc(size, now.Add(age)))
		}
		for len(alive) > 0 && now.After(alive[0].Expires) {
			alive = alive[1:]
		}
		time.Sleep(10 * time.Millisecond)
	}
}

type Alloc struct {
	Data    []byte
	Expires time.Time
}

//go:noinline
func NewAlloc(size int, expires time.Time) *Alloc {
	return &Alloc{
		Data:    make([]byte, size),
		Expires: expires,
	}
}
