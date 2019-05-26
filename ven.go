/*
Copyright (c) 2019 Bailey Thompson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"github.com/bkthomps/Ven/screen"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		print("Usage: ven <file_name>\n")
		return
	}
	userArg := os.Args[1]
	if userArg == "-v" || userArg == "--version" {
		print("Ven version 0.3.0\n")
		print("Created by Bailey Thompson\n")
		print("Available at github.com/bkthomps/Ven\n")
		return
	}
	s, e := tcell.NewScreen()
	if e != nil {
		log.Fatal(e)
	}
	encoding.Register()
	quit := make(chan struct{})
	screen.Init(s, quit, userArg)
	<-quit
	s.Fini()
}
