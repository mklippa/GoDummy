package main

import (
	"fmt"
	"runtime"
)

type job func(in, out chan interface{})

func startWorker(in <-chan job) {
	inCh := make(chan interface{}, 0)
	outCh := make(chan interface{}, 0)
	for input := range in {
		input(inCh, outCh)
		inCh <- outCh
		runtime.Gosched() // попробуйте закомментировать
	}
}

func main() {
	runtime.GOMAXPROCS(0)
	workers := make(chan job, 2)
	go startWorker(workers)

	jobs := []job{
		job(func(in, out chan interface{}) { out <- 1 }),
		job(func(in, out chan interface{}) { out <- in }),
		job(func(in, out chan interface{}) { out <- in }),
		job(func(in, out chan interface{}) { out <- in }),
		job(func(in, out chan interface{}) { out <- in }),
		job(func(in, out chan interface{}) { out <- in }),
		job(func(in, out chan interface{}) { <-in }),
	}

	for _, worker := range jobs {
		workers <- worker
	}
	close(workers)

	fmt.Scanln()
}
