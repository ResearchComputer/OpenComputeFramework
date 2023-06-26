import json
import asyncio
from loguru import logger
from nats.aio.client import Client as NATS
from utils import get_visible_gpus_specs

class InferenceWorker():
    def __init__(self, model_name) -> None:
        self.model_name = model_name
        # todo(xiaozhe): get gpu specs from nvml
    
    async def run(self, loop):
        self.nc = NATS()
        await self.nc.connect("nats://localhost:8094")        
        await self.nc.subscribe(f"inference:{self.model_name}", "workers", self.process_request)
        connection_notice = {
            'service': f'inference:{self.model_name}',
            'gpus': get_visible_gpus_specs(),
            'client_id': self.nc.client_id,
            'status': 'connected'
        }
        await self.nc.publish("worker:status", bytes(f"{json.dumps(connection_notice)}", encoding='utf-8'))

    async def process_request(self, msg):
        processed_msg = json.loads(msg.data.decode())
        result = await self.handle_requests(processed_msg['params'])
        await self.reply(msg, result)

    async def handle_requests(self, msg):
        raise NotImplementedError

    async def reply(self, msg, data):
        data = json.dumps(data)
        await self.nc.publish(msg.reply, bytes(data, encoding='utf-8'))
    
    def start(self):
        logger.info(f"Starting {self.model_name} worker...")
        
        loop = asyncio.get_event_loop()
        loop.run_until_complete(self.run(loop))
        loop.run_forever()
        loop.close()