/*
Aplicacao para disparar execução de servidor central em sistema distribuído de editor de texto
Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito
*/

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sistema-distribuido/editor"
)

func main() {

	if len(os.Args) < 5 {
		fmt.Println("Please specify server address:port!")
		fmt.Println("go run runServer.go <number of lines> <number of columns> <file name> <address:port>")
		fmt.Println("Example: go run runServer.go 20 20 test.txt 127.0.0.1:5000")
		return
	}

	numberOfLinesText, _ := strconv.Atoi(os.Args[1])
	numberOfColumnsText, _ := strconv.Atoi(os.Args[2])
	fileName := os.Args[3]
	address := os.Args[4]

	// Instancia modulo do servidor central
	var editor *editor.Editor_Server_Module = editor.NewServer(address, true, numberOfLinesText, numberOfColumnsText)

	for {
		response := <- editor.Ind

		// Armazena conteudo do texto sendo editado em arquivo
		bytesToWrite := []byte(strings.Join(response.Text, "\n"))
		err := os.WriteFile(fileName, bytesToWrite, 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			log.Fatal()
		}
	}
}
