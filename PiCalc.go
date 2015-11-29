//Baby steps in Golang
//Pi calculation based on
//http://doc.akka.io/docs/akka/2.0/intro/getting-started-first-scala.html
//Costas Yeditech 2015

package main

import (
	"fmt"
	"math"
	"sync"
)

type f64 float64

//------------------------------------------------------------------------------------------------------------------------------------------------------------
func calculatePiFor(start int, nrOfElements int) f64 {
	var acc f64
	acc = 0
	var fi f64
	var mod_ float64
	for i := start; i < start+nrOfElements; i++ {
		fi = f64(i)
		mod_ = math.Mod(float64(i), 2)

		acc += 4.0 * (1 - (f64(mod_))*2) / (2*fi + 1)
	}

	return acc
}

//------------------------------------------------------------------------------------------------------------------------------------------------------------
func worker(id int, no_of_elements int, jobs <-chan int, results chan<- f64) {
	var value f64
	value = 0
	for j := range jobs {
		value = calculatePiFor(j, no_of_elements)
		results <- value
	}
}

//------------------------------------------------------------------------------------------------------------------------------------------------------------
func main() {
	var message_no int
	var elements_no int
	var workers_no int

	//parametes based on original article
	message_no = 10000
	elements_no = 10000
	workers_no = 4

	jobs := make(chan int, message_no)
	results := make(chan f64)

	for w := 0; w < workers_no; w++ {
		go worker(w, elements_no, jobs, results)
	}

	for j := 0; j < message_no; j++ {
		jobs <- j * elements_no
	}

	close(jobs)

	var pi_result f64
	pi_result = 0

	var wg sync.WaitGroup

	wg.Add(1)

	go func(message_no int, result *f64, results <-chan f64) {
		var values_no int
		values_no = 0
		defer wg.Done()
		for {
			select {
			case j := <-results:
				*result += j
				values_no += 1

				if values_no == message_no {
					return
				}
			default:
			}
		}

	}(message_no, &pi_result, results)

	wg.Wait()
	fmt.Printf("Result = %.16f \n", pi_result)
}
