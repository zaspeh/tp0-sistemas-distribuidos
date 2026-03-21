from common.message import BatchMessage, NotifyDoneMessage, AskWinnersMessage

SEND_BATCH  = 0
NOTIFY_DONE = 1
ASK_WINNERS = 2

def build_message(body, msg_type):
    if msg_type == SEND_BATCH:
        return BatchMessage(body)
    elif msg_type == NOTIFY_DONE:
        return NotifyDoneMessage()
    elif msg_type == ASK_WINNERS:
        return AskWinnersMessage()