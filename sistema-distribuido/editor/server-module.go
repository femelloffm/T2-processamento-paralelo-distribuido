/*
Modulo do servidor central no sistema distribuído de um editor de texto
Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito
*/

package editor

import (
	PP2PLink "sistema-distribuido/PP2PLink"
	"fmt"
	"strings"
	"strconv"
)

// ------------------------------------------------------------------------------------
// ------- principais tipos
// ------------------------------------------------------------------------------------

type AppServerRequest struct {
}

type AppServerResponse struct {
}

type Editor_Server_Module struct {
	Req       chan AppServerRequest   // canal para receber pedidos da aplicacao (REQ e EXIT)
	Ind       chan AppServerResponse  // canal para informar aplicacao que pode acessar
	processes []string          // endereco de todos os processos
	id        int               // identificador do processo - é o indice no array de enderecos acima
	dbg       bool
	text      []string          // texto sendo editado

	Pp2plink *PP2PLink.PP2PLink // acesso a comunicacao enviar por PP2PLinq.Req  e receber por PP2PLinq.Ind
}

// ------------------------------------------------------------------------------------
// ------- inicializacao
// ------------------------------------------------------------------------------------

func NewServer(_addresses []string, _id int, _dbg bool, _numLines int, _numColumns int) *Editor_Server_Module {

	p2p := PP2PLink.NewPP2PLink(_addresses[_id], _dbg)

	server := &Editor_Server_Module{
		Req: make(chan AppServerRequest, 1),
		Ind: make(chan AppServerResponse, 1),

		processes: _addresses,
		id:        _id,
		dbg:       _dbg,
		text:      initializeText(_numLines, _numColumns),

		Pp2plink: p2p}
	
	server.Start()
	server.outDbg("Init text editor server!")
	return server
}

func initializeText(numLines int, numColumns int) []string {
	text := make([]string, numLines)
	for i := range numLines {
		text[i] = strings.Repeat("0", numColumns)
	}
	return text
}

// ------------------------------------------------------------------------------------
// ------- nucleo do funcionamento
// ------------------------------------------------------------------------------------

func (module *Editor_Server_Module) Start() {

	go func() {
		for {
			select {
			case msgOutro := <-module.Pp2plink.Ind: // vindo de outro processo
				fmt.Println("SERVER recebe da rede: ", msgOutro)
				if strings.Contains(msgOutro.Message, "READ") {
					module.outDbg("         <<<---- responde! " + msgOutro.Message)
					module.handleUponDeliverRead(msgOutro) // ENTRADA DO ALGORITMO

				} else if strings.Contains(msgOutro.Message, "WRITE") {
					module.outDbg("          <<<---- pede??  " + msgOutro.Message)
					module.handleUponDeliverWrite(msgOutro) // ENTRADA DO ALGORITMO

				}
			}
		}
	}()
}

// ------------------------------------------------------------------------------------
// ------- tratamento de mensagens de outros processos
// ------- UPON read
// ------- UPON write
// ------------------------------------------------------------------------------------

func (module *Editor_Server_Module) handleUponDeliverRead(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	clientId, error := strconv.Atoi(strings.TrimPrefix(msgOutro.Message, "READ,"))
	if error == nil {
		messageToSend := "UPDATE," + strings.Join(module.text, ",")
		module.sendToLink(module.processes[clientId], messageToSend, strconv.Itoa(module.id));
	}
}

func (module *Editor_Server_Module) handleUponDeliverWrite(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	messageElements := strings.Split(msgOutro.Message, ",")
	lineToUpdate, error := strconv.Atoi(messageElements[1])
	lineUpdatedValue := messageElements[2]
	if error == nil {
		module.text[lineToUpdate] = lineUpdatedValue
		module.broadcastTextToAllProcesses()
	} else {
		fmt.Println("ERRO ao processar mensagem WRITE recebida", error)
	}
}

// ------------------------------------------------------------------------------------
// ------- funcoes de ajuda
// ------------------------------------------------------------------------------------

func (module *Editor_Server_Module) broadcastTextToAllProcesses() {
	messageToSend := "UPDATE," + strings.Join(module.text, ",")
	for i := 1; i < len(module.processes); i++ {
		module.sendToLink(module.processes[i], messageToSend, strconv.Itoa(module.id));
	}
}

func (module *Editor_Server_Module) sendToLink(address string, content string, space string) {
	module.outDbg(space + " ---->>>>   to: " + address + "     msg: " + content)
	module.Pp2plink.Req <- PP2PLink.PP2PLink_Req_Message{
		To:      address,
		Message: content}
}

func (module *Editor_Server_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . . . . [ SERVER : " + s + " ]")
	}
}