import socket
import logging
from common.utils import Bet, store_bets

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)

    def close(self):
        if self._server_socket:
            try:
                self._server_socket.close()
                logging.debug("action: close_server_socket | result: success")
            except Exception as e:
                logging.error(f"action: close_server_socket | result: fail | error: {e}")
            finally:
                self._server_socket = None # Evito errores por doble llamados (posibles casos futuros)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while True:
            client_sock = self.__accept_new_connection()
            self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            msg = self.__recv_until_newline(client_sock)

            if not msg:
                return

            bet = self.__parse_bet(msg)

            store_bets([bet])

            logging.info(
                f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}"
            )

            client_sock.sendall(b"ok\n")

        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")

        finally:
            client_sock.close()


    def __recv_until_newline(self, sock):
        data = b"" # inicio el buffer vacío, para luego ir armando el mensaje de la apuesta

        while not data.endswith(b"\n"):
            chunk = sock.recv(1024) 

            if not chunk: # puede ocurrir que el cliente cierre la conexión
                raise ConnectionError("client disconnected before end of message") # esto lo captura el try/except de handle_client_connection
            
            data += chunk

        return data.decode("utf-8").rstrip("\n")
    
    def __parse_bet(self, msg):
        parts = [p.strip() for p in msg.strip().split(";")]

        if len(parts) != 5:
            raise ValueError(f"invalid bet format: {msg}")

        return Bet(
            "1",        # agencia hardcodeada por ahora
            parts[0],   # nombre
            parts[1],   # apellido
            parts[2],   # dni
            parts[3],   # nacimiento
            parts[4],   # numero
        )

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
