package main

import (
	"log"
	"os"

	"github.com/bkthomps/Ven/screen"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

const version = "0.3.1"

func main() {
	if len(os.Args) != 2 {
		print("Usage: ven <file_name>\n")
		return
	}
	userArg := os.Args[1]
	if userArg == "-v" || userArg == "--version" {
		print("Ven version " + version + "\n")
		print("Created by Bailey Thompson\n")
		print("Available at github.com/bkthomps/Ven\n")
		return
	}
	tCellScreen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	encoding.Register()
	quit := make(chan struct{})
	s := &screen.Screen{}
	s.Init(tCellScreen, quit, userArg)
	<-quit
	tCellScreen.Fini()
}
