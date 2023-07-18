import os
from loguru import logger

import traceback
from typing import Union
import pynvml


def gpu_measure() -> Union[dict, None]:
    try:
        pynvml.nvmlInit()
        metrics = {"gpu": []}
        deviceCount = pynvml.nvmlDeviceGetCount()
        for i in range(deviceCount):
            handle = pynvml.nvmlDeviceGetHandleByIndex(i)
            name = pynvml.nvmlDeviceGetName(handle)
            mem = pynvml.nvmlDeviceGetMemoryInfo(handle)
            power = pynvml.nvmlDeviceGetPowerUsage(handle)
            utilitization = pynvml.nvmlDeviceGetUtilizationRates(handle)
            try:
                name = name.decode("utf-8")
            except Exception as e:
                pass
            metrics["gpu"].append(
                {
                    "product_name": name,
                    "fb_memory_usage": {
                        "total": mem.total / 1024 / 1024,
                        "used": mem.used / 1024 / 1024,
                        "free": mem.free / 1024 / 1024,
                    },
                    "utilization": utilitization.gpu,
                    "power_readings": {"power_draw": power / 1000},
                }
            )
    except pynvml.NVMLError as error:
        traceback.print_exc()
        print(error)
        metrics = None
    finally:
        pynvml.nvmlShutdown()
        return metrics

def get_visible_gpus_specs():
    # https://github.com/gpuopenanalytics/pynvml/issues/28
    os.environ["CUDA_DEVICE_ORDER"] = "PCI_BUS_ID"
    gpus = []
    try:
        from pynvml import nvmlInit, nvmlDeviceGetCount, nvmlDeviceGetHandleByIndex, nvmlDeviceGetMemoryInfo, nvmlDeviceGetName
        nvmlInit()
        if "CUDA_VISIBLE_DEVICES" in os.environ:
            ids = list(map(int, os.environ.get("CUDA_VISIBLE_DEVICES", "").split(",")))
        else:
            deviceCount = nvmlDeviceGetCount()
            ids = range(deviceCount)
        for i in ids:
            handle = nvmlDeviceGetHandleByIndex(i)
            meminfo = nvmlDeviceGetMemoryInfo(handle)
            gpus.append({
                'name': nvmlDeviceGetName(handle),
                'memory': meminfo.total,
                'memory_free': meminfo.free,
                'memory_used': meminfo.used,
            })
    except Exception as e:
        logger.info(f"No GPU found: {e}")
    return gpus