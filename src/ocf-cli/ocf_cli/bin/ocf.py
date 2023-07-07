#!/usr/bin/env python
import os
import time
import sched
import typer
from loguru import logger
import netifaces as ni
from typing import Optional
from netifaces import AF_INET
from ocf_cli.lib.core.utils import gpu_measure
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

@app.command()
def join(host: str='localhost',
         nic_name: str ='access',
         report_interval: int=10, 
         working_dir: str = "."
        ):
    print("> Joining OCF network")
    ip_addr = ni.ifaddresses(nic_name)[AF_INET][0]['addr']
    gpu_stats = gpu_measure()
    if gpu_stats is not None and 'gpu' in gpu_stats:
        print(">> GPU found, joining OCF network")
        total_gpus = len(gpu_stats)
        # s = sched.scheduler(time.time, time.sleep)
        
        # s.enter(report_interval, 1, clock_watch, (s, tom_client, ip_addr, str(idx)))
        # s.run(blocking=True)
    else:
        logger.error("No GPU found, exiting...")

if __name__ == "__main__":
    app()