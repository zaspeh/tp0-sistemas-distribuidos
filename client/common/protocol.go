package common

import (
	"bufio"
	"fmt"
	"net"
)

func SendBet(conn net.Conn, bet *Bet) error {
	return SendBatch(conn, []*Bet{bet})
}

func SendBatch(conn net.Conn, bets []*Bet) error {
	data := SerializeBatch(bets)

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

func SerializeBet(b *Bet) []byte {
	msg := fmt.Sprintf("%s;%s;%s;%s;%s\n",
		b.config.Nombre,
		b.config.Apellido,
		b.config.DNI,
		b.config.Nacimiento,
		b.config.Numero,
	)
	return []byte(msg)
}

// para serializar el batch, es como serializar N bets...
func SerializeBatch(bets []*Bet) []byte {
	var result []byte

	for _, b := range bets {
		result = append(result, SerializeBet(b)...)
	}

	return result
}

func ReceiveConfirmation(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return response, nil
}
