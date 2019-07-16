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
			// mu := &sync.Mutex{}
			res := make(chan interface{}, 3)
			wg.Add(3)
			for i := 0; i < 3; i++ {
				go func() {
					for val := range in {
						data := val.(int) + 1
						// crc32md5 := make(chan int, 1)
						// go func(res chan<- int) {
						// 	mu.Lock()
						time.Sleep(1 * time.Second)
						res <- data
						// 	mu.Unlock()
						// }(crc32md5)
						// res <- (<-crc32md5)
					}
					wg.Done()
				}()
			}
			start := time.Now()
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
