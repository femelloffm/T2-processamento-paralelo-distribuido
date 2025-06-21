// Aplicacao para disparar execução de servidor central em sistema distribuído de editor de texto
// Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito

/*
  Cria um processo servidor 
  para cada processo: seu id único e a mesma lista de processos.
  	o endereco de cada processo é o dado na lista, na posicao do seu id.
*/

package main

import (
	"sistema-distribuido/editor"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	if len(os.Args) < 5 {
		fmt.Println("Please specify at least one address:port!")
		fmt.Println("go run runServer.go <number of lines> <number of columns> <id> <addresses>")
		fmt.Println("go run runServer.go 20 20 0 127.0.0.1:5000  127.0.0.1:6001  127.0.0.1:7002")
		return
	}

	numberOfLinesText, _ := strconv.Atoi(os.Args[1])
	numberOfColumnsText, _ := strconv.Atoi(os.Args[2])
	id, _ := strconv.Atoi(os.Args[3])
	addresses := os.Args[4:]

	var dmx *editor.Editor_Server_Module = editor.NewServer(addresses, id, true, numberOfLinesText, numberOfColumnsText)
	fmt.Println(dmx)

	time.Sleep(5 * time.Second)

	for {
		<- dmx.Req
	}
}
