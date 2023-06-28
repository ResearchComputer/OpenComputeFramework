#!/usr/bin/env python
import os
import typer
from typing import Optional
from ocf_cli.lib.core.config import read_config
from ocf_cli.lib.pprint.nodes import pprint_nodes
from ocf_cli.lib.pprint.service import pprint_service

app = typer.Typer()
home_dir = os.path.expanduser("~")
default_ocf_home = os.path.join(home_dir, ".config", "ocf")
config_path = os.path.join(default_ocf_home, "cli.json")
config = read_config(config_path)

@app.command()
def list(entity: str):
    if entity == "node":
        pprint_nodes()
    elif entity == "service":
        pprint_service()
    else:
        print("[ERROR] Unknown: ", entity)

@app.command()
def main():
    print("OCF CLI")


if __name__ == "__main__":
    app()