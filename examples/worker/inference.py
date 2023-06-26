from loguru import logger
from _base import InferenceWorker
class HFWorker(InferenceWorker):
    def __init__(self, model_name) -> None:
        super().__init__(model_name)        

    async def handle_requests(self, msg):
        logger.info(f"Processing request {msg}")
        return {"result": "hello world"}

if __name__=="__main__":
    worker = HFWorker("test")
    worker.start()