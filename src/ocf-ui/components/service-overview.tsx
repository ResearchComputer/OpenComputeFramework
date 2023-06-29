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

function constructService(nodes: any) {
  // format: ["service_name": "", "gpu_str": ""]
  // gpu_str: 1x NVIDIA GeForce RTX 3090
  let service:any = {}
  for (let node of nodes) {
    if (service[node.service]) {
      service[node.service].push(...node.gpus)
    } else {
      service[node.service] = node.gpus
    }
  }
  for (let s in service) {
    service[s] = gpuStr(service[s])
  }
  let serviceArr = []
  for (let s in service) {
    serviceArr.push({"service_name": s, "gpu_str": service[s]})
  }
  return serviceArr
}

export function ServiceOverview(nodes: any) {
  let services = constructService(nodes.nodes)
  return (
    <div className="space-y-8">
      {services.map((service: any, idx: Number) => (
        <div className="flex items-center" key={idx.toString()}>
        <Avatar className="h-9 w-9">
          <AvatarImage src="/avatars/01.png" alt="Avatar" />
          <AvatarFallback>INF</AvatarFallback>
        </Avatar>
        <div className="ml-4 space-y-1">
          <p className="text-sm font-medium leading-none">{service.service_name}</p>
          <p className="text-sm text-muted-foreground">
            {service.gpu_str}
          </p>
        </div>
      </div>
      ))}
    </div>
  )
}
