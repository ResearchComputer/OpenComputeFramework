import json
import requests

URL = "https://api.autoai.dev/inference"

def inference():
    resp = requests.post(URL, json={
        'model': 'microsoft/deberta-large-mnli',
        'params': {
            'prompt': "tell me about computer science?",
            'max_tokens': 32,
            'temperature': 0.7,
            'top_p': 1.0,
            'top_k': 40,
        }
    })
    resp = resp.json()
    print(resp)
    return resp

if __name__ == "__main__":
    inference()