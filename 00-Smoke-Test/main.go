package main

import (
	"io"
	"log/slog"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		slog.Info("failed to listen target address", "error", err.Error())
		return
	}
	defer listen.Close()

	slog.Info("server started", "address", ":8080")

	for {
		conn, err := listen.Accept()
		if err != nil {
			slog.Error("failed to accept connection", "error", err.Error())
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()

	slog.Info("new connection", "remote_addr", addr)

	_, err := io.Copy(conn, conn)

	if err != nil {
		slog.Error("connection error", "error", err.Error())
		return
	}

	slog.Info("client disconnected", "addr", addr)
}
