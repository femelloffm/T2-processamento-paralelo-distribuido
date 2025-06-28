/*
Constantes utilizadas por modulos do servidor e cliente de editor de texto distribuido
Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito
*/

package editor

const MESSAGE_SEPARATOR string = "," // separa informacoes recebidas em mensagens atraves do modulo link perfeito

const DISCONNECT_REQ string = "disconnectReq"
const CONNECT_REQ string = "connectReq"
const READ_REQ string = "readReq"
const ENTRY_REQ string = "entryReq"
const EXIT_REQ string = "exitReq"
const WRITE_REQ string = "writeReq"

const TEXT_OK_RESP string = "respOk"
const ENTRY_OK_RESP string = "entryOk"
const ENTRY_ERROR_RESP string = "entryError"
const EXIT_OK_RESP string = "exitOk"
const DISCONNECT_OK_RESP string = "disconnectOk"
const CONNECT_OK_RESP string = "connectOk"