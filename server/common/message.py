import logging
from common.protocol import parse_batch, recv_msg, send_message
from common.utils import Bet, store_bets

class Message:
    def handle(self, server, client_sock):
        raise NotImplementedError
    
class BatchMessage(Message):
    def __init__(self, body: str):
        self.bets = parse_batch(body)

    def handle(self, server, client_sock):
        self.client_agency[client_sock] = self.bets[0].agency
        cantidad = len(self.bets)

        try:
            store_bets(self.bets)

            logging.info(
                f"action: apuesta_recibida | result: success | cantidad: {cantidad}"
            )

            send_message(client_sock, "ok")
            return True
        except Exception:
            logging.error(
                f"action: apuesta_recibida | result: fail | cantidad: {cantidad}"
            )

            send_message(client_sock, "error")
            client_sock.close()


class NotifyDoneMessage(Message):
    def handle(self, server, client_sock):
        return server.mark_client_done(client_sock)

class AskWinnersMessage(Message):
    def handle(self, server, client_sock):
        return server.handle_winners_request(client_sock)
