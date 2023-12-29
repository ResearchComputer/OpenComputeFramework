import requests

endpoint = "http://140.238.223.13:8092"
peerId = "QmVVw6oykwy8q1siFBd7vdJ9hod2qkTJHHEqxajeUd4Y3N"

def test_forward():
    peer = {
        "service": [{
            "name": "triteia",
            "status": "online",
            "hardware": []
        }]
    }
    res = requests.get(endpoint + f"/v1/p2p/{peerId}/v1/service/triteia/metrics", json=peer)
    print(res.text)

test_forward()