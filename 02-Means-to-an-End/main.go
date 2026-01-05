package main

import (
	"encoding/binary"
	"io"
	"log/slog"
	"net"
)

type Transaction struct {
	timestamp int32
	price     int32
}

type Session struct {
	ip           string
	transactions []Transaction
	uniqueSet    map[int32]bool
}

func main() {
	slog.Info("Server listening on :8080")

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		return
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			slog.Error("accept error", "message", err.Error())
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	buf := make([]byte, 9)
	session := Session{
		ip:           conn.RemoteAddr().String(),
		transactions: make([]Transaction, 0, 1000),
		uniqueSet:    make(map[int32]bool),
	}

	for {
		_, err := io.ReadFull(conn, buf)
		if err != nil {
			UndefineBehaviour(conn)
			return
		}

		number1 := int32(binary.BigEndian.Uint32(buf[1:5]))
		number2 := int32(binary.BigEndian.Uint32(buf[5:9]))

		switch buf[0] {
		case 'I':
			success := Insert(&session, number1, number2)
			if !success {
				UndefineBehaviour(conn)
				return
			}
		case 'Q':
			average := Query(&session, number1, number2)
			WriteFixed4(conn, average)
		default:
			UndefineBehaviour(conn)
			return
		}
	}
}

func Insert(session *Session, timestamp int32, price int32) bool {
	if session.uniqueSet[timestamp] {
		slog.Error("duplicate price", "timestamp", timestamp, "price", price)
		return false
	}

	session.transactions = append(session.transactions, Transaction{timestamp, price})

	return true
}

func Query(session *Session, minimum int32, maximum int32) int32 {
	if maximum < minimum {
		return 0
	}

	transactions := FilterTransactions(session, minimum, maximum)

	length := len(transactions)

	if length == 0 {
		return 0
	}

	if length == 1 {
		return transactions[0].price
	}

	average := ComputeAverage(transactions)

	return average
}

func FilterTransactions(session *Session, minimum int32, maximum int32) []Transaction {
	result := make([]Transaction, 0, 1000)

	for _, transaction := range session.transactions {
		if transaction.timestamp >= minimum && transaction.timestamp <= maximum {
			result = append(result, transaction)
		}
	}

	return result
}

func UndefineBehaviour(conn net.Conn) {
	_, _ = conn.Write([]byte("undefine behaviour"))
}

func WriteFixed4(conn net.Conn, average int32) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(average))

	_, err := conn.Write(buf)
	if err != nil {
		return
	}
}

func ComputeAverage(transactions []Transaction) int32 {
	sum := int64(0)

	for _, transaction := range transactions {
		sum += int64(transaction.price)
	}

	return int32(sum / int64(len(transactions)))
}
