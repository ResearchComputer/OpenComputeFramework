import aiohttp
import asyncio

def parse_hardware_info(hardware_info):
    """
    Parse hardware information and return a string representation.
    
    Args:
        hardware_info (dict): Dictionary containing hardware information
        
    Returns:
        str: String representation of the hardware in the format "Nx[Spec]"
    """
    if not hardware_info or "gpus" not in hardware_info or not hardware_info["gpus"]:
        return "Unknown"
    # Group GPUs by name
    gpu_counts = {}
    for gpu in hardware_info["gpus"]:
        name = gpu.get("name", "Unknown GPU")
        gpu_counts[name] = gpu_counts.get(name, 0) + 1
    # Format the output
    result = []
    for gpu_name, count in gpu_counts.items():
        result.append(f"{count}x {gpu_name}")
    return ", ".join(result)

async def get_all_models(endpoint: str, with_details: bool=False):
    available_models = []
    async with aiohttp.ClientSession() as session:
        async with session.get(endpoint) as response:
            if response.status != 200:
                raise Exception(f"Failed to fetch data from endpoint: {response.status}")
            data = await response.json()

    models = []
    for node_info in data.values():
        # Only include nodes that are currently connected
        # if not node_info.get('connected', False):
        #     continue
        if not node_info.get('service'):
            continue
        device_info = parse_hardware_info(node_info.get("hardware"))
        for service in node_info['service']:
            if not service.get('identity_group'):
                continue
            model_names = [identity[len('model='):] for identity in service['identity_group'] if identity.startswith('model=')]
            # Add each model to the list
            if with_details:
                models.extend({
                    'id': model_name,
                    'device': device_info,
                    'object': 'model',
                    'created': '0x',
                    'owner': '0x',
                    } for model_name in model_names)
            else:
                models.extend({
                    'id': model_name,
                    'object': 'model',
                    'created': '0x',
                    'owner': '0x',
                } for model_name in model_names)
    return models