//Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito

package main

import (
	"fmt"
	"math/rand"
	"time"
)

var u = [3]int{300, 600, 900}

const MAX = 999

func main() {
	for _, u := range u {
		fmt.Printf("------ Execução sequencial com %d números ------\n", u)

		v := make([]int, u)
		rand.Seed(time.Now().UnixNano())

		inicio := time.Now()

		for i := 0; i < u; i++ {
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

		duracao := time.Since(inicio).Seconds()

		sorted := true
		for i := 1; i < u; i++ {
			if v[i] < v[i-1] {
				sorted = false
				break
			}
		}

		fmt.Printf("Execução sequencial com %d números: %.6f segundos - Ordenado: %t\n", u, duracao, sorted)
	}
}
