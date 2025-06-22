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
	Text []string
}

type Editor_Server_Module struct {
	Req              chan AppServerRequest   // canal para receber pedidos da aplicacao
	Ind              chan AppServerResponse  // canal para entregar respostas para a aplicacao
	address          string                  // endereco do servidor central
	processes        []string                // endereco de todos os processos cliente conectados ao servidores
	dbg              bool                    // utilizado para logs
	text             []string                // texto sendo editado
	criticalSections []string                // cada indice do array representa a secao critica da linha correspondente, preenchida com o endereco do processo que esta acessando ela ou ""

	Pp2plink *PP2PLink.PP2PLink              // acesso a comunicacao enviar por PP2PLinq.Req  e receber por PP2PLinq.Ind
}

// ------------------------------------------------------------------------------------
// ------- inicializacao
// ------------------------------------------------------------------------------------

func NewServer(_serverAddress string, _dbg bool, _numLines int, _numColumns int) *Editor_Server_Module {

	p2p := PP2PLink.NewPP2PLink(_serverAddress, _dbg)

	server := &Editor_Server_Module{
		Req: make(chan AppServerRequest, 1),
		Ind: make(chan AppServerResponse, 1),

		address: _serverAddress,
		processes: make([]string, 0),
		dbg: _dbg,
		text: initializeText(_numLines, _numColumns),
		criticalSections: initializeCriticalSections(_numLines),

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

func initializeCriticalSections(numLines int) []string {
	csArray := make([]string, numLines)
	for i := range csArray {
		csArray[i] = ""
	}
	return csArray
}

// ------------------------------------------------------------------------------------
// ------- nucleo do funcionamento
// ------------------------------------------------------------------------------------

func (module *Editor_Server_Module) Start() {

	go func() {
		for {
			select {
				case msgOutro := <-module.Pp2plink.Ind: // vindo de outro processo por meio do modulo link perfeito
					module.outDbg("          <<<---- pede??  " + msgOutro.Message)
					if strings.Contains(msgOutro.Message, "disconnectReq") {
						module.handleUponDeliverDisconnect(msgOutro) // cliente desconectando
					} else if strings.Contains(msgOutro.Message, "connectReq") {
						module.handleUponDeliverConnect(msgOutro) // cliente conectando
					} else if strings.Contains(msgOutro.Message, "readReq") {
						module.handleUponDeliverRead(msgOutro) // leitura
					} else if strings.Contains(msgOutro.Message, "entryReq") {
						module.handleUponDeliverEntry(msgOutro) // entrada em secao critica
					} else if strings.Contains(msgOutro.Message, "exitReq") {
						module.handleUponDeliverExit(msgOutro) // saida da secao critica
					} else if strings.Contains(msgOutro.Message, "writeReq") {
						module.handleUponDeliverWrite(msgOutro) // escrita em linha do texto
					}
			}
		}
	}()
}

// ------------------------------------------------------------------------------------
// ------- tratamento de mensagens de outros processos
// ------- UPON connectReq
// ------- UPON disconnectReq
// ------- UPON readReq
// ------- UPON entryReq
// ------- UPON exitReq
// ------- UPON writeReq
// ------------------------------------------------------------------------------------

func (module *Editor_Server_Module) handleUponDeliverConnect(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	clientAddress := strings.TrimPrefix(msgOutro.Message, "connectReq,")
	module.processes = append(module.processes, clientAddress)
	module.sendToLink(clientAddress, "connectOk", module.address)
}

func (module *Editor_Server_Module) handleUponDeliverDisconnect(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	clientAddress := strings.TrimPrefix(msgOutro.Message, "disconnectReq,")
	numberOfClients := len(module.processes)
	for i := 0; i < numberOfClients; i++ {
		if module.processes[i] == clientAddress {
			module.processes[i] = module.processes[numberOfClients-1]
			module.processes = module.processes[0:numberOfClients-1]
			break
		}
	}
	for i := 0; i < len(module.criticalSections); i++ {
		if module.criticalSections[i] == clientAddress {
			module.criticalSections[i] = ""
		}
	}
	module.sendToLink(clientAddress, "disconnectOk", module.address)
}

func (module *Editor_Server_Module) handleUponDeliverRead(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	clientAddress := strings.TrimPrefix(msgOutro.Message, "readReq,")
	messageToSend := "respOk," + strings.Join(module.text, "\n")
	module.sendToLink(clientAddress, messageToSend, module.address)
}

func (module *Editor_Server_Module) handleUponDeliverEntry(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	messageElements := strings.Split(msgOutro.Message, ",")
	clientAddress := messageElements[1]
	lineIndex, lineError := strconv.Atoi(messageElements[2])

	// Se nao conseguir obter index valido da linha, enviar entryError para o cliente
	if lineError != nil {
		module.outDbg("ERRO ao obter linha para acessar em evento entryReq: " + lineError.Error())
		module.sendToLink(clientAddress, "entryError,Unexpected error", module.address)
		return
	}

	// Se nenhum outro processo estiver editando aquela linha, pode acessar a secao critica
	if module.criticalSections[lineIndex] == "" {
		module.criticalSections[lineIndex] = clientAddress
		messageToSend := "entryOk"
		module.sendToLink(clientAddress, messageToSend, module.address)
	} else { // Senao, nao pode acessar a secao critica para editar a linha
		messageToSend := "entryError,User " + module.criticalSections[lineIndex] + " is editing this line"
		module.sendToLink(clientAddress, messageToSend, module.address)
	}
}

func (module *Editor_Server_Module) handleUponDeliverExit(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	messageElements := strings.Split(msgOutro.Message, ",")
	clientAddress := messageElements[1]
	lineIndex, lineError := strconv.Atoi(messageElements[2])

	if lineError != nil {
		module.outDbg("ERRO ao obter linha para acessar em evento entryReq: " + lineError.Error())
		return
	}

	// Se processo estava editando aquela linha, libera o acesso a linha
	if module.criticalSections[lineIndex] == clientAddress {
		module.criticalSections[lineIndex] = ""
		module.sendToLink(clientAddress, "exitOk", module.address)
	}
}

func (module *Editor_Server_Module) handleUponDeliverWrite(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	messageElements := strings.Split(msgOutro.Message, ",")
	clientAddress := messageElements[1]
	lineToUpdate, lineError := strconv.Atoi(messageElements[2])
	lineUpdatedValue := messageElements[3:]

	if lineError != nil {
		module.outDbg("ERRO ao obter linha para editar em evento writeReq: " + lineError.Error())
	} else if module.criticalSections[lineToUpdate] != clientAddress { // Se processo nao tem acesso a secao critica da linha
		module.outDbg("ERRO ao processar evento writeReq recebido: processo " + clientAddress + " nao tem acesso a secao critica")
	} else {
		// Se conseguiu obter os dados do evento sem erros e se este processo tem acesso a secao critica da linha a editar
		module.criticalSections[lineToUpdate] = "" // sai da secao critica
		module.text[lineToUpdate] = strings.Join(lineUpdatedValue, "")
		module.broadcastTextToAllProcesses()
		module.Ind <- AppServerResponse{ module.text }
	}
}

// ------------------------------------------------------------------------------------
// ------- funcoes de ajuda
// ------------------------------------------------------------------------------------

func (module *Editor_Server_Module) broadcastTextToAllProcesses() {
	messageToSend := "respOk," + strings.Join(module.text, "\n")
	for i := 1; i < len(module.processes); i++ {
		module.sendToLink(module.processes[i], messageToSend, module.address);
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