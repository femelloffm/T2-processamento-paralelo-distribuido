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
	CONNECT
	DISCONNECT
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
	CONNECT_OK
	DISCONNECT_OK
)

type AppClientResponse struct {
	Type AppClientResponseType
	Text []string
	Err  *string
}

type Editor_Client_Module struct {
	Req           chan AppClientRequest  // canal para receber pedidos da aplicacao (READ, ENTRY e WRITE)
	Ind           chan AppClientResponse // canal para entregar informacao para a aplicacao
	serverAddress string                 // endereco do servidor central
	address       string                 // endereco do processo cliente
	dbg           bool                   // utilizado para logs

	Pp2plink      *PP2PLink.PP2PLink     // acesso aa comunicacao enviar por PP2PLinq.Req  e receber por PP2PLinq.Ind
}

// ------------------------------------------------------------------------------------
// ------- inicializacao
// ------------------------------------------------------------------------------------

func NewClient(_serverAddress string, _clientAddress string, _dbg bool) *Editor_Client_Module {

	p2p := PP2PLink.NewPP2PLink(_clientAddress, _dbg)

	client := &Editor_Client_Module{
		Req: make(chan AppClientRequest, 1),
		Ind: make(chan AppClientResponse, 1),

		serverAddress: _serverAddress,
		address: _clientAddress,
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
				if appReq.Type == CONNECT {
					module.outDbg("APP quer se conectar ao servidor central")
					module.handleUponReqConnect()
				} else if appReq.Type == DISCONNECT {
					module.outDbg("APP quer se desconectar do servidor central")
					module.handleUponReqDisconnect()
				} else if appReq.Type == READ {
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
					module.handleUponDeliverEntryOk()
				} else if strings.HasPrefix(msgOutro.Message, "entryError") {
					module.handleUponDeliverEntryError(msgOutro)
				} else if strings.HasPrefix(msgOutro.Message, "disconnectOk") {
					module.handleUponDeliverDisconnectOk()
				} else if strings.HasPrefix(msgOutro.Message, "connectOk") {
					module.handleUponDeliverConnectOk()
				}
			}
		}
	}()
}

// ------------------------------------------------------------------------------------
// ------- tratamento de pedidos vindos da aplicacao
// ------- UPON connect
// ------- UPON disconnect
// ------- UPON read
// ------- UPON entry
// ------- UPON exit
// ------- UPON write
// ------------------------------------------------------------------------------------

func (module *Editor_Client_Module) handleUponReqConnect() {
	// Envia evento para servidor central contendo: endereco do processo
	messageToSend := "connectReq," + module.address
	module.sendToLink(module.serverAddress, messageToSend, module.address);
}

func (module *Editor_Client_Module) handleUponReqDisconnect() {
	// Envia evento para servidor central contendo: endereco do processo
	messageToSend := "disconnectReq," + module.address
	module.sendToLink(module.serverAddress, messageToSend, module.address);
}

func (module *Editor_Client_Module) handleUponReqRead() {
	// Envia evento para servidor central contendo: endereco do processo
	messageToSend := "readReq," + module.address
	module.sendToLink(module.serverAddress, messageToSend, module.address);
}

func (module *Editor_Client_Module) handleUponReqEntry(appReq AppClientRequest) {
	// Envia evento para servidor central contendo: endereco do processo e index da linha para acessar
	messageToSend := "entryReq," + module.address + "," + strconv.Itoa(*appReq.Cursor)
	module.sendToLink(module.serverAddress, messageToSend, module.address);
}

func (module *Editor_Client_Module) handleUponReqExit(appReq AppClientRequest) {
	// Envia evento para servidor central contendo: endereco do processo e index da linha para acessar
	messageToSend := "exitReq," + module.address + "," + strconv.Itoa(*appReq.Cursor)
	module.sendToLink(module.serverAddress, messageToSend, module.address);
}

func (module *Editor_Client_Module) handleUponReqWrite(appReq AppClientRequest) {
	// Envia evento para servidor central contendo: endereco do processo, index da linha para editar, e novo conteudo da linha
	message := "writeReq," + module.address + "," + strconv.Itoa(*appReq.Cursor) + "," + *appReq.Line
	module.sendToLink(module.serverAddress, message, module.address);
}

// ------------------------------------------------------------------------------------
// ------- tratamento de mensagens de outros processos
// ------- UPON respOk
// ------- UPON entryOk
// ------- UPON entryError
// ------- UPON connectOk
// ------- UPON disconnectOk
// ------------------------------------------------------------------------------------

func (module *Editor_Client_Module) handleUponDeliverRespOk(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	updatedText := strings.Split(strings.TrimPrefix(msgOutro.Message, "respOk,"), "\n")
	module.Ind <- AppClientResponse{ RESP, updatedText, nil }
}

func (module *Editor_Client_Module) handleUponDeliverEntryOk() {
	module.Ind <- AppClientResponse{ ENTRY_OK, nil, nil }
}

func (module *Editor_Client_Module) handleUponDeliverEntryError(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	errorMessage := strings.TrimPrefix(msgOutro.Message, "entryError,")
	module.Ind <- AppClientResponse{ ENTRY_ERROR, nil, &errorMessage }
}

func (module *Editor_Client_Module) handleUponDeliverConnectOk() {
	module.Ind <- AppClientResponse{ CONNECT_OK, nil, nil }
}

func (module *Editor_Client_Module) handleUponDeliverDisconnectOk() {
	module.Ind <- AppClientResponse{ DISCONNECT_OK, nil, nil }
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