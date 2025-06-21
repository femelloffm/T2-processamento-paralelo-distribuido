package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var n = [3]int{200, 500, 800}
var p = [3]int{1, 2, 4}

const MAX = 999

func main() {
	for _, p := range p {
		runtime.GOMAXPROCS(p)
		for _, z := range n {

			inicio := time.Now()

			result := make(chan int, z)
			canais := make([]chan int, z+1)
			for i := 0; i <= z; i++ {
				canais[i] = make(chan int, 2)
			}

			for i := 0; i < z; i++ {
				go cellSorter(i, canais[i], canais[i+1], result, MAX)
			}

			rand.Seed(time.Now().UnixNano())
			for i := 0; i < z; i++ {
				valor := rand.Intn(MAX) - rand.Intn(MAX)
				canais[0] <- valor
			}
			canais[0] <- MAX + 1

			for i := 0; i < z; i++ {
				<-result
			}
			<-canais[z]

			duracao := time.Since(inicio).Seconds()

			fmt.Printf("Execução paralela com %d processadores e %d números: %.6f segundos\n", p, z, duracao)

		}
	}
}

func cellSorter(i int, in chan int, out chan int, result chan int, max int) {
	var myVal int
	var indefinido = true
	for {
		n := <-in
		if n == max+1 {
			result <- myVal
			out <- n
			break
		}
		if indefinido {
			myVal = n
			indefinido = false
		} else if n >= myVal {
			out <- n
		} else {
			out <- myVal
			myVal = n
		}
	}
}
