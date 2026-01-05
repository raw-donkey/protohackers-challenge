package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net"
)

var database map[string]string

func init() {
	database = make(map[string]string)

	database["version"] = "Ken's Key-Value Store 1.0"
}

func main() {
	slog.Info("listening on :8080")
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		slog.Error("failed to resolve UDP address")
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		slog.Error("error listening on :8080", "message", err.Error())
		return
	}

	buffer := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buffer)

		if err != nil {
			slog.Error("error reading from UDP", "message", err.Error())
			return
		}

		slog.Info("read from UDP", "message", string(buffer[:n]))

		data := make([]byte, n)

		copy(data, buffer[:n])

		handle(conn, addr, data)
	}
}

func handle(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	if bytes.Contains(data, []byte{'='}) {
		insert(data)
		return
	}

	retrieve(conn, addr, data)
}

func insert(data []byte) {
	slog.Info("inserting data to database")

	firstIndex := bytes.IndexByte(data, '=')
	key := string(data[:firstIndex])

	if key == "version" {
		slog.Info("can't modify version key")
		return
	}

	value := string(data[firstIndex+1:])

	database[key] = value
}

func retrieve(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	slog.Info("retrieving data from UDP")

	key := string(data)

	slog.Info("key value", "value", key)

	value, exists := database[key]
	if !exists {
		slog.Error("key not found", "key", key)
		response := fmt.Sprintf("%s=", key)
		_, err := conn.WriteToUDP([]byte(response), addr)
		if err != nil {
			slog.Error("error writing to UDP", "message", err.Error())
		}

		return
	}

	response := fmt.Sprintf("%s=%s", key, value)

	slog.Info("retrieved data", "value", value)

	_, err := conn.WriteToUDP([]byte(response), addr)
	if err != nil {
		slog.Error("error writing to UDP", "message", err.Error())
		return
	}
}
