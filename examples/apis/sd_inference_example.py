import json
import requests

URL = "https://inference.autoai.dev/api/v1/request/inference"

def inference():
    resp = requests.post(URL, json={
        'model_name': 'stabilityai/stable-diffusion-xl-base-0.9',
        'params': {
            'prompt': "An astronaut is running on mars",
        }
    })
    resp = resp.json()
    return resp

if __name__ == "__main__":
    inference()