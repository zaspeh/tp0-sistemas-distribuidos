package common

import (
	"bufio"
	"net"
	"os"
	"strings"
	"time"

	"github.com/op/go-logging"
)

const MaxBatchBytes = 8 * 1024 // 8KB

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}

	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(datasetPath string, maxAmount int) {
	err := c.createClientSocket()
	if err != nil { // debería validar que se cree bien el socket
		return
	}
	defer c.conn.Close()

	// ahora envío de a batches para no cargarlo en memoria
	err = c.ProcessAndSendBatches(datasetPath, maxAmount)
	if err != nil {
		log_send_error(c.config.ID, err)
		return
	}

	// log.Infof("action: enviar_todos_los_batches | result: success ")
	// aviso que envié todas las bets
	err = SendDone(c.conn)
	if err != nil {
		log_send_error(c.config.ID, err)
		return
	}

	// log.Infof("action: nofiticar_envio_de_todas_las_apuestas | result: success ")
	// consulto los ganadores
	err = SendAskWinners(c.conn)
	if err != nil {
		log_send_error(c.config.ID, err)
		return
	}

	// log.Infof("action: esperando_a_recibir_ganadores | result: success ")
	// recibo los ganadores
	response, err := ReceiveMessage(c.conn)
	// log.Infof("action: recibir_ganadores | result: success ")
	if err != nil {
		log_send_error(c.config.ID, err)
		return
	}

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d",
		countLines(response),
	)

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) ProcessAndSendBatches(datasetPath string, maxAmount int) error {
	file, err := os.Open(datasetPath)
	if err != nil {
		return err
	}
	defer file.Close() // aseguro que se cierre el archivo

	scanner := bufio.NewScanner(file)

	var batch []*Bet
	currentSize := 0
	totalBets := 0

	for scanner.Scan() { // leo por lineas
		line := scanner.Text()
		parts := strings.Split(line, ",")

		if len(parts) != 5 { // salteo si no tiene 5 partes para evitar que rompa la ejecución (no mando esa bet)
			continue
		}

		bet := NewBet(BetConfig{
			Nombre:     parts[0],
			Apellido:   parts[1],
			DNI:        parts[2],
			Nacimiento: parts[3],
			Numero:     parts[4],
		})

		betSize := len(SerializeBet(c.config.ID, bet))

		// si no me entra más nada en el batch lo envío y lo reinicio
		if len(batch) >= maxAmount || currentSize+betSize > MaxBatchBytes {
			if err := c.sendBatchAndWait(batch); err != nil {
				return err
			}

			totalBets += len(batch)
			batch = nil
			currentSize = 0
		}

		batch = append(batch, bet)
		currentSize += betSize
	}

	// envío lo que resta
	if len(batch) > 0 {
		if err := c.sendBatchAndWait(batch); err != nil {
			return err
		}
		totalBets += len(batch)
	}

	return scanner.Err()
}

func (c *Client) sendBatchAndWait(batch []*Bet) error {
	err := SendBatch(c.conn, c.config.ID, batch)
	if err != nil {
		return err
	}

	data, err := ReceiveMessage(c.conn)
	if err != nil {
		return err
	}

	switch data {
	case ERROR_MSG:
		log.Errorf(
			"action: apuesta_enviada | result: fail | cantidad: %d",
			len(batch),
		)
		return c.sendBatchAndWait(batch) // reenvío si hubo un error

	case OK_MSG:
		log.Infof(
			"action: apuesta_enviada | result: success | cantidad: %d",
			len(batch),
		)
	}

	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			log.Errorf("action: close_socket | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return err
		}

		c.conn = nil
		log.Infof("action: close_socket | result: success | client_id: %v", c.config.ID)
	}
	return nil
}
