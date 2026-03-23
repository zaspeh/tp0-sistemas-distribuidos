package common

import (
	"bufio"
	"fmt"
	"net"
)

const (
	ERROR_MSG = "error"
	OK_MSG    = "ok"
)

func SendBatch(conn net.Conn, clientID string, bets []*Bet) error {
	data := SerializeBatch(clientID, bets)

	totalWritten := 0
	for totalWritten < len(data) {
		n, err := conn.Write(data[totalWritten:])
		if err != nil {
			return err
		}
		totalWritten += n
	}

	return nil
}

func SerializeBet(clientID string, b *Bet) []byte {
	msg := fmt.Sprintf("%s;%s;%s;%s;%s;%s\n",
		clientID,
		b.config.Nombre,
		b.config.Apellido,
		b.config.DNI,
		b.config.Nacimiento,
		b.config.Numero,
	)
	return []byte(msg)
}

// para serializar el batch, es como serializar N bets...
func SerializeBatch(clientID string, bets []*Bet) []byte {
	body := []byte{}

	for _, b := range bets {
		body = append(body, SerializeBet(clientID, b)...)
	}

	header := fmt.Sprintf("LEN:%d\n", len(body))

	return append([]byte(header), body...)
}

func recvHeader(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return line[:len(line)-1], nil
}

func recvExact(reader *bufio.Reader, size int) (string, error) {
	data := make([]byte, size)
	totalRead := 0

	for totalRead < size {
		n, err := reader.Read(data[totalRead:])
		if err != nil {
			return "", err
		}
		totalRead += n
	}

	return string(data), nil
}

func parseLength(header string) (int, error) {
	var length int

	_, err := fmt.Sscanf(header, "LEN:%d", &length)
	if err != nil {
		return 0, fmt.Errorf("invalid header: %s", header)
	}

	return length, nil
}

func ReceiveConfirmation(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	header, err := recvHeader(reader)
	if err != nil {
		return "", err
	}

	length, err := parseLength(header)
	if err != nil {
		return "", err
	}

	return recvExact(reader, length)
}
