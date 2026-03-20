import socket
from common.utils import Bet


def recv_batch(sock: socket.socket) -> str:
    data = b""

    while not data.endswith(b"\n\n"):
        chunk = sock.recv(1024)

        if not chunk:
            raise ConnectionError("client disconnected before end of batch")

        data += chunk

    return data.decode("utf-8").strip()


def parse_batch(msg: str) -> list[Bet]:
    bets = []

    lines = msg.split("\n")

    for line in lines:
        if not line.strip():
            continue

        parts = [p.strip() for p in line.split(";")]

        if len(parts) != 5:
            raise ValueError(f"invalid bet format: {line}")

        bet = Bet(
            "1",        # agencia hardcodeada
            parts[0],
            parts[1],
            parts[2],
            parts[3],
            parts[4],
        )

        bets.append(bet)

    return bets


def send_ok(sock: socket.socket):
    sock.sendall(b"ok\n")


def send_error(sock: socket.socket):
    sock.sendall(b"error\n")