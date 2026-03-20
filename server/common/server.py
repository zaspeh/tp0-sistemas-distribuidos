import socket
import logging
from common.utils import Bet, store_bets
from common.protocol import parse_batch, recv_batch, send_message

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True

    def close(self):
        self._running = False  

        if self._server_socket:
            try:
                self._server_socket.close()
                logging.debug("action: close_server_socket | result: success")
            except Exception as e:
                logging.error(f"action: close_server_socket | result: fail | error: {e}")
            finally:
                self._server_socket = None

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        while self._running:
            try:
                client_sock = self.__accept_new_connection()
                self.__handle_client_connection(client_sock)

            except OSError:
                # El socket de cierra GRACEFUL entonces entro acá
                if not self._running:
                    break  
                logging.error("action: accept_connections | result: fail")

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            while True:
                cantidad = 0

                try:
                    msg = recv_batch(client_sock)

                    if not msg:
                        break

                    lines = [l for l in msg.split("\n") if l.strip()]
                    cantidad = len(lines)

                    bets = parse_batch(msg)

                    store_bets(bets)

                    logging.info(
                        f"action: apuesta_recibida | result: success | cantidad: {cantidad}"
                    )

                    send_message(client_sock, "ok")

                except Exception as e:
                    logging.error(
                        f"action: apuesta_recibida | result: fail | cantidad: {cantidad} | error: {e}"
                    )

                    try:
                        send_message(client_sock, "error")
                    except:
                        pass

                    break  # si algo falla cierro la conexión

        finally:
            client_sock.close()
        
    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
