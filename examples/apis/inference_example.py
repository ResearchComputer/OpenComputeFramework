import json
import requests

URL = "https://inference.autoai.dev/api/v1/request/handler"

def inference():
    resp = requests.post(URL, json={
        'model_name': 'inference:mosaicml/mpt-7b-chat',
        'params': {
            'prompt': "<human>: tell me about computer science?\n<bot>: ",
            'max_tokens': 32,
            'temperature': 0.7,
            'top_p': 1.0,
            'top_k': 40,
        }
    })
    resp = resp.text
    print(resp)
    return resp

if __name__ == "__main__":
    inference()