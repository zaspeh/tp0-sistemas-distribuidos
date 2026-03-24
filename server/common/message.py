import logging
from common.protocol import send_ok, send_error, parse_batch, RESPONSE_OK, RESPONSE_ERROR
from common.utils import Bet, store_bets

class Message:
    def handle(self, server, client_sock):
        raise NotImplementedError
    
class BatchMessage(Message):
    def __init__(self, body: str):
        self.bets = parse_batch(body)

    def handle(self, server, client_sock):
        with server.client_lock:
            server.client_agency[client_sock] = self.bets[0].agency
        cantidad = len(self.bets)

        try:
            with server.file_lock:
                store_bets(self.bets)

            logging.info(
                f"action: apuesta_recibida | result: success | cantidad: {cantidad}"
            )

            send_ok(client_sock)

            return False
        except Exception:
            logging.error(
                f"action: apuesta_recibida | result: fail | cantidad: {cantidad}"
            )

            send_error(client_sock)
            return True


class NotifyDoneMessage(Message):
    def handle(self, server, client_sock):
        return server.mark_client_done(client_sock)

class AskWinnersMessage(Message):
    def handle(self, server, client_sock):
        return server.handle_winners_request(client_sock)
