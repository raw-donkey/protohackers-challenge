package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net"
	"strings"
)

type PrimeTime struct {
	Method string  `json:"method"`
	Number float64 `json:"number"`
}

func main() {
	slog.Info("Server listening on :8080")

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		slog.Error("listening error", "message", err.Error())
		return
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			slog.Error("connection error", "message", err.Error())
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		slog.Info("read content", "data", scanner.Text())

		if !checkFormat(scanner.Text()) {
			slog.Error("format failed", "data", scanner.Text())

			sendSingleMalformedResponse(conn)

			return
		}

		var primeTime PrimeTime

		err := json.Unmarshal(scanner.Bytes(), &primeTime)
		if err != nil {
			slog.Error("unmarshal error", "message", err.Error())

			sendSingleMalformedResponse(conn)

			return
		}

		slog.Info("json unmarshal data", "primeTime", primeTime)

		if isMalFormedResponse(primeTime) {
			slog.Error("data is malformed")

			sendSingleMalformedResponse(conn)

			return
		}

		if isNotInteger(primeTime) {
			sendPrimeResultResponse(conn, false)

			continue
		}

		sendPrimeResultResponse(conn, isPrime(primeTime.Number))
	}
}

func checkFormat(text string) bool {
	if !strings.Contains(text, "method") || !strings.Contains(text, "number") {
		return false
	}

	return true
}

func isPrime(number float64) bool {
	if number <= 1 {
		return false
	}

	for i := 2.0; i <= math.Sqrt(number); i++ {
		if math.Mod(number, i) == 0 {
			return false
		}
	}

	return true
}

func isMalFormedResponse(primeTime PrimeTime) bool {
	if primeTime.Method != "isPrime" {
		return true
	}

	return false
}

func isNotInteger(primeTime PrimeTime) bool {
	if _, frac := math.Modf(primeTime.Number); frac != 0.0 {
		return true
	}

	return false
}

func sendSingleMalformedResponse(conn net.Conn) {
	_, err := fmt.Fprintln(conn, "malformed")
	if err != nil {
		slog.Error("write error", "message", err.Error())
	}
}

func sendPrimeResultResponse(conn net.Conn, isPrime bool) {
	data, err := json.Marshal(map[string]any{
		"method": "isPrime",
		"prime":  isPrime,
	})

	slog.Info("prime result response", "data", data)

	if err != nil {
		return
	}

	_, err = fmt.Fprintln(conn, string(data))

	if err != nil {
		slog.Error("write error", "message", err.Error())
	}
}
