package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type job func(in, out chan interface{})

func main() {
	runtime.GOMAXPROCS(0)

	jobs := []job{
		job(func(in, out chan interface{}) {
			out <- 1
			out <- 2
			out <- 3
		}),
		job(func(in, out chan interface{}) {
			wg := &sync.WaitGroup{}
			res := make(chan interface{}, 100)
			start := time.Now()
			for val := range in {
				data := val.(int) + 1
				wg.Add(1)
				go func() {
					time.Sleep(time.Second)
					res <- data
					wg.Done()
				}()
			}
			wg.Wait()
			close(res)
			end := time.Since(start)
			fmt.Println(end)
			for r := range res {
				out <- r
			}
		}),
		job(func(in, out chan interface{}) {
			for val := range in {
				out <- val.(int) * 2
			}
		}),
		job(func(in, out chan interface{}) {
			for val := range in {
				out <- val
			}
		}),
		job(func(in, out chan interface{}) {
			for val := range in {
				fmt.Println(val)
			}
		}),
	}

	in := make(chan interface{})
	out := make(chan interface{})

	for i := 0; i < len(jobs); i++ {
		in, out = out, make(chan interface{})
		go func(i int, in, out chan interface{}) {
			jobs[i](in, out)
			close(out)
		}(i, in, out)
	}

	fmt.Scanln()
}
