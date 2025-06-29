//Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
 	"runtime"
)

func main() {
	processadores := []int{1, 2, 4} 
	tamanhos := []int{10000, 50000, 100000, 500000, 1000000}
	granularidades := []int{1, 500, 1000, 5000}

	for _, tamanho := range tamanhos {
		baseSlice := generateSlice(tamanho)

		for _, P := range processadores {
			runtime.GOMAXPROCS(P)

			for _, G := range granularidades {
				slicePar := make([]int, len(baseSlice))
				copy(slicePar, baseSlice)

				var wg sync.WaitGroup
				wg.Add(1)
				startPar := time.Now()
				go func() {
					_ = mergeSortPar(slicePar, &wg, G)
				}()
				wg.Wait()
				tempoPar := time.Since(startPar).Seconds()

				fmt.Printf("Tamanho: %d - P: %d - G: %d - Tempo Paralelo: %.6f s\n",
					tamanho, P, G, tempoPar)
			}
		}
	}
}

func generateSlice(size int) []int {
	slice := make([]int, size)
	rand.Seed(time.Now().UnixNano())
	for i := range slice {
		slice[i] = rand.Intn(999999)
	}
	return slice
}

func mergeSortPar(arr []int, wg *sync.WaitGroup, G int) []int {
	defer wg.Done()

	if len(arr) <= 1 {
		return arr
	}

	if len(arr) <= G {
		return mergeSortSeq(arr)
	}

	mid := len(arr) / 2
	var left, right []int
	var wgInner sync.WaitGroup
	wgInner.Add(2)

	go func() {
		left = mergeSortPar(arr[:mid], &wgInner, G)
	}()
	go func() {
		right = mergeSortPar(arr[mid:], &wgInner, G)
	}()

	wgInner.Wait()
	return merge(left, right)
}

func mergeSortSeq(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := mergeSortSeq(arr[:mid])
	right := mergeSortSeq(arr[mid:])
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
