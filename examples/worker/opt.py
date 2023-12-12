import mii

from _base import InferenceWorker

class MIIWorker(InferenceWorker):
    def __init__(self, model_name) -> None:
        super().__init__(model_name)
        self.pipe = mii.pipeline(model_name)
        
    async def handle_requests(self, msg):
        def callback(response):
            print(f"recv: {response.response[0]}")
        self.pipe(["test"], streaming_fn=callback)
        return await super().handle_requests(msg)

if __name__=="__main__":
    worker = MIIWorker("facebook/opt-125m")
    worker.start()