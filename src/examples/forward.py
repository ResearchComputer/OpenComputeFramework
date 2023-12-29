import requests

endpoint = "http://localhost:8092"

def test_forward():
    peer = {
        "service": [{
            "name": "dnt",
            "status": "online",
            "hardware": []
        }]
    }
    res = requests.post(endpoint + "/v1/proxy/QmWxgDBrscNmiURmba196goATfG6fHrMniNDMei13YTCay/v1/chat", json=peer)
    print(res.json())

test_forward()