import { Card, CardContent, CardHeader } from "@/components/ui/card";
import {
  Network,
  Cpu,
  Shield,
  Zap,
  Globe,
  Server,
} from "lucide-react";

const features = [
  {
    icon: Network,
    title: "Peer-to-Peer Network",
    description:
      "Built on LibP2P, connecting computing resources in a decentralized mesh network without single points of failure.",
  },
  {
    icon: Cpu,
    title: "GPU Cluster Management",
    description:
      "Automatically discover and manage GPU resources across the network with hardware specification tracking.",
  },
  {
    icon: Shield,
    title: "Conflict-Free Protocol",
    description:
      "CRDT-based distributed state management ensures all nodes eventually agree on network state.",
  },
  {
    icon: Globe,
    title: "Global Service Discovery",
    description:
      "Dynamic service registration and discovery allows automatic routing of requests to available resources.",
  },
  {
    icon: Server,
    title: "ML Inference at Scale",
    description:
      "Run machine learning inference across distributed nodes with automatic load balancing and fault tolerance.",
  },
  {
    icon: Zap,
    title: "High Performance",
    description:
      "Optimized for low-latency communication with efficient request routing and parallel processing capabilities.",
  },
];

const Features = () => {
  return (
    <div
      id="features"
      className="max-w-screen-xl mx-auto w-full py-12 xs:py-20 px-6"
    >
      <h2 className="text-3xl xs:text-4xl md:text-5xl md:leading-[3.5rem] font-bold tracking-tight sm:max-w-xl sm:text-center sm:mx-auto">
        Boost Your Strategy with Smart Features
      </h2>
      <div className="mt-8 xs:mt-14 w-full mx-auto grid md:grid-cols-2 lg:grid-cols-3 gap-x-10 gap-y-12">
        {features.map((feature) => (
          <Card
            key={feature.title}
            className="flex flex-col border rounded-xl overflow-hidden shadow-none"
          >
            <CardHeader>
              <feature.icon />
              <h4 className="!mt-3 text-xl font-bold tracking-tight">
                {feature.title}
              </h4>
              <p className="mt-1 text-muted-foreground text-sm xs:text-[17px]">
                {feature.description}
              </p>
            </CardHeader>
            <CardContent className="mt-auto px-0 pb-0">
              <div className="bg-muted h-52 ml-6 rounded-tl-xl" />
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
};

export default Features;
