import json
import requests

URL = "http://140.238.214.47:8092/api/v1/request/inference"

def inference():
    resp = requests.post(URL, json={
        'model_name': 'microsoft/deberta-large-mnli',
        'params': {
            'prompt': "tell me about computer science?",
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