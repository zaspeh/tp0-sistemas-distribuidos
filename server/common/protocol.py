import socket
from common.utils import Bet



def recv_until_newline(sock: socket.socket) -> str:
        data = b"" # inicio el buffer vacío, para luego ir armando el mensaje de la apuesta

        while not data.endswith(b"\n"):
            chunk = sock.recv(1024) 

            if not chunk: # puede ocurrir que el cliente cierre la conexión
                raise ConnectionError("client disconnected before end of message") # esto lo captura el try/except de handle_client_connection
            
            data += chunk

        return data.decode("utf-8").rstrip("\n")

def parse_bet(msg: str) -> Bet:
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


def send_ok(sock: socket.socket) -> None:
    sock.sendall(b"ok\n")