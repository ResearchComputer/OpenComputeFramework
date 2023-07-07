import json
import signal
import asyncio
from loguru import logger
from nats.aio.client import Client as NATS
from utils import get_visible_gpus_specs
import atexit

async def shutdown(signal, loop, nc, model_name, connection_notice):
    """Cleanup tasks tied to the service's shutdown."""
    logger.info(f"Gracefully shutting down {model_name} worker...")
    tasks = [t for t in asyncio.all_tasks() if t is not
             asyncio.current_task()]
    [task.cancel() for task in tasks]
    await asyncio.gather(*tasks)
    connection_notice['status'] = 'disconnected'
    await nc.publish("worker:status", bytes(f"{json.dumps(connection_notice)}", encoding='utf-8'))
    await nc.close()
    loop.stop()

class InferenceWorker():
    def __init__(self, model_name) -> None:
        self.model_name = model_name
        # todo(xiaozhe): get gpu specs from nvml
        self.nc = NATS()
        self.connection_notice = {}

    async def run(self, loop):
        await self.nc.connect("nats://localhost:8094")        
        await self.nc.subscribe(f"inference:{self.model_name}", "workers", self.process_request)
        self.connection_notice = {
            'service': f'inference:{self.model_name}',
            'gpus': get_visible_gpus_specs(),
            'client_id': self.nc.client_id,
            'status': 'connected'
        }
        await self.nc.publish("worker:status", bytes(f"{json.dumps(self.connection_notice)}", encoding='utf-8'))

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
        # atexit.register(exit_handler, self.nc, self.model_name, self.connection_notice)
        logger.info(f"Starting {self.model_name} worker...")
        loop = asyncio.get_event_loop()
        signals = (signal.SIGHUP, signal.SIGTERM, signal.SIGINT, signal.SIGQUIT, signal.SIGABRT, signal.SIGTSTP)
        for s in signals:
            loop.add_signal_handler(
                s, lambda s=s: asyncio.create_task(shutdown(s, loop, self.nc, self.model_name, self.connection_notice)))
        loop.run_until_complete(self.run(loop))
        loop.run_forever()
        loop.close()