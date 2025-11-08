package main

import (
	"bufio"
	"fmt"
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

func messageReceiver(conn net.Conn, messages chan string) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		messages <- scanner.Text()
	}
	close(messages)
}

// Lida com as mensagens, mantendo o prompt
func messageHandler(done chan struct{}, messages chan string, prompt string) {
	for message := range messages {
		// Limpamos o prompt e repetimos o digite a mensagem. Dessa forma,
		// o prompt continua aparecendo mesmo recebendo uma mensagem.
		fmt.Print("\r\033[K")
		fmt.Println(message)
		fmt.Print(prompt)
	}
	done <- struct{}{}
}

func main() {
	prompt := "> "
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(Green + "Conectado com sucesso! Bem-vindo!" + Reset)

	done := make(chan struct{})
	messages := make(chan string)

	go messageReceiver(conn, messages)
	go messageHandler(done, messages, prompt)

	// Loop principal para enviar mensagens
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Print(prompt)
		message := scanner.Text()
		fmt.Fprintln(conn, message)
	}

	conn.Close()
	<-done
}
