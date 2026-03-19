package common

import (
	"net"
	"time"

	"github.com/op/go-logging"
)

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
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(batches [][]*Bet) {
	for _, batch := range batches { // realizo procesamiento por batches

		err := c.createClientSocket()
		if err != nil { // debería validar que se cree bien el socket
			return
		}

		// ahora envío el batch (se encarga el protocolo)
		err = SendBatch(c.conn, batch)
		if err != nil {
			log.Errorf(
				"action: send_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		response, err := ReceiveConfirmation(c.conn)
		c.conn.Close()

		if err != nil {
			log.Errorf(
				"action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		if response == "ok\n" {
			log.Infof( // no se especificó que se diga nada en el cliente pero no puedo seguir mostrando el DNI y NUMERO
				"action: apuesta_enviada | result: success | cantidad: %d",
				len(batch),
			)
		}

		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
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
