import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/registry/new-york/ui/avatar"

function gpuStr(gpus: any) {
  // format: 1x NVIDIA GeForce RTX 3090
  let gpu_count:any = {}
  for (let gpu of gpus) {
    if (gpu_count[gpu.name]) {
      gpu_count[gpu.name] += 1
    } else {
      gpu_count[gpu.name] = 1
    }
  }
  let gpu_str = ""
  for (let gpu in gpu_count) {
    gpu_str += `${gpu_count[gpu]}x ${gpu}, `
  }
  return gpu_str.slice(0, -2)
}

export function NodeOverview(nodes: any) {
  let nnodes = nodes.nodes
  return (
    <div className="space-y-8">
      {nnodes.map((node: any) => (
        <div className="flex items-center">
        <Avatar className="h-9 w-9">
          <AvatarImage src="/avatars/01.png" alt="Avatar" />
          <AvatarFallback>{node.client_id}</AvatarFallback>
        </Avatar>
        <div className="ml-4 space-y-1">
          <p className="text-sm font-medium leading-none">{node.peer_id.slice(0, 4)}*{node.peer_id.slice(-8)}</p>
          <p className="text-sm text-muted-foreground">
            {node.service}
          </p>
        </div>
        <div className="ml-auto font-medium">{gpuStr(node.gpus)}</div>
      </div>
      ))}
    </div>
  )
}
