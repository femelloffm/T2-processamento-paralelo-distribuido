//Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	tamanhos := []int{10000, 50000, 100000, 500000, 1000000}

	for _, tamanho := range tamanhos {
		slice := generateSlice(tamanho)

		start := time.Now()
		_ = mergeSortSeq(slice)
		duracao := time.Since(start).Seconds()

		fmt.Printf("Tamanho: %d - Tempo: %.6f segundos\n", tamanho, duracao)
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
