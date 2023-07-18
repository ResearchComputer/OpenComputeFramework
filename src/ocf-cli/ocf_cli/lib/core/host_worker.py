from ocf_cli.lib.core.base import BaseWorker
from ocf_cli.lib.core.utils import get_visible_gpus_specs

"""
Host worker is the meta-worker that will always connect to the ocf-node, and manages the start/stop of other workers.
It also periodically report the status of the workers to the ocf-node.
"""

class HostWorker(BaseWorker):
    def __init__(self) -> None:
        self.service_name = "host"
        super().__init__(self.service_name)
    
    def get_connection_notice(self):
        notice = {
            'service': f"{self.service_name}.{self.nc.client_id}",
            'gpus': get_visible_gpus_specs(),
            'client_id': self.nc.client_id,
            'status': 'connected',
            'offering': []
        }
        return notice

    async def handle_requests(self, msgs):
        # i.e., only the last message is processed
        msgs = msgs[-1]
        print(msgs)