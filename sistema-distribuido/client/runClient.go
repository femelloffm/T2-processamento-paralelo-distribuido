/*
Aplicacao para disparar execução de cliente em sistema distribuído de editor de texto
Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito
*/

package main

import (
	"fmt"
	"log"
	"os"
	"sistema-distribuido/editor"
	"github.com/gdamore/tcell/v2"
)

const CURSOR rune = '*' // identifica em qual linha o cursor do usuario está

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please specify server and client address:port!")
		fmt.Println("go run runClient.go <server address:port> <client address:port>")
		fmt.Println("Example: go run runClient.go 127.0.0.1:5000  127.0.0.1:6001")
		return
	}

	serverAddress := os.Args[1]
	clientAddress := os.Args[2]

	// Instancia modulo do cliente
	var client *editor.Editor_Client_Module = editor.NewClient(serverAddress, clientAddress, false)

	screen, err := tcell.NewScreen()
	handleError(err)
	defer screen.Fini()
	err = screen.Init()
	handleError(err)

	// Conecta com servidor central e lê texto
	client.Req <- editor.AppClientRequest{ Type: editor.CONNECT, Cursor: nil, Line: nil }
	<- client.Ind
	client.Req <- editor.AppClientRequest{ Type: editor.READ, Cursor: nil, Line: nil }
	text := (<- client.Ind).Text
	maxLines := len(text)
	
	currentLine := 0
	running := true

	for running {
		select {
			// Recebe evento de modulo do cliente
			case editorModuleResponse := <- client.Ind:
				switch editorModuleResponse.Type {
					case editor.DISCONNECT_OK:
						running = false
					case editor.ENTRY_OK:
						screen.Suspend()
						text = editLine(client, text, currentLine)
						screen.Resume()
					case editor.ENTRY_ERROR:
						screen.Suspend()
						showErrorScreen(*editorModuleResponse.Err)
						screen.Resume()
					case editor.RESP:
						text = editorModuleResponse.Text
				}
			default:
				// Exibe linhas do texto
				drawText(screen, text, currentLine)
				if (screen.HasPendingEvent()) {
					ev := screen.PollEvent()
					switch ev := ev.(type) {
						case *tcell.EventKey:
							if ev.Key() == tcell.KeyEscape {
								client.Req <- editor.AppClientRequest{ Type: editor.DISCONNECT, Cursor: nil, Line: nil }
							} else if ev.Key() == tcell.KeyUp && currentLine > 0 {
								currentLine--
							} else if ev.Key() == tcell.KeyDown && currentLine < (maxLines - 1) {
								currentLine++
							} else if ev.Key() == tcell.KeyEnter {
								client.Req <- editor.AppClientRequest{ Type: editor.ENTRY, Cursor: &currentLine, Line: nil }
							}
					}
				}
		}
	}
}

// Metodo para edicao de uma das linhas do texto atraves de uma tela no terminal
func editLine(client *editor.Editor_Client_Module, text []string, line int) []string {
	screen, err := tcell.NewScreen()
	handleError(err)
	defer screen.Fini()
	err = screen.Init()
	handleError(err)

	lineContent := text[line]
	cursor := 0

	for {
		drawTextLine(screen, lineContent, cursor)

		ev := screen.PollEvent()
		switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape { // Sair de edicao de linha
					client.Req <- editor.AppClientRequest{ Type: editor.EXIT, Cursor: &line, Line: nil }
					return text
				} else if ev.Key() == tcell.KeyLeft && cursor > 0 { // Mover cursor para esquerda
					cursor--
				} else if ev.Key() == tcell.KeyRight && cursor < len(lineContent) { // Mover cursor para direita
					cursor++
				} else if ev.Key() == tcell.KeyBackspace && cursor > 0 && len(lineContent) > 0 { // Deletar caracter da linha
					if (cursor == len(lineContent)) {
						lineContent = lineContent[0:cursor-1]
					} else if (cursor == 1) {
						lineContent = lineContent[cursor:]
					} else {
						lineContent = lineContent[0:cursor-1] + lineContent[cursor:]
					}
					cursor--
				} else if ev.Key() == tcell.KeyEnter { // Salvar alteracoes e sair da edicao da linha
					sendWriteRequest(client, line, lineContent)
					text[line] = lineContent
					return text
				} else if ev.Key() == tcell.KeyRune { // Escrever novo caracter na linha
					lineContent = lineContent[0:cursor] + string(ev.Rune()) + lineContent[cursor:]
					cursor++
				}
		}
	}
}

// Metodo para exibir uma tela de erro no terminal quando já existir um outro cliente editando a linha selecionada
func showErrorScreen(entryErr string) {
	screen, err := tcell.NewScreen()
	handleError(err)
	defer screen.Fini()
	err = screen.Init()
	handleError(err)

	screen.Clear()
	screen.SetContent(0, 0, '[', []rune("ERROR: " + entryErr + "]"), tcell.StyleDefault)
	screen.SetContent(0, 1, '[', []rune("Press ESC to go back to editor]"), tcell.StyleDefault)
	screen.Show()

	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape {
					return
				}
		}
	}
}

// Metodo para enviar um evento de requisicao de escrita para o modulo cliente
func sendWriteRequest(client *editor.Editor_Client_Module, line int, lineContent string) {
	client.Req <- editor.AppClientRequest{
		Type: editor.WRITE,
		Cursor: &line,
		Line: &lineContent,
	}
}

// Metodo para exibir o texto completo em tela
func drawText(screen tcell.Screen, text []string, line int) {
	screen.Clear()
	for i := range text {
		lineIdentifier := ' '
		if line == i {
			lineIdentifier = CURSOR
		}

		screen.SetContent(0, i, lineIdentifier, []rune(text[i]), tcell.StyleDefault)
	}
	screen.Show()
}

// Metodo para exibir uma linha selecionada em tela
func drawTextLine(screen tcell.Screen, lineContent string, cursor int) {
	screen.Clear()
	if len(lineContent) > 0 {
		for i, r := range lineContent {
			screen.SetContent(i, 0, r, nil, tcell.StyleDefault)
		}
	}
	screen.Show()
	screen.ShowCursor(cursor, 0)
}

// Metodo para tratar erros da aplicacao
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
