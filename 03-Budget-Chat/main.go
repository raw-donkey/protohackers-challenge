package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net"
	"slices"
	"strings"
	"unicode"
)

// key for username, value for conn
var users map[string]net.Conn

func main() {
	slog.Info("listening on :8080")

	users = make(map[string]net.Conn)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("connection failed", "message", err.Error())
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	name := ""

	defer func() {
		slog.Info("connection closed")

		conn.Close()
	}()

	_, err := conn.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))
	if err != nil {
		slog.Error("write failed", "error", err.Error())
		return
	}

	slog.Info("reading for read name")

	scanner := bufio.NewScanner(conn)

	scanner.Scan()

	name = scanner.Text()

	slog.Info("connection input name", "name", name)

	if isInvalidName(name) {
		slog.Error("invalid name", "name", name)
		return
	}

	isPresenceNotification := presenceNotification(conn)
	if !isPresenceNotification {
		return
	}

	isNotificationSendSuccessful := newUserJoinsNotification(name)
	if !isNotificationSendSuccessful {
		return
	}

	users[name] = conn

	for scanner.Scan() {
		isSendMessageSuccessful := sendMessage(name, scanner.Text())
		if !isSendMessageSuccessful {
			slog.Error("sendMessage failed", "name", name, "message", scanner.Text())
			return
		}
	}

	isSuccess := userLeaveNotification(name)
	if !isSuccess {
		slog.Error("userLeaveNotification failed", "name", name)
	}
}

func isInvalidName(name string) bool {
	if len(name) < 1 {
		return true
	}

	return strings.ContainsFunc(name, func(r rune) bool {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return false
		}

		return true
	})
}

func presenceNotification(writer io.Writer) bool {
	baseMessage := "* The room contains:"

	usersIter := maps.Keys(users)

	usersList := slices.Collect(usersIter)

	length := len(usersList)

	usernames := ""

	for i := 0; i < length; i++ {
		usernames += usersList[i]
		if i < length-1 {
			usernames += ", "
		}
	}

	message := fmt.Sprintf("%s %s\n", baseMessage, usernames)

	slog.Info("presenceNotification", "message", message)

	_, err := writer.Write([]byte(message))
	if err != nil {
		slog.Error("presenceNotification write failed", "message", err.Error())
		return false
	}

	return true
}

func newUserJoinsNotification(name string) bool {
	message := fmt.Sprintf("* %s has entered the room\n", name)

	slog.Info("newUserJoinsNotification", "message", message)

	for _, conn := range users {
		_, err := conn.Write([]byte(message))
		if err != nil {
			slog.Error("NewUserJoinsNotification write failed", "message", err.Error())
			return false
		}
	}

	return true
}

func userLeaveNotification(name string) bool {
	message := fmt.Sprintf("* %s has left the room\n", name)

	slog.Info("userLeaveNotification", "message", message)

	delete(users, name)

	for _, conn := range users {
		_, err := conn.Write([]byte(message))
		if err != nil {
			slog.Error("userLeaveNotification write failed", "message", err.Error())
			return false
		}
	}

	return true
}

func sendMessage(username, message string) bool {
	slog.Info("sendMessage", "username", username, "message", message)

	formatMessage := fmt.Sprintf("[%s] %s\n", username, message)

	for name, conn := range users {
		if name != username {
			_, err := conn.Write([]byte(formatMessage))
			if err != nil {
				slog.Error("sendMessage write failed", "message", err.Error())
				return false
			}
		}
	}

	return true
}
