import socket
from common.utils import Bet

def _recv_header(sock):
    data = b""

    while not data.endswith(b"\n"):
        chunk = sock.recv(1)
        if not chunk:
            raise ConnectionError("client disconnected")
        data += chunk

    return data.decode("utf-8").strip()

def _parse_length(header: str) -> int:
    if not header.startswith("LEN:"):
        raise ValueError("invalid header")

    return int(header.split(":")[1])

def _recv_exact(sock, size: int) -> str:
    data = b""

    while len(data) < size:
        chunk = sock.recv(size - len(data))
        if not chunk:
            raise ConnectionError("client disconnected")
        data += chunk

    return data.decode("utf-8")


def recv_batch(sock: socket.socket) -> str:
    header = _recv_header(sock)
    length = _parse_length(header)
    return _recv_exact(sock, length)

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
            parts[0],        # agencia deja de estar hardcodeada
            parts[1],
            parts[2],
            parts[3],
            parts[4],
            parts[5],
        )

        bets.append(bet)

    return bets

def send_message(sock: socket.socket, message: str):
    body = message.encode("utf-8")
    header = f"LEN:{len(body)}\n".encode("utf-8")

    sock.sendall(header + body)