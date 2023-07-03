import os
from loguru import logger

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