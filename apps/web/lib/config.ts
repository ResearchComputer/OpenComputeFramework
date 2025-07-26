// OCF Website Configuration
export const config = {
  // Project Information
  name: "Open Compute Framework",
  shortName: "OCF",
  description: "A peer-to-peer framework for decentralized computing built on LibP2P",
  
  // Links
  links: {
    github: "https://github.com/ResearchComputer/OpenComputeFramework",
    api: "https://api.research.computer",
    triteia: "https://api.research.computer/triteia/",
    inferencia: "https://api.research.computer/inferencia/",
    docker: "https://ghcr.io/xiaozheyao/ocf",
    releases: "https://github.com/ResearchComputer/OpenComputeFramework/releases",
    researchComputer: "https://research.computer",
    libp2p: "https://libp2p.io",
    issues: "https://github.com/ResearchComputer/OpenComputeFramework/issues",
    discussions: "https://github.com/ResearchComputer/OpenComputeFramework/discussions",
  },
  
  // API Endpoints
  endpoints: {
    nodeTable: (host: string) => `http://${host}:8092/v1/dnt/table`,
    peers: (host: string) => `http://${host}:8092/v1/dnt/peers`,
  },
  
  // Default Ports
  ports: {
    api: 8092,
    p2p: 43905,
  },
  
  // Supported Backends
  backends: {
    triteia: {
      name: "Triteia",
      description: "Large generative transformer models with OpenAI-compatible API",
      endpoint: "https://api.research.computer/triteia/",
    },
    inferencia: {
      name: "Inferencia",
      description: "General HuggingFace models for various ML tasks",
      endpoint: "https://api.research.computer/inferencia/",
    },
    deepspeed: {
      name: "DeepSpeed-MII",
      description: "Text-to-image and other specialized models",
      endpoint: "https://api.research.computer/deepspeed/",
    },
  },
};

export default config;
