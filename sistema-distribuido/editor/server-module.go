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
	Req              chan AppServerRequest   // canal para receber pedidos da aplicacao
	Ind              chan AppServerResponse  // canal para entregar respostas para a aplicacao
	processes        []string                // endereco de todos os processos
	id               int                     // identificador do processo - é o indice no array de enderecos acima
	dbg              bool                    // utilizado para logs
	text             []string                // texto sendo editado
	criticalSections []int                   // cada indice do array representa a secao critica da linha correspondente, preenchida com o id de um processo que esta acessando ela ou -1

	Pp2plink *PP2PLink.PP2PLink              // acesso a comunicacao enviar por PP2PLinq.Req  e receber por PP2PLinq.Ind
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

func initializeCriticalSections(numLines int) []int {
	csArray := make([]int, numLines)
	for i := range csArray {
		csArray[i] = -1
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
					if strings.Contains(msgOutro.Message, "readReq") {
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
// ------- UPON readReq
// ------- UPON entryReq
// ------- UPON exitReq
// ------- UPON writeReq
// ------------------------------------------------------------------------------------

func (module *Editor_Server_Module) handleUponDeliverRead(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	clientId, err := strconv.Atoi(strings.TrimPrefix(msgOutro.Message, "readReq,"))
	if err == nil {
		messageToSend := "respOk," + strings.Join(module.text, "\n")
		module.sendToLink(module.processes[clientId], messageToSend, strconv.Itoa(module.id))
	}
}

func (module *Editor_Server_Module) handleUponDeliverEntry(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	messageElements := strings.Split(msgOutro.Message, ",")
	clientId, clientIdError := strconv.Atoi(messageElements[1])
	lineIndex, lineError := strconv.Atoi(messageElements[2])

	//Se nao conseguir obter id do cliente, nao consegue enviar entryError para ele
	if clientIdError != nil {
		module.outDbg("ERRO ao obter endereco do processo que enviou evento entryReq: " + clientIdError.Error())
		return
	}

	// Se nao conseguir obter index valido da linha, enviar entryError para o cliente
	if lineError != nil {
		module.outDbg("ERRO ao obter linha para acessar em evento entryReq: " + lineError.Error())
		module.sendToLink(module.processes[clientId], "entryError,Unexpected error", strconv.Itoa(module.id))
		return
	}

	// Se nenhum outro processo estiver editando aquela linha, pode acessar a secao critica
	if module.criticalSections[lineIndex] == -1 {
		module.criticalSections[lineIndex] = clientId
		messageToSend := "entryOk"
		module.sendToLink(module.processes[clientId], messageToSend, strconv.Itoa(module.id))
	} else { // Senao, nao pode acessar a secao critica para editar a linha
		messageToSend := "entryError,User " + module.processes[module.criticalSections[lineIndex]] + " is editing this line"
		module.sendToLink(module.processes[clientId], messageToSend, strconv.Itoa(module.id))
	}
}

func (module *Editor_Server_Module) handleUponDeliverExit(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	messageElements := strings.Split(msgOutro.Message, ",")
	clientId, clientIdError := strconv.Atoi(messageElements[1])
	lineIndex, lineError := strconv.Atoi(messageElements[2])

	if clientIdError != nil {
		module.outDbg("ERRO ao obter endereco do processo que enviou evento entryReq: " + clientIdError.Error())
		return
	}
	if lineError != nil {
		module.outDbg("ERRO ao obter linha para acessar em evento entryReq: " + lineError.Error())
		return
	}

	// Se processo estava editando aquela linha, libera o acesso a linha
	if module.criticalSections[lineIndex] == clientId {
		module.criticalSections[lineIndex] = -1
		module.sendToLink(module.processes[clientId], "exitOk", strconv.Itoa(module.id))
	}
}

func (module *Editor_Server_Module) handleUponDeliverWrite(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	messageElements := strings.Split(msgOutro.Message, ",")
	clientId, clientIdError := strconv.Atoi(messageElements[1])
	lineToUpdate, lineError := strconv.Atoi(messageElements[2])
	lineUpdatedValue := messageElements[3:]

	if clientIdError != nil {
		module.outDbg("ERRO ao obter endereco do processo que enviou evento writeReq: " + clientIdError.Error())
	} else if lineError != nil {
		module.outDbg("ERRO ao obter linha para editar em evento writeReq: " + lineError.Error())
	} else if module.criticalSections[lineToUpdate] != clientId { // Se processo nao tem acesso a secao critica da linha
		module.outDbg("ERRO ao processar evento writeReq recebido: processo " + module.processes[clientId] + " nao tem acesso a secao critica")
	} else {
		// Se conseguiu obter os dados do evento sem erros e se este processo tem acesso a secao critica da linha a editar
		module.criticalSections[lineToUpdate] = -1 // sai da secao critica
		module.text[lineToUpdate] = strings.Join(lineUpdatedValue, "")
		module.broadcastTextToAllProcesses()
	}
}

// ------------------------------------------------------------------------------------
// ------- funcoes de ajuda
// ------------------------------------------------------------------------------------

func (module *Editor_Server_Module) broadcastTextToAllProcesses() {
	messageToSend := "respOk," + strings.Join(module.text, "\n")
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