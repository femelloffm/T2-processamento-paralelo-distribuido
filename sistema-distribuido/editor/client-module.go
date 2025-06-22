/*
Modulo do cliente no sistema distribuído de um editor de texto
Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito
*/

package editor

import (
	"fmt"
	PP2PLink "sistema-distribuido/PP2PLink"
	"strconv"
	"strings"
)

// ------------------------------------------------------------------------------------
// ------- principais tipos
// ------------------------------------------------------------------------------------

type AppClientRequestType int
const (
	READ AppClientRequestType = iota
	ENTRY
	EXIT
	WRITE
)

type AppClientRequest struct {
	Type   AppClientRequestType
	Cursor *int
	Line   *string
}

type AppClientResponseType int
const (
	RESP AppClientResponseType = iota
	ENTRY_OK
	ENTRY_ERROR
)

type AppClientResponse struct {
	Type AppClientResponseType
	Text []string
	Err  *string
}

type Editor_Client_Module struct {
	Req       chan AppClientRequest   // canal para receber pedidos da aplicacao (READ, ENTRY e WRITE)
	Ind       chan AppClientResponse  // canal para entregar informacao para a aplicacao
	processes []string                // endereco de todos, na mesma ordem
	id        int                     // identificador do processo - é o indice no array de enderecos acima
	dbg       bool                    // utilizado para logs

	Pp2plink *PP2PLink.PP2PLink       // acesso aa comunicacao enviar por PP2PLinq.Req  e receber por PP2PLinq.Ind
}

// ------------------------------------------------------------------------------------
// ------- inicializacao
// ------------------------------------------------------------------------------------

func NewClient(_addresses []string, _id int, _dbg bool) *Editor_Client_Module {

	p2p := PP2PLink.NewPP2PLink(_addresses[_id], _dbg)

	client := &Editor_Client_Module{
		Req: make(chan AppClientRequest, 1),
		Ind: make(chan AppClientResponse, 1),

		processes: _addresses,
		id:        _id,
		dbg:       _dbg,

		Pp2plink: p2p}

	client.Start()
	client.outDbg("Init text editor client!")
	return client
}

// ------------------------------------------------------------------------------------
// ------- nucleo do funcionamento
// ------------------------------------------------------------------------------------

func (module *Editor_Client_Module) Start() {

	go func() {
		for {
			select {
			case appReq := <-module.Req: // vindo da aplicação
				if appReq.Type == READ {
					module.outDbg("APP quer ler texto")
					module.handleUponReqRead()
				} else if appReq.Type == ENTRY {
					module.outDbg("APP quer ter acesso a linha do texto para editar")
					module.handleUponReqEntry(appReq)
				} else if appReq.Type == EXIT {
					module.outDbg("APP quer liberar o acesso a linha do texto")
					module.handleUponReqExit(appReq)
				} else if appReq.Type == WRITE {
					module.outDbg("APP quer editar linha de texto")
					module.handleUponReqWrite(appReq)
				}

			case msgOutro := <-module.Pp2plink.Ind: // vindo de outro processo via modulo link perfeito
				module.outDbg("         <<<---- responde! " + msgOutro.Message)
				if strings.HasPrefix(msgOutro.Message, "respOk") {
					module.handleUponDeliverRespOk(msgOutro)
				} else if strings.HasPrefix(msgOutro.Message, "entryOk") {
					module.handleUponDeliverEntryOk(msgOutro)
				} else if strings.HasPrefix(msgOutro.Message, "entryError") {
					module.handleUponDeliverEntryError(msgOutro)
				}
			}
		}
	}()
}

// ------------------------------------------------------------------------------------
// ------- tratamento de pedidos vindos da aplicacao
// ------- UPON read
// ------- UPON entry
// ------- UPON exit
// ------- UPON write
// ------------------------------------------------------------------------------------

func (module *Editor_Client_Module) handleUponReqRead() {
	// Envia evento para servidor central contendo: id do processo
	messageToSend := "readReq," + strconv.Itoa(module.id)
	module.sendToLink(module.processes[0], messageToSend, strconv.Itoa(module.id));
}

func (module *Editor_Client_Module) handleUponReqEntry(appReq AppClientRequest) {
	// Envia evento para servidor central contendo: id do processo e index da linha para acessar
	messageToSend := "entryReq," + strconv.Itoa(module.id) + "," + strconv.Itoa(*appReq.Cursor)
	module.sendToLink(module.processes[0], messageToSend, strconv.Itoa(module.id));
}

func (module *Editor_Client_Module) handleUponReqExit(appReq AppClientRequest) {
	// Envia evento para servidor central contendo: id do processo e index da linha para acessar
	messageToSend := "exitReq," + strconv.Itoa(module.id) + "," + strconv.Itoa(*appReq.Cursor)
	module.sendToLink(module.processes[0], messageToSend, strconv.Itoa(module.id));
}

func (module *Editor_Client_Module) handleUponReqWrite(appReq AppClientRequest) {
	// Envia evento para servidor central contendo: id do processo, index da linha para editar, e novo conteudo da linha
	message := "writeReq," + strconv.Itoa(module.id) + "," + strconv.Itoa(*appReq.Cursor) + "," + *appReq.Line
	module.sendToLink(module.processes[0], message, strconv.Itoa(module.id));
}

// ------------------------------------------------------------------------------------
// ------- tratamento de mensagens de outros processos
// ------- UPON respOk
// ------- UPON entryOk
// ------- UPON entryError
// ------------------------------------------------------------------------------------

func (module *Editor_Client_Module) handleUponDeliverRespOk(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	updatedText := strings.Split(strings.TrimPrefix(msgOutro.Message, "respOk,"), "\n")
	module.Ind <- AppClientResponse{ RESP, updatedText, nil }
}

func (module *Editor_Client_Module) handleUponDeliverEntryOk(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	module.Ind <- AppClientResponse{ ENTRY_OK, nil, nil }
}

func (module *Editor_Client_Module) handleUponDeliverEntryError(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	errorMessage := strings.TrimPrefix(msgOutro.Message, "entryError,")
	module.Ind <- AppClientResponse{ ENTRY_ERROR, nil, &errorMessage }
}

// ------------------------------------------------------------------------------------
// ------- funcoes de ajuda
// ------------------------------------------------------------------------------------

func (module *Editor_Client_Module) sendToLink(address string, content string, space string) {
	module.outDbg(space + " ---->>>>   to: " + address + "     msg: " + content)
	module.Pp2plink.Req <- PP2PLink.PP2PLink_Req_Message{
		To:      address,
		Message: content}
}

func (module *Editor_Client_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . . . . [ CLIENT : " + s + " ]")
	}
}