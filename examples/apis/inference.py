import requests
from constant import RELAY_URL, HOST_ID
from multiprocessing import Pool

def inference(i):
    url = f"{RELAY_URL}/api/v1/proxy/{HOST_ID}/api/v1/request/inference"
    resp = requests.post(url, json={
        'model_name': 'openlm-research/open_llama_7b',
        'params': {
            'prompt': "Hello!"
        }
    })
    print(resp.text)
    return resp.text

def global_inference(i):
    url = f"{RELAY_URL}/api/v1/request/inference"
    resp = requests.post(url, json={
        'model_name': 'openlm-research/open_llama_7b',
        'params': {
            'prompt': "Hello!"
        }
    })
    print(resp.json())
    return resp.text

if __name__ == "__main__":
    with Pool(1) as p:
        p.map(global_inference, range(5))