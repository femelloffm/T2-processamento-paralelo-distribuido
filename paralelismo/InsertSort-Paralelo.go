package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const MAX = 999

func main() {
	var N int
	fmt.Print("Digite o valor de N (quantidade de números): ")
	fmt.Scanf("%d", &N)

	if N <= 0 {
		fmt.Println("N deve ser maior que 0")
		return
	}

	fmt.Println("------ Pipeline Sort Dinâmico -------")
	fmt.Printf("N = %d\n", N)

	// Testa com 1, 2 e 4 núcleos
	nucleos := []int{1, 2, 4}

	for _, numNucleos := range nucleos {
		fmt.Printf("\n=== Executando com %d núcleo(s) ===\n", numNucleos)
		executarPipelineSort(N, numNucleos)
	}
}

func executarPipelineSort(N int, numNucleos int) {

	runtime.GOMAXPROCS(numNucleos)

	maxGoroutines := numNucleos * 2
	numProcessos := N
	if numProcessos > maxGoroutines && maxGoroutines >= N {
		numProcessos = maxGoroutines
	}

	fmt.Printf("Usando %d goroutines para %d números\n", numProcessos, N)

	start := time.Now()

	canais := make([]chan int, numProcessos+1)
	for i := 0; i <= numProcessos; i++ {
		canais[i] = make(chan int, N+1)
	}

	result := make(chan int, N)

	var wg sync.WaitGroup
	for i := 0; i < numProcessos; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cellSorter(id, canais[id], canais[id+1], result, MAX)
		}(i)
	}

	valores := make([]int, N)
	go func() {
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < N; i++ {
			valor := rand.Intn(MAX) - rand.Intn(MAX)
			valores[i] = valor
			canais[0] <- valor
		}

		canais[0] <- MAX + 1
	}()

	resultados := make([]int, N)
	for i := 0; i < N; i++ {
		v := <-result
		resultados[i] = v
	}

	<-canais[numProcessos]

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Print("Entrada (10 primeiros): ")
	limite := 10
	if N < 10 {
		limite = N
	}
	for i := 0; i < limite; i++ {
		fmt.Printf("%d ", valores[i])
	}
	fmt.Println()

	fmt.Print("Saída (10 primeiros):   ")
	for i := 0; i < limite; i++ {
		fmt.Printf("%d ", resultados[i])
	}
	fmt.Println()

	fmt.Printf("Tempo de execução: %v | Núcleos: %d\n", elapsed, numNucleos)
}

func cellSorter(id int, in chan int, out chan int, result chan int, max int) {
	var myVal int
	var undef bool = true

	for {
		n := <-in

		if n == max+1 {
			if !undef {
				result <- myVal
			}
			out <- n
			break
		}

		if undef {
			myVal = n
			undef = false
		} else if n >= myVal {
			out <- n
		} else {
			out <- myVal
			myVal = n
		}
	}
}
