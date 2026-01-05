package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
)

var address = "7YWHMfk9JZe0LM0g1ZauHuiSxhI"

func main() {
	slog.Info("listening on :8080")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		slog.Error("error listening on :8080", "error", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("error accepting on :8080", "error", err)
			return
		}

		slog.Info("accepted new connection", "ip", conn.RemoteAddr())
		go handle(conn)
	}
}

func handle(client net.Conn) {
	name := ""

	server, err := net.Dial("tcp", "chat.protohackers.com:16963")
	if err != nil {
		slog.Error("error handling on :16963", "error", err)
		return
	}

	defer func() {
		message := fmt.Sprintf("* %s has left the room\n", name)
		slog.Info(message)

		client.Close()
		server.Close()
	}()

	//scanner := bufio.NewScanner(client)
	serverScanner := bufio.NewScanner(server)

	ch := make(chan struct{})

	// only read sever and send receive data from server
	go func(client net.Conn, serverScanner *bufio.Scanner) {
		first := true
		for serverScanner.Scan() {
			slog.Info("read from server", "data", serverScanner.Text())
			data := make([]byte, len(serverScanner.Bytes()))
			copy(data, serverScanner.Bytes())
			_, err := client.Write(append(modifyData(data), '\n'))
			if err != nil {
				slog.Error("error handling on :16963", "error", err)
				return
			}

			if first {
				first = false
				ch <- struct{}{}
			}
		}
	}(client, serverScanner)

	slog.Info("waiting to client")

	first := true

	<-ch

	reader := bufio.NewReader(client)

	for {
		text, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			slog.Error("error handling on :16963", "error", err)
			return
		} else if err == io.EOF {
			break
		}

		text = strings.TrimSuffix(text, "\n")

		if first {
			name = text
			first = false
		}

		slog.Info("read from client", "data", text)

		data := make([]byte, len(text))

		copy(data, text)

		data = append(modifyData(data), '\n')

		slog.Info("after modify", "data", data)

		_, err = server.Write(data)
		if err != nil {
			slog.Error("error writing to server", "error", err)
			return
		}
	}
}

func modifyData(data []byte) []byte {
	items := bytes.Split(data, []byte(" "))

	result := make([][]byte, 0, len(data))

	for _, item := range items {
		if isCryptoAddress(item) {
			slog.Info("it is crypto address", "data", item)
			result = append(result, []byte(address))
			continue
		}

		slog.Info("it is not crypto address", "data", item)

		result = append(result, item)
	}

	return bytes.Join(result, []byte(" "))
}

func isCryptoAddress(data []byte) bool {
	if len(data) < 26 || len(data) > 35 {
		return false
	}

	if data[0] != '7' {
		return false
	}

	return true
}
