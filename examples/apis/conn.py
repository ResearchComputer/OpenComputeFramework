import requests
from constant import RELAY_URL, HOST_ID

def get_conn():
    url = f"{RELAY_URL}/api/v1/proxy/{HOST_ID}/api/v1/status/connections"
    resp = requests.get(url)
    return resp.text

def get_global_view():
    url = f"{RELAY_URL}/api/v1/proxy/{HOST_ID}/api/v1/status/global_view"
    resp = requests.get(url)
    return resp.text

if __name__ == "__main__":
    print(get_conn())