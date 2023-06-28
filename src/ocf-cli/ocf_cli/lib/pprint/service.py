import requests
from rich.pretty import pprint
from rich.table import Table
from ocf_cli.lib.pprint._base import console

RELAY_URL = "https://inference.autoai.dev"

def pprint_service():
    url = f"{RELAY_URL}/api/v1/status/table"
    resp = requests.get(url)
    resp = resp.json()
    table = Table(title="Service")
    table.add_column("Service")
    table.add_column("Providers")
    services = {}
    for node in resp['nodes']:
        if node['service'] not in services:
            services[node['service']] = {'providers': []}
        services[node['service']]['providers'].append(node)
    for service in services:
        gpu_specs = {}
        for node in services[service]['providers']:
            for gpu in node["gpus"]:
                if gpu["name"] not in gpu_specs:
                    gpu_specs[gpu["name"]] = 0
                gpu_specs[gpu["name"]] += 1
        gpu_specs_str = ""
        for gpu_name, gpu_count in gpu_specs.items():
            gpu_specs_str += f"{gpu_count}x {gpu_name}, "
        gpu_specs_str = gpu_specs_str[:-2]
        table.add_row(
            service,
            gpu_specs_str,
        )
    console.print(table)