import requests
from rich.pretty import pprint
from rich.table import Table
from ocf_cli.lib.pprint._base import console

RELAY_URL = "https://inference.autoai.dev"

def pprint_nodes():
    url = f"{RELAY_URL}/api/v1/status/table"
    resp = requests.get(url)
    resp = resp.json()
    table = Table(title="Connected Nodes")
    table.add_column("Node")
    table.add_column("Worker")
    table.add_column("Service")
    table.add_column("Status")
    table.add_column("GPU Device")
    table.add_column("GPU Memory (used / total)")
    for node in resp['nodes']:
        # make it 1x NVIDIA GeForce RTX 3090... etc.
        gpus_specs = {}
        for gpu in node["gpus"]:
            if gpu["name"] not in gpus_specs:
                gpus_specs[gpu["name"]] = 0
            gpus_specs[gpu["name"]] += 1
        gpu_specs_str = ""
        for gpu_name, gpu_count in gpus_specs.items():
            gpu_specs_str += f"{gpu_count}x {gpu_name}, "
        gpu_specs_str = gpu_specs_str[:-2]

        used_memory = 0
        total_memory = 0
        for gpu in node["gpus"]:
            used_memory += gpu["memory_used"]
            total_memory += gpu["memory"]
        memory_str = f"{used_memory/1024/1024:.2f} / {total_memory/1024/1024:.2f} MB"
        table.add_row(
            node["peer_id"],
            str(node["client_id"]),
            node["service"],
            node["status"],
            gpu_specs_str,
            memory_str,
        )
    console.print(table)