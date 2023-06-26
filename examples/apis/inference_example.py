import json
import requests
from multiprocessing import Pool

URL = "https://inference.autoai.dev/api/v1/request/inference"

def inference():
    resp = requests.post(URL, json={
        'model_name': 'togethercomputer/RedPajama-INCITE-7B-Chat',
        'params': {
            'prompt': "<human>: tell me about computer science?\n<bot>: ",
            'max_tokens': 32,
            'temperature': 0.7,
            'top_p': 1.0,
            'top_k': 40,
        }
    })
    resp = json.loads(resp.json()['data'])
    print(resp)
    return resp

if __name__ == "__main__":
    inference()