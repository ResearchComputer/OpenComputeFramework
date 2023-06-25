import requests
from multiprocessing import Pool

URL = "https://inference.autoai.dev/api/v1/request/inference"

def global_inference(i):
    resp = requests.post(URL, json={
        'model_name': 'togethercomputer/RedPajama-INCITE-7B-Base',
        'params': {
            'prompt': "Alan Turing was a "
        }
    })
    print(resp.json())
    return resp.text

if __name__ == "__main__":
    with Pool(1) as p:
        p.map(global_inference, range(1))