package common

import (
	"bufio"
	"fmt"
	"net"
)

const (
	SEND_BATCH       = 0
	NOTIFY_DONE      = 1
	ASK_WINNERS      = 2
	RESPONSE_OK      = 3
	RESPONSE_ERROR   = 4
	RESPONSE_WINNERS = 5
	ERROR_MSG        = "error"
	OK_MSG           = "ok"
)

func writeAll(conn net.Conn, data []byte) error {
	total := 0

	for total < len(data) {
		n, err := conn.Write(data[total:])
		if err != nil {
			return err
		}
		total += n
	}

	return nil
}

func SendBatch(conn net.Conn, clientID string, bets []*Bet) error {
	data := SerializeBatch(clientID, bets)

	return writeAll(conn, data)
}

func SerializeBet(clientId string, b *Bet) []byte {
	msg := fmt.Sprintf("%s;%s;%s;%s;%s;%s\n",
		clientId,
		b.config.Nombre,
		b.config.Apellido,
		b.config.DNI,
		b.config.Nacimiento,
		b.config.Numero,
	)
	return []byte(msg)
}

// para serializar el batch, es como serializar N bets...
func SerializeBatch(clientId string, bets []*Bet) []byte {
	body := []byte{}

	for _, b := range bets {
		body = append(body, SerializeBet(clientId, b)...)
	}

	header := fmt.Sprintf("LEN:%d;TYPE:%d\n", len(body), SEND_BATCH)

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

func parseHeader(header string) (int, int, error) {
	var length, msgType int

	n, err := fmt.Sscanf(header, "LEN:%d;TYPE:%d", &length, &msgType)
	if err != nil || n == RESPONSE_ERROR {
		return 0, 0, fmt.Errorf("invalid header: %s", header)
	}

	return length, msgType, nil
}

func ReceiveMessage(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	header, err := recvHeader(reader)
	if err != nil {
		return "", err
	}

	length, _, err := parseHeader(header)
	if err != nil {
		return "", err
	}

	return recvExact(reader, length)
}

func SerializeDone() []byte {
	body := []byte("DONE")

	header := fmt.Sprintf("LEN:%d;TYPE:%d\n", len(body), NOTIFY_DONE)

	return append([]byte(header), body...)
}

func SerializeAskWinners() []byte {
	body := []byte("ASK")

	header := fmt.Sprintf("LEN:%d;TYPE:%d\n", len(body), ASK_WINNERS)

	return append([]byte(header), body...)
}

func SendDone(conn net.Conn) error {
	data := SerializeDone()
	return writeAll(conn, data)
}

func SendAskWinners(conn net.Conn) error {
	data := SerializeAskWinners()
	return writeAll(conn, data)
}
