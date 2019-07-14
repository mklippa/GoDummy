package main

import (
	"fmt"
	"runtime"
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
			for val := range in {
				out <- val
			}
		}),
		job(func(in, out chan interface{}) {
			for val := range in {
				out <- val
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

	workerInput := make(chan interface{}, 2)
	combination := make(chan interface{})
	for i := 0; i < 100; i++ {
		go func(in <-chan interface{}) {
			for input := range in {
				inCh := make(chan interface{})
				outCh := make(chan interface{}, 1)
				outCh <- input
				close(outCh)
				for i := 1; i < len(jobs)-1; i++ {
					inCh, outCh = outCh, make(chan interface{})
					go func(i int, in, out chan interface{}) {
						jobs[i](in, out)
						fmt.Println("foo")
						close(out)
					}(i, inCh, outCh)
				}
				combination <- (<-outCh)
			}
		}(workerInput)
	}

	works := make(chan interface{}, 100)
	jobs[0](nil, works)
	close(works)

	for data := range works {
		workerInput <- data
	}
	close(workerInput)

	go func() {
		jobs[len(jobs)-1](combination, nil)
		close(combination)
	}()

	// in := make(chan interface{})
	// out := make(chan interface{})

	// for i := 0; i < len(jobs); i++ {
	// 	in, out = out, make(chan interface{})
	// 	go func(i int, in, out chan interface{}) {
	// 		jobs[i](in, out)
	// 		close(out)
	// 	}(i, in, out)
	// }

	// fmt.Scanln()
}
