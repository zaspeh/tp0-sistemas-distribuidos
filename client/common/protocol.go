package common

import (
	"bufio"
	"fmt"
	"net"
)

func (c *Client) SendBet() error {
	data := SerializeBet(c.config.Bet)

	totalWritten := 0
	for totalWritten < len(data) {
		n, err := c.conn.Write(data[totalWritten:])
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

func ReceiveConfirmation(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)

	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return response, nil
}
