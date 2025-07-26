"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Copy } from "lucide-react";
import { Button } from "@/components/ui/button";

const CodeExample = () => {
  const dockerCommand = `# Run OCF with Docker
docker run -it -p 8092:8092 -p 43905:43905 --rm --name ocf \\
  ghcr.io/xiaozheyao/ocf:dev start --mode standalone`;

  const pythonExample = `# Use OCF for ML inference
import requests

response = requests.post(
    url="https://api.research.computer/triteia/v1/chat/completions",
    headers={"Content-Type": "application/json"},
    json={
        "model": "meta-llama/Llama-2-7b-chat-hf",
        "messages": [
            {"role": "user", "content": "Hello, world!"}
        ],
        "max_tokens": 100
    }
)

print(response.json())`;

  const joinNetworkCommand = `# Join an existing OCF network
./ocf start --bootstrap.addr=/ip4/<ip>/tcp/43905/p2p/<peer-id>`;

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  return (
    <div className="w-full max-w-screen-xl mx-auto py-8 xs:py-16 px-6">
      <div className="text-center mb-12">
        <h2 className="text-3xl xs:text-4xl md:text-5xl !leading-[1.15] font-bold tracking-tighter">
          Quick Start
        </h2>
        <p className="mt-1.5 xs:text-lg text-muted-foreground max-w-2xl mx-auto">
          Get started with Open Compute Framework in minutes. Deploy standalone, join a network, or use our public API.
        </p>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              Quick Deploy
              <Badge variant="secondary">Docker</Badge>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="bg-muted p-4 rounded-lg relative">
              <Button
                size="sm"
                variant="ghost"
                className="absolute top-2 right-2 h-6 w-6 p-0"
                onClick={() => copyToClipboard(dockerCommand)}
              >
                <Copy className="h-3 w-3" />
              </Button>
              <pre className="text-sm overflow-x-auto">
                <code>{dockerCommand}</code>
              </pre>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              Use Public API
              <Badge variant="secondary">Python</Badge>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="bg-muted p-4 rounded-lg relative">
              <Button
                size="sm"
                variant="ghost"
                className="absolute top-2 right-2 h-6 w-6 p-0"
                onClick={() => copyToClipboard(pythonExample)}
              >
                <Copy className="h-3 w-3" />
              </Button>
              <pre className="text-xs overflow-x-auto">
                <code>{pythonExample}</code>
              </pre>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              Join Network
              <Badge variant="secondary">CLI</Badge>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="bg-muted p-4 rounded-lg relative">
              <Button
                size="sm"
                variant="ghost"
                className="absolute top-2 right-2 h-6 w-6 p-0"
                onClick={() => copyToClipboard(joinNetworkCommand)}
              >
                <Copy className="h-3 w-3" />
              </Button>
              <pre className="text-sm overflow-x-auto">
                <code>{joinNetworkCommand}</code>
              </pre>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default CodeExample;
