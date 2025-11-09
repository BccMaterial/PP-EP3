package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func reverseString(s string) string {
	if len(s) <= 1 {
		return s
	}

	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// processMessage processa uma mensagem seguindo o padrão conhecido
func processMessage(line string) string {
	// Remove o \n no final se existir
	line = strings.TrimSuffix(line, "\n")

	// Procura pelo início da tag ([)
	startBracket := strings.Index(line, "[")
	if startBracket == -1 {
		// Se não tem o padrão esperado, inverte tudo
		return reverseString(line) + "\n"
	}

	// A cor é tudo antes do "["
	color := line[:startBracket]

	// Procura pelo final da tag (]:)
	endTag := strings.Index(line, "]:")
	if endTag == -1 {
		// Padrão incompleto, inverte tudo
		return reverseString(line) + "\n"
	}

	// O apelido está entre [ e ]:
	apelido := line[startBracket : endTag+2] // Inclui o "]:"

	// O restante é mensagem + Reset
	restante := line[endTag+2:]

	// Encontra o Reset no final
	reset := "\033[0m"
	mensagem := restante
	if strings.HasSuffix(restante, reset) {
		mensagem = restante[:len(restante)-len(reset)]
	}

	mensagemInvertida := reverseString(mensagem)
	return color + apelido + " " + mensagemInvertida + reset + "\n"
}

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	fmt.Println("Connected!")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		// Usamos um reader bufio para ler linha por linha
		buf := make([]byte, 1024)
		var buffer string

		for {
			n, err := conn.Read(buf)
			if n > 0 {
				buffer += string(buf[:n])

				// Processa linhas completas
				for {
					newlineIndex := strings.Index(buffer, "\n")
					if newlineIndex == -1 {
						break
					}

					line := buffer[:newlineIndex+1]
					buffer = buffer[newlineIndex+1:]

					// Processa a linha conforme o padrão
					processed := processMessage(line)
					fmt.Print(processed)
				}
			}

			if err != nil {
				// Processa qualquer dado restante no buffer
				if buffer != "" {
					processed := processMessage(buffer)
					fmt.Print(processed)
				}
				break
			}
		}

		log.Println("done")
		done <- struct{}{}
	}()

	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done
}
