package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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
	entering     = make(chan client)
	leaving      = make(chan client)
	messages     = make(chan string)
	conn_clients = make(map[string]chan string)
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
		fmt.Fprint(conn, msg)
	}
}

func handleNewConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	ch <- Yellow + "[Servidor]: Por favor, digite seu apelido: " + Reset
	nick_input := bufio.NewScanner(conn)
	ip := conn.RemoteAddr().String()
	apelido := ip

	for nick_input.Scan() {
		apelido = nick_input.Text()
		if apelido == "" || apelido == " " {
			ch <- Yellow + "[Servidor]: Por favor, digite um apelido válido: " + Reset
			continue
		}
		break
	}

	conn_clients[apelido] = ch
	ch <- Yellow + "[Servidor]: Bem vindo, " + apelido + "! Digite /help para os comandos." + Reset + "\n"
	messages <- Yellow + "[Servidor]: " + apelido + " entrou no chat" + Reset + "\n"
	fmt.Println(Green + apelido + " (" + ip + ") " + "entrou" + Reset)
	entering <- ch

	input := bufio.NewScanner(conn)

loop:
	for input.Scan() {
		text_input := input.Text()
		tokens := strings.Split(text_input, " ")
		first_word := tokens[0]

		switch first_word {
		case "/msg":

			if len(tokens) > 2 && strings.HasPrefix(tokens[1], "@") {
				channel_to_send := strings.Replace(tokens[1], "@", "", 1)
				message := strings.Join(tokens[2:], " ")
				pv_ch := conn_clients[channel_to_send]

				if pv_ch == nil {
					ch <- Red + "ERRO: Usuário " + channel_to_send + " não encontrado." + Reset + "\n"
					continue
				}

				pv_ch <- Magenta + "[" + apelido + "]: " + message + Reset + "\n"
				ch <- Magenta + "[" + apelido + "]: " + message + Reset + "\n"
				fmt.Println(Magenta + "[" + apelido + " (" + ip + ")" + "]" + " enviou uma mensagem privada para " + channel_to_send + Reset)
				continue
			}

			message := strings.Join(tokens[1:], " ")
			messages <- Cyan + "[" + apelido + "]" + ": " + message + Reset + "\n"
			fmt.Println(Cyan + "[" + apelido + " (" + ip + ")" + "]" + " enviou: " + text_input + Reset)
		case "/help":
			ch <- Blue + "------ COMANDOS ------" + Reset + "\n"
			ch <- Blue + "/help: " + "Mostra essa tela" + Reset + "\n"
			ch <- Blue + "/changenick APELIDO: " + "Altera seu apelido" + Reset + "\n"
			ch <- Blue + "/msg [@nick] @msg: " + "Envia uma mensagem, insira o @nick para uma mensagem privada." + Reset + "\n"
			ch <- Blue + "/exit: " + "Sai do chat" + Reset + "\n"
		case "/changenick":
			if len(tokens) < 2 {
				ch <- Red + "ERRO: /changenick requer um nome. Uso: /changenickname NOME" + Reset + "\n"
				continue
			}
			old_apelido := apelido
			nickname := tokens[1]
			apelido = nickname
			ch <- Green + "Nome alterado com sucesso!" + Reset + "\n"
			messages <- Blue + "[Servidor]: " + old_apelido + " alterou o apelido para " + apelido + "!" + Reset + "\n"
			fmt.Println(Blue + "[Servidor]: " + old_apelido + " (" + ip + ") " + "alterou o apelido para " + apelido + Reset)
		case "/exit":
			messages <- Yellow + "[Servidor]: " + apelido + " saiu do chat" + Reset + "\n"
			ch <- Yellow + "[Servidor]: Adeus, volte sempre!" + Reset + "\n"
			fmt.Println(Red + apelido + " (" + ip + ") " + "saiu" + Reset)
			leaving <- ch
			conn.Close()
			break loop
		default:
			ch <- Red + "ERRO: Insira um comando válido (/help lista os comandos)" + Reset + "\n"
		}
	}
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
