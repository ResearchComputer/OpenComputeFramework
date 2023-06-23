import asyncio
from nats.aio.client import Client as NATS

class InferenceWorker():
    def __init__(self, model_name) -> None:
        self.model_name = model_name
        # todo(xiaozhe): get gpu specs from nvml
        self.nc = None
    
    async def run(self, loop):
        self.nc = NATS()
        await self.nc.connect("nats://localhost:8094")        
        await self.nc.subscribe(f"inference:{self.model_name}", "workers", self.handle_requests)

    async def handle_requests(self, msg):
        raise NotImplementedError

    async def reply(self, msg, data):
        await self.nc.publish(reply=msg.reply, payload=bytes(data))
    
    def start(self):
        loop = asyncio.get_event_loop()
        loop.run_until_complete(self.run(loop))
        loop.run_forever()
        loop.close()