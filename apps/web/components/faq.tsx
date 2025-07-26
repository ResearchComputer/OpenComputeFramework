import {
  Accordion,
  AccordionContent,
  AccordionItem,
} from "@/components/ui/accordion";
import { cn } from "@/lib/utils";
import * as AccordionPrimitive from "@radix-ui/react-accordion";
import { PlusIcon } from "lucide-react";

const faq = [
  {
    question: "What is Open Compute Framework?",
    answer:
      "OCF is a peer-to-peer framework for decentralized computing built on LibP2P. It enables you to connect computing resources across a distributed network and run ML inference at scale without single points of failure.",
  },
  {
    question: "How do I get started with OCF?",
    answer:
      "Download the binary from our GitHub releases or build from source with Go 1.22.5+. You can run a standalone instance with Docker or join an existing network by connecting to bootstrap peers.",
  },
  {
    question: "What types of workloads does OCF support?",
    answer:
      "OCF currently supports machine learning inference workloads including large language models (via Triteia), text-to-image models (via DeepSpeed-MII), and general HuggingFace models (via Inferencia).",
  },
  {
    question: "How does the distributed architecture work?",
    answer:
      "OCF uses a CRDT-based distributed state system to maintain network consensus. Each node can discover services, route requests, and automatically balance load across available GPU resources in the network.",
  },
  {
    question: "Can I contribute my GPU resources?",
    answer:
      "Yes! You can join the network as a worker node by connecting to existing bootstrap peers. Your GPU resources will be automatically discovered and made available for inference tasks across the network.",
  },
  {
    question: "Is OCF suitable for production use?",
    answer:
      "OCF is currently in active development. We run public inference APIs at api.research.computer, but please evaluate thoroughly for your production requirements and consider the alpha status.",
  },
];

const FAQ = () => {
  return (
    <div id="faq" className="w-full max-w-screen-xl mx-auto py-8 xs:py-16 px-6">
      <h2 className="md:text-center text-3xl xs:text-4xl md:text-5xl !leading-[1.15] font-bold tracking-tighter">
        Frequently Asked Questions
      </h2>
      <p className="mt-1.5 md:text-center xs:text-lg text-muted-foreground">
        Quick answers to common questions about Open Compute Framework.
      </p>

      <div className="min-h-[550px] md:min-h-[320px] xl:min-h-[300px]">
        <Accordion
          type="single"
          collapsible
          className="mt-8 space-y-4 md:columns-2 gap-4"
        >
          {faq.map(({ question, answer }, index) => (
            <AccordionItem
              key={question}
              value={`question-${index}`}
              className="bg-accent py-1 px-4 rounded-xl border-none !mt-0 !mb-4 break-inside-avoid"
            >
              <AccordionPrimitive.Header className="flex">
                <AccordionPrimitive.Trigger
                  className={cn(
                    "flex flex-1 items-center justify-between py-4 font-semibold tracking-tight transition-all hover:underline [&[data-state=open]>svg]:rotate-45",
                    "text-start text-lg"
                  )}
                >
                  {question}
                  <PlusIcon className="h-5 w-5 shrink-0 text-muted-foreground transition-transform duration-200" />
                </AccordionPrimitive.Trigger>
              </AccordionPrimitive.Header>
              <AccordionContent className="text-[15px]">
                {answer}
              </AccordionContent>
            </AccordionItem>
          ))}
        </Accordion>
      </div>
    </div>
  );
};

export default FAQ;
