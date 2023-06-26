from loguru import logger
def get_visible_gpus_specs():
    gpus = []
    try:
        from pynvml import nvmlInit, nvmlDeviceGetCount, nvmlDeviceGetHandleByIndex, nvmlDeviceGetMemoryInfo, nvmlDeviceGetName
        nvmlInit()
        deviceCount = nvmlDeviceGetCount()
        for i in range(deviceCount):
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