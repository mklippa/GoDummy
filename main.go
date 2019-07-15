package main

import (
	"fmt"
	"runtime"
	"sync"
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
				out <- val.(int) + 1
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

	gen := func() chan interface{} {
		out := make(chan interface{})
		go func() {
			jobs[0](nil, out)
			close(out)
		}()
		return out
	}

	in := gen()

	calc := func(i int, in chan interface{}) chan interface{} {
		out := make(chan interface{})
		go func() {
			jobs[i](in, out)
			close(out)
		}()
		return out
	}

	outs := make([]chan interface{}, 0)
	for i := 0; i < 5; i++ {
		outs = append(outs, calc(2, calc(1, in)))
	}

	merge := func(cs ...chan interface{}) chan interface{} {
		var wg sync.WaitGroup
		out := make(chan interface{})

		// Start an output goroutine for each input channel in cs.  output
		// copies values from c to out until c is closed, then calls wg.Done.
		output := func(c <-chan interface{}) {
			for n := range c {
				out <- n
			}
			wg.Done()
		}
		wg.Add(len(cs))
		for _, c := range cs {
			go output(c)
		}

		// Start a goroutine to close out once all the output goroutines are
		// done.  This must start after the wg.Add call.
		go func() {
			wg.Wait()
			close(out)
		}()
		return out
	}

	jobs[len(jobs)-1](merge(outs...), nil)

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
