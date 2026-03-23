import socket
from common.utils import Bet

RESPONSE_OK = 3
RESPONSE_ERROR = 4
RESPONSE_WINNERS = 5

def _recv_header(sock):
    data = b""

    while not data.endswith(b"\n"):
        chunk = sock.recv(1)
        if not chunk:
            raise ConnectionError("client disconnected")
        data += chunk

    return data.decode("utf-8").strip()

def _parse_header(header: str) -> int:
    if not header.startswith("LEN:"):
        raise ValueError("invalid header")

    # A mi me llega "LEN:size;TYPE:number"
    parts = header.split(";")

    length = int(parts[0].split(":")[1])
    msg_type = int(parts[1].split(":")[1])

    return length, msg_type

def _recv_exact(sock, size: int) -> str:
    data = b""

    while len(data) < size:
        chunk = sock.recv(size - len(data))
        if not chunk:
            raise ConnectionError("client disconnected")
        data += chunk

    return data.decode("utf-8")


def recv_raw(sock):
    header = _recv_header(sock)
    length, msg_type = _parse_header(header)
    body = _recv_exact(sock, length)

    return body, msg_type

def parse_batch(msg: str) -> list[Bet]:
    bets = []

    lines = msg.split("\n")

    for line in lines:
        if not line.strip():
            continue

        parts = [p.strip() for p in line.split(";")]

        if len(parts) != 6:
            raise ValueError(f"invalid bet format: {line}")

        bet = Bet(
            parts[0],        # no está más hardcodeada la agencia
            parts[1],
            parts[2],
            parts[3],
            parts[4],
            parts[5],
        )

        bets.append(bet)

    return bets

def _send_message(sock: socket.socket, message: str, msg_type: int):
    body = message.encode("utf-8")
    header = f"LEN:{len(body)};TYPE:{msg_type}\n".encode("utf-8")

    sock.sendall(header + body)
    
def send_winners(sock: socket.socket, winners: []):
    message = "\n".join(winners)
    _send_message(sock, message, RESPONSE_WINNERS)

def send_error(sock: socket.socket):
    _send_message(sock, "error", RESPONSE_ERROR)

def send_ok(sock: socket.socket):
    _send_message(sock, "ok", RESPONSE_OK)