// Aplicacao para disparar execução de servidor central em sistema distribuído de editor de texto
// Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito

package main

import (
	"fmt"
	"log"
	"os"
	"sistema-distribuido/editor"
	"strconv"
	"strings"
	"time"
)

func main() {

	if len(os.Args) < 5 {
		fmt.Println("Please specify at least one address:port!")
		fmt.Println("go run runServer.go <number of lines> <number of columns> <file name> <id> <addresses>")
		fmt.Println("go run runServer.go 20 20 0 127.0.0.1:5000  127.0.0.1:6001  127.0.0.1:7002")
		return
	}

	numberOfLinesText, _ := strconv.Atoi(os.Args[1])
	numberOfColumnsText, _ := strconv.Atoi(os.Args[2])
	fileName := os.Args[3]
	id, _ := strconv.Atoi(os.Args[4])
	addresses := os.Args[5:]

	var editor *editor.Editor_Server_Module = editor.NewServer(addresses, id, true, numberOfLinesText, numberOfColumnsText)
	time.Sleep(5 * time.Second)

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
