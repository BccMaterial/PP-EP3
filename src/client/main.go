package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(Green + "Conectado com sucesso! Bem-vindo!" + Reset)

	done := make(chan struct{})

	go func() {
		io.Copy(os.Stdout, conn)
		log.Println(Yellow + "[Servidor]: Fechando o chat" + Reset)
		done <- struct{}{} // sinaliza para a gorrotina principal
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done // espera a gorrotina terminar
}
