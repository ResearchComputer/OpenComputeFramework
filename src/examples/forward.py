import requests

endpoint = "http://140.238.223.13:8092"

def test_forward():
    peer = {
        "service": [{
            "name": "triteia",
            "status": "online",
            "hardware": []
        }]
    }
    res = requests.get(endpoint + "/v1/service/triteia/metrics", json=peer)
    print(res.text)

test_forward()