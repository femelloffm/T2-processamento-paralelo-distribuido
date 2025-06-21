/*
Modulo do cliente no sistema distribuído de um editor de texto
Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito
*/

package editor

import (
	PP2PLink "sistema-distribuido/PP2PLink"
	"fmt"
	"strconv"
	"strings"
)

// ------------------------------------------------------------------------------------
// ------- principais tipos
// ------------------------------------------------------------------------------------

type AppClientRequestType int
const (
	READ AppClientRequestType = iota
	WRITE
)

type AppClientRequest struct {
	RequestType AppClientRequestType
	Cursor *int
	Line *string
}

type AppClientResponse struct {
	Text string
}

type Editor_Client_Module struct {
	Req       chan AppClientRequest   // canal para receber pedidos da aplicacao (REQ e EXIT)
	Ind       chan AppClientResponse  // canal para informar aplicacao que pode acessar
	processes []string          // endereco de todos, na mesma ordem
	id        int               // identificador do processo - é o indice no array de enderecos acima
	dbg       bool

	Pp2plink *PP2PLink.PP2PLink // acesso aa comunicacao enviar por PP2PLinq.Req  e receber por PP2PLinq.Ind
}

// ------------------------------------------------------------------------------------
// ------- inicializacao
// ------------------------------------------------------------------------------------

func NewClient(_addresses []string, _id int, _dbg bool) *Editor_Client_Module {

	p2p := PP2PLink.NewPP2PLink(_addresses[_id], _dbg)

	dmx := &Editor_Client_Module{
		Req: make(chan AppClientRequest, 1),
		Ind: make(chan AppClientResponse, 1),

		processes: _addresses,
		id:        _id,
		dbg:       _dbg,

		Pp2plink: p2p}

	dmx.Start()
	dmx.outDbg("Init text editor client!")
	return dmx
}

// ------------------------------------------------------------------------------------
// ------- nucleo do funcionamento
// ------------------------------------------------------------------------------------

func (module *Editor_Client_Module) Start() {

	go func() {
		for {
			select {
			case appReq := <-module.Req: // vindo da  aplicação
				fmt.Println("CLIENT recebe da app: ", appReq)
				if appReq.RequestType == READ {
					module.outDbg("APP quer ler texto")
					module.handleUponReqRead() // ENTRADA DO ALGORITMO

				} else if appReq.RequestType == WRITE {
					module.outDbg("APP quer editar linha de texto")
					module.handleUponReqWrite(appReq) // ENTRADA DO ALGORITMO
				}

			case msgOutro := <-module.Pp2plink.Ind: // vindo de outro processo
				fmt.Println("CLIENT recebe da rede: ", msgOutro)
				if strings.HasPrefix(msgOutro.Message, "UPDATE") {
					module.outDbg("         <<<---- responde! " + msgOutro.Message)
					module.handleUponDeliverUpdate(msgOutro) // ENTRADA DO ALGORITMO
				}
			}
		}
	}()
}

// ------------------------------------------------------------------------------------
// ------- tratamento de pedidos vindos da aplicacao
// ------- UPON read
// ------- UPON write
// ------------------------------------------------------------------------------------

func (module *Editor_Client_Module) handleUponReqRead() {
	messageToSend := "READ," + strconv.Itoa(module.id)
	module.sendToLink(module.processes[0], messageToSend, strconv.Itoa(module.id));
}

func (module *Editor_Client_Module) handleUponReqWrite(appReq AppClientRequest) {
	message := "WRITE," + strconv.Itoa(*appReq.Cursor) + "," + *appReq.Line
	module.sendToLink(module.processes[0], message, strconv.Itoa(module.id));
}

// ------------------------------------------------------------------------------------
// ------- tratamento de mensagens de outros processos
// ------- UPON update
// ------------------------------------------------------------------------------------

func (module *Editor_Client_Module) handleUponDeliverUpdate(msgOutro PP2PLink.PP2PLink_Ind_Message) {
	updatedText := strings.TrimPrefix(msgOutro.Message, "UPDATE,")
	module.Ind <- AppClientResponse{ updatedText }
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