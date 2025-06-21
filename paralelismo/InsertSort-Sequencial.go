package main

import (
	"fmt"
	"math/rand"
	"time"
)

const N = 200
const MAX = 999

func main() {
	var v [N]int
	fmt.Println("  ------ sequencial -------")
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < N; i++ {
		valor := rand.Intn(MAX) - rand.Intn(MAX)

		j := 0
		for j = 0; j < i; j++ {
			if v[j] >= valor {
				break
			}
		}
		for k := i; k > j; k-- {
			v[k] = v[k-1]
		}
		v[j] = valor
	}

	fmt.Println("Primeiros 20 valores ordenados:")
	for i := 0; i < 20 && i < N; i++ {
		fmt.Printf("%d ", v[i])
	}
	fmt.Println()

	sorted := true
	for i := 1; i < N; i++ {
		if v[i] < v[i-1] {
			sorted = false
			break
		}
	}
	fmt.Printf("Array está ordenado: %t\n", sorted)
}
