import socket
import logging
from common.utils import Bet, load_bets, has_won
from common.message_factory import build_message
from common.protocol import recv_raw, send_message, RESPONSE_WINNERS

class Server:
    def __init__(self, port, listen_backlog, total_clients):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True
        self.client_agency = {}
        self.finished_clients = set()
        self.waiting_winners = []
        self.total_clients = total_clients
        self.sorteo_done = False

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
        while True: # termino cuando me pregunta los ganadores
            try:
                body, msg_type = recv_raw(client_sock)
                msg = build_message(body, msg_type)

                should_break = msg.handle(self, client_sock)

                if should_break == True:
                    break

            except ConnectionError:
                # client_sock.close() -> debería cerrar el socket?
                break # si el cliente se desconecta... entonces me voy

        
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

    def mark_client_done(self, client_sock):
        logging.info(f"action: agregar_cliente_listo | result: success | clientes_totales = {self.total_clients}")
        self.finished_clients.add(client_sock)

        if len(self.finished_clients) == self.total_clients:
            logging.info("action: sorteo | result: success")
            self.sorteo_done = True

            self._choose_winners()

            for sock in self.waiting_winners:
                self._send_winners(sock)

            self.waiting_winners = {}

        return False


    def _choose_winners(self):
        self.winners_by_agency = {}

        for bet in load_bets():
            if has_won(bet):
                self.winners_by_agency.setdefault(bet.agency, []).append(bet.document)

    def handle_winners_request(self, client_sock):
        if not self.sorteo_done:
            self.waiting_winners.append(client_sock)
            return True

        self._send_winners(client_sock)
        return True

    def _send_winners(self, client_sock):
        agency = self.client_agency[client_sock]

        winners = [
            b.document
            for b in load_bets()
            if b.agency == agency and has_won(b)
        ]

        response = "\n".join(winners)
        send_message(client_sock, response, RESPONSE_WINNERS)
        client_sock.close()