// Aplicacao para disparar execução de cliente em sistema distribuído de editor de texto
// Grupo: Fernanda Ferreira de Mello, Gaya Isabel Pizoli, Vitor Lamas Esposito

/*
  Cria um processo cliente
  para cada processo: seu id único e a mesma lista de processos - o primeiro processo da lista deve ser sempre o servidor
  	o endereco de cada processo é o dado na lista, na posicao do seu id.
*/

package main

import (
	"fmt"
	"log"
	"os"
	"sistema-distribuido/editor"
	"strconv"
	"strings"
	"time"
	"github.com/gdamore/tcell/v2"
)

const CURSOR rune = '*' // identificar em qual linha o cursor do usuario está

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Please specify at least one address:port!")
		fmt.Println("go run runClient.go <id> <server address> <client addresses>")
		fmt.Println("go run runClient.go 1 127.0.0.1:5000  127.0.0.1:6001  127.0.0.1:7002")
		return
	}

	id, _ := strconv.Atoi(os.Args[1])
	addresses := os.Args[2:]

	var client *editor.Editor_Client_Module = editor.NewClient(addresses, id, true)
	fmt.Println(client)

	time.Sleep(5 * time.Second)

	screen, err := tcell.NewScreen()
	handleError(err)
	defer screen.Fini()

	err = screen.Init()
	handleError(err)

	client.Req <- editor.AppClientRequest{
		RequestType: editor.READ,
		Cursor: nil,
		Line: nil,
	}

	response := <- client.Ind
	text := strings.Split(response.Text, ",")
	currentLine := 0

	maxLines := len(text)
	maxColumns := len(text[0])

	//editor loop
	//while running, update state, draw screen, else close
	running := true
	for running {
		select { 
			case editorModuleResponse := <- client.Ind:
				text = strings.Split(editorModuleResponse.Text, ",")
			default:
				screen.Clear()
				drawText(screen, text, currentLine)
				screen.Show()

				// get event
				ev := screen.PollEvent()
				// check event type
				switch ev := ev.(type) {
					// check event key
					case *tcell.EventKey:
						if ev.Key() == tcell.KeyEscape {
							running = false
						} else if ev.Key() == tcell.KeyUp && currentLine > 0 {
							currentLine--
						} else if ev.Key() == tcell.KeyDown && currentLine < (maxLines - 1) {
							currentLine++
						} else if ev.Key() == tcell.KeyEnter {
							screen.Suspend()
							editLine(client, text, currentLine, maxColumns)
							screen.Resume()
						}
				}
		}
	}
}

func editLine(client *editor.Editor_Client_Module, text []string, line int, maxColumns int) {
	screen, err := tcell.NewScreen()
	handleError(err)
	defer screen.Fini()

	err = screen.Init()
	handleError(err)

	lineContent := text[line]
	cursor := 0

	for {
		screen.Clear()
		if len(lineContent) > 0 {
			for i, r := range lineContent {
				screen.SetContent(i, 0, r, nil, tcell.StyleDefault)
			}
		}
		screen.Show()
		screen.ShowCursor(cursor, 0)

		ev := screen.PollEvent()
		switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape {
					return
				} else if ev.Key() == tcell.KeyLeft && cursor > 0 {
					cursor--
				} else if ev.Key() == tcell.KeyRight && cursor < len(lineContent) {
					cursor++
				} else if ev.Key() == tcell.KeyBackspace && cursor > 0 && len(lineContent) > 0 {
					if (cursor == len(lineContent)) {
						lineContent = lineContent[0:cursor-1]
					} else if (cursor == 1) {
						lineContent = lineContent[cursor:]
					} else {
						lineContent = lineContent[0:cursor-1] + lineContent[cursor:]
					}
					cursor--
				} else if ev.Key() == tcell.KeyEnter {
					sendWriteRequest(client, line, lineContent)
					return
				} else if ev.Key() == tcell.KeyRune {
					lineContent = lineContent[0:cursor] + string(ev.Rune()) + lineContent[cursor:]
					cursor++
				}
		}
	}
}

func sendWriteRequest(client *editor.Editor_Client_Module, line int, lineContent string) {
	client.Req <- editor.AppClientRequest{
		RequestType: editor.WRITE,
		Cursor: &line,
		Line: &lineContent,
	}
}

func drawText(s tcell.Screen, text []string, line int) {

	for i := range text {
		lineIdentifier := ' '
		if line == i {
			lineIdentifier = CURSOR
		}

		s.SetContent(0, i, lineIdentifier, []rune(text[i]), tcell.StyleDefault)
	}
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
