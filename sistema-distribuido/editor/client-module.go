/*
Modulo do cliente no sistema distribuído de um editor de texto
Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito
*/

package editor

import (
	"fmt"
	"strconv"
	"strings"
	PP2PLink "sistema-distribuido/PP2PLink"
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
	Req           chan AppClientRequest  // canal para receber pedidos da aplicacao (READ, ENTRY, WRITE, EXIT, CONNECT e DISCONNECT)
	Ind           chan AppClientResponse // canal para entregar informacao para a aplicacao
	serverAddress string                 // endereco do servidor central
	address       string                 // endereco do processo cliente
	dbg           bool                   // indica se modulo deve ser executado em modo de debug (exibindo logs)
	Pp2plink      *PP2PLink.PP2PLink     // acesso a comunicacao - enviar por PP2PLinq.Req e receber por PP2PLinq.Ind
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
		dbg: _dbg,
		Pp2plink: p2p}

	client.Start()
	client.outDbg("Init text editor client module!")
	return client
}

// ------------------------------------------------------------------------------------
// ------- nucleo do funcionamento
// ------------------------------------------------------------------------------------

func (module *Editor_Client_Module) Start() {
	go func() {
		for {
			select {
				// Vindo da aplicacao
				case appReq := <-module.Req:
					module.handleAppRequest(appReq)
				// Vindo de outro processo via modulo link perfeito
				case msgOutro := <-module.Pp2plink.Ind:
					module.handlePerfectLinkMessage(msgOutro)
			}
		}
	}()
}

func (module *Editor_Client_Module) handleAppRequest(appReq AppClientRequest) {
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
}

func (module *Editor_Client_Module) handlePerfectLinkMessage(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	module.outDbg("         <<<---- responde! " + msgOutro.Message)
	if strings.HasPrefix(msgOutro.Message, TEXT_OK_RESP) {
		module.handleUponDeliverRespOk(msgOutro)
	} else if strings.HasPrefix(msgOutro.Message, ENTRY_OK_RESP) {
		module.handleUponDeliverEntryOk()
	} else if strings.HasPrefix(msgOutro.Message, ENTRY_ERROR_RESP) {
		module.handleUponDeliverEntryError(msgOutro)
	} else if strings.HasPrefix(msgOutro.Message, DISCONNECT_OK_RESP) {
		module.handleUponDeliverDisconnectOk()
	} else if strings.HasPrefix(msgOutro.Message, CONNECT_OK_RESP) {
		module.handleUponDeliverConnectOk()
	}
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
	messageToSend := CONNECT_REQ + MESSAGE_SEPARATOR + module.address
	module.sendToLink(module.serverAddress, messageToSend, module.address);
}

func (module *Editor_Client_Module) handleUponReqDisconnect() {
	// Envia evento para servidor central contendo: endereco do processo
	messageToSend := DISCONNECT_REQ + MESSAGE_SEPARATOR + module.address
	module.sendToLink(module.serverAddress, messageToSend, module.address);
}

func (module *Editor_Client_Module) handleUponReqRead() {
	// Envia evento para servidor central contendo: endereco do processo
	messageToSend := READ_REQ + MESSAGE_SEPARATOR + module.address
	module.sendToLink(module.serverAddress, messageToSend, module.address);
}

func (module *Editor_Client_Module) handleUponReqEntry(appReq AppClientRequest) {
	// Envia evento para servidor central contendo: endereco do processo e indice da linha para acessar
	messageToSend := ENTRY_REQ + MESSAGE_SEPARATOR + module.address + MESSAGE_SEPARATOR + strconv.Itoa(*appReq.Cursor)
	module.sendToLink(module.serverAddress, messageToSend, module.address);
}

func (module *Editor_Client_Module) handleUponReqExit(appReq AppClientRequest) {
	// Envia evento para servidor central contendo: endereco do processo e indice da linha para acessar
	messageToSend := EXIT_REQ + MESSAGE_SEPARATOR + module.address + MESSAGE_SEPARATOR + strconv.Itoa(*appReq.Cursor)
	module.sendToLink(module.serverAddress, messageToSend, module.address);
}

func (module *Editor_Client_Module) handleUponReqWrite(appReq AppClientRequest) {
	// Envia evento para servidor central contendo: endereco do processo, indice da linha para editar, e novo conteudo da linha
	message := WRITE_REQ + MESSAGE_SEPARATOR + module.address + MESSAGE_SEPARATOR + strconv.Itoa(*appReq.Cursor) + MESSAGE_SEPARATOR + *appReq.Line
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
	// Entrega texto atualizado para aplicacao cliente
	updatedText := strings.Split(strings.TrimPrefix(msgOutro.Message, TEXT_OK_RESP + MESSAGE_SEPARATOR), "\n")
	module.Ind <- AppClientResponse{ RESP, updatedText, nil }
}

func (module *Editor_Client_Module) handleUponDeliverEntryOk() {
	module.Ind <- AppClientResponse{ ENTRY_OK, nil, nil }
}

func (module *Editor_Client_Module) handleUponDeliverEntryError(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	// Entrega erro para aplicacao cliente
	errorMessage := strings.TrimPrefix(msgOutro.Message, ENTRY_ERROR_RESP + MESSAGE_SEPARATOR)
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

// Envia mensagem para outro processo atraves do modulo link perfeito
func (module *Editor_Client_Module) sendToLink(address string, content string, space string) {
	module.outDbg(space + " ---->>>>   to: " + address + "     msg: " + content)
	module.Pp2plink.Req <- PP2PLink.PP2PLink_Req_Message{
		To:      address,
		Message: content}
}

// Realiza print da string recebida por parametro, para debug
func (module *Editor_Client_Module) outDbg(s string) {
	if module.dbg {
		fmt.Println(". . . . . . . . . . . . [ CLIENT : " + s + " ]")
	}
}