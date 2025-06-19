package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
  "runtime"
)

func main() {
	tamanhos := []int{10000, 50000, 100000, 500000, 1000000}
	processadores := []int{1, 2, 4} 

	for _, P := range processadores {
		runtime.GOMAXPROCS(P)
		fmt.Printf("\n Teste com %d processadores \n", P)

		for _, tamanho := range tamanhos {
			slice := generateSlice(tamanho)

			var wg sync.WaitGroup
			wg.Add(1)
			start := time.Now()
			go func() {
				_ = mergeSortPar(slice, &wg)
			}()
			wg.Wait()
			duracao := time.Since(start).Seconds()

			fmt.Printf("Tamanho: %d - Tempo: %.6f segundos - Processadores: %d\n", tamanho, duracao, P)
		}
	}
}

// Gera um slice com valores aleatórios
func generateSlice(size int) []int {
	slice := make([]int, size)
	rand.Seed(time.Now().UnixNano())
	for i := range slice {
		slice[i] = rand.Intn(999999)
	}
	return slice
}

func mergeSortPar(arr []int, wg *sync.WaitGroup) []int {
	defer wg.Done()

	if len(arr) <= 1 {
		return arr
	}

	mid := len(arr) / 2

	var left, right []int
	var wgInner sync.WaitGroup
	wgInner.Add(2)

	go func() {
		left = mergeSortPar(arr[:mid], &wgInner)
	}()
	go func() {
		right = mergeSortPar(arr[mid:], &wgInner)
	}()

	wgInner.Wait()
	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}
