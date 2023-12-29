import requests

endpoint = "http://localhost:8092"

def update_peer():
    peer = {
        "service": [{
            "name": "triteia",
            "status": "online",
            "hardware": [],
            "port": "8000",
        }]
    }
    res = requests.post(endpoint + "/v1/dnt/_node", json=peer)
    print(res.json())

update_peer()