package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// Cores - Firula pro chat =)
var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

type client chan<- string // canal de mensagem

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	// Mapeia os clientes conectados
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			// Broadcast de mensagens. Envio para todos
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleNewConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	ch <- Yellow + "[Servidor]: Por favor, digite seu apelido: " + Reset
	nick_input := bufio.NewScanner(conn)
	apelido := conn.RemoteAddr().String()

	for nick_input.Scan() {
		apelido = nick_input.Text()
		if apelido == "" || apelido == " " {
			ch <- Yellow + "[Servidor]: Por favor, digite um apelido válido: " + Reset
			continue
		}
		break
	}

	ch <- Yellow + "[Servidor]: Bem vindo, " + apelido + "!" + Reset
	messages <- Yellow + "[Servidor]: " + apelido + " entrou no chat" + Reset
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- Cyan + "[" + apelido + "]" + ": " + input.Text() + Reset
	}

	leaving <- ch
	messages <- Yellow + "[Servidor]: " + apelido + " saiu do chat" + Reset
	conn.Close()
}

func main() {
	fmt.Println("Iniciando servidor...")
	listener, err := net.Listen("tcp", "localhost:3000")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(Green + "Servidor iniciado!" + Reset)
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleNewConn(conn)
	}
}
