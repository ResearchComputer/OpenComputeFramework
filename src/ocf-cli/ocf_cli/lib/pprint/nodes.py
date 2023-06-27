import requests
from rich.pretty import pprint
RELAY_URL = "https://inference.autoai.dev/"

def pprint_nodes():
    url = f"{RELAY_URL}/api/v1/status/table"
    resp = requests.get(url)
    resp = resp.json()
    pprint(resp)