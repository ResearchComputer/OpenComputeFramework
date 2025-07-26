"use client";

"use client";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  Carousel,
  CarouselApi,
  CarouselContent,
  CarouselItem,
} from "@/components/ui/carousel";
import { cn } from "@/lib/utils";
import { StarIcon } from "lucide-react";
import Image from "next/image";
import { useEffect, useState } from "react";

const testimonials = [
  {
    id: 1,
    name: "Dr. Sarah Chen",
    designation: "ML Research Lead",
    company: "Academic AI Lab",
    testimonial:
      "OCF revolutionized our research workflow. We can now run large-scale ML experiments across distributed GPU clusters without worrying about infrastructure costs. " +
      "The peer-to-peer architecture makes it incredibly resilient and cost-effective for academic research.",
    avatar: "https://randomuser.me/api/portraits/women/1.jpg",
  },
  {
    id: 2,
    name: "Marcus Rodriguez",
    designation: "DevOps Engineer",
    company: "TechStartup Inc",
    testimonial:
      "Setting up our distributed inference pipeline with OCF was surprisingly straightforward. The automatic service discovery and load balancing saved us months of development time. " +
      "Our ML services now scale seamlessly across multiple regions.",
    avatar: "https://randomuser.me/api/portraits/men/2.jpg",
  },
  {
    id: 3,
    name: "Dr. Priya Patel",
    designation: "Computer Vision Researcher",
    company: "Research Institute",
    testimonial:
      "The CRDT-based consensus mechanism in OCF is brilliant. We can have multiple researchers contributing compute resources without worrying about network partitions or conflicts. " +
      "It's like having a global supercomputer for our vision models.",
    avatar: "https://randomuser.me/api/portraits/women/3.jpg",
  },
  {
    id: 4,
    name: "Alex Thompson",
    designation: "Infrastructure Architect",
    company: "CloudScale Systems",
    testimonial:
      "OCF's P2P architecture eliminated our single points of failure. We can now offer ML inference services with 99.9% uptime across our distributed nodes. " +
      "The automatic failover and request routing make our infrastructure incredibly resilient.",
    avatar: "https://randomuser.me/api/portraits/men/4.jpg",
  },
  {
    id: 5,
    name: "Linda Wang",
    designation: "Data Science Lead",
    company: "FinTech Solutions",
    testimonial:
      "The cost savings from using OCF are incredible. Instead of paying for expensive cloud GPUs, we leverage community resources and contribute our idle compute back. " +
      "It's sustainable, cost-effective, and perfect for our fluctuating ML workloads.",
    avatar: "https://randomuser.me/api/portraits/women/5.jpg",
  },
  {
    id: 6,
    name: "David Kim",
    designation: "Research Engineer",
    company: "AI Research Group",
    testimonial:
      "OCF's service discovery is amazing. We can deploy new models and they're automatically available across the network. " +
      "The LibP2P foundation gives us incredible flexibility while maintaining performance. " +
      "Perfect for our research collaboration needs.",
    avatar: "https://randomuser.me/api/portraits/men/6.jpg",
  },
];
const Testimonial = () => {
  const [api, setApi] = useState<CarouselApi>();
  const [current, setCurrent] = useState(0);
  const [count, setCount] = useState(0);

  useEffect(() => {
    if (!api) {
      return;
    }

    setCount(api.scrollSnapList().length);
    setCurrent(api.selectedScrollSnap() + 1);

    api.on("select", () => {
      setCurrent(api.selectedScrollSnap() + 1);
    });
  }, [api]);

  return (
    <div
      id="testimonials"
      className="w-full max-w-screen-xl mx-auto py-6 xs:py-12 px-6"
    >
      <h2 className="mb-8 xs:mb-14 text-4xl md:text-5xl font-bold text-center tracking-tight">
        Testimonials
      </h2>
      <div className="container w-full mx-auto">
        <Carousel setApi={setApi}>
          <CarouselContent>
            {testimonials.map((testimonial) => (
              <CarouselItem key={testimonial.id}>
                <TestimonialCard testimonial={testimonial} />
              </CarouselItem>
            ))}
          </CarouselContent>
        </Carousel>
        <div className="flex items-center justify-center gap-2">
          {Array.from({ length: count }).map((_, index) => (
            <button
              key={index}
              onClick={() => api?.scrollTo(index)}
              className={cn("h-3.5 w-3.5 rounded-full border-2", {
                "bg-primary border-primary": current === index + 1,
              })}
            />
          ))}
        </div>
      </div>
    </div>
  );
};

const TestimonialCard = ({
  testimonial,
}: {
  testimonial: (typeof testimonials)[number];
}) => (
  <div className="mb-8 bg-accent rounded-xl py-8 px-6 sm:py-6">
    <div className="flex items-center justify-between gap-20">
      <div className="hidden lg:block relative shrink-0 aspect-[3/4] max-w-[18rem] w-full bg-muted-foreground/20 rounded-xl">
        <Image
          src="/placeholder.svg"
          fill
          alt=""
          className="object-cover rounded-xl"
        />

        <div className="absolute top-1/4 right-0 translate-x-1/2 h-12 w-12 bg-primary rounded-full flex items-center justify-center">
          <svg
            width="102"
            height="102"
            viewBox="0 0 102 102"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            className="h-6 w-6"
          >
            <path
              d="M26.0063 19.8917C30.0826 19.8625 33.7081 20.9066 36.8826 23.024C40.057 25.1414 42.5746 28.0279 44.4353 31.6835C46.2959 35.339 47.2423 39.4088 47.2744 43.8927C47.327 51.2301 44.9837 58.4318 40.2444 65.4978C35.4039 72.6664 28.5671 78.5755 19.734 83.2249L2.54766 74.1759C8.33598 71.2808 13.2548 67.9334 17.3041 64.1335C21.2515 60.3344 23.9203 55.8821 25.3105 50.7765C20.5179 50.4031 16.6348 48.9532 13.6612 46.4267C10.5864 44.0028 9.03329 40.5999 9.00188 36.2178C8.97047 31.8358 10.5227 28.0029 13.6584 24.7192C16.693 21.5381 20.809 19.9289 26.0063 19.8917ZM77.0623 19.5257C81.1387 19.4965 84.7641 20.5406 87.9386 22.6581C91.1131 24.7755 93.6306 27.662 95.4913 31.3175C97.3519 34.9731 98.2983 39.0428 98.3304 43.5268C98.383 50.8642 96.0397 58.0659 91.3004 65.1319C86.4599 72.3005 79.6231 78.2095 70.79 82.859L53.6037 73.8099C59.392 70.9149 64.3108 67.5674 68.3601 63.7676C72.3075 59.9685 74.9763 55.5161 76.3665 50.4105C71.5739 50.0372 67.6908 48.5873 64.7172 46.0608C61.6424 43.6369 60.0893 40.2339 60.0579 35.8519C60.0265 31.4698 61.5787 27.6369 64.7145 24.3532C67.7491 21.1722 71.865 19.563 77.0623 19.5257Z"
              className="fill-primary-foreground"
            />
          </svg>
        </div>
      </div>
      <div className="flex flex-col justify-center">
        <div className="flex items-center justify-between gap-1">
          <div className="hidden sm:flex md:hidden items-center gap-4">
            <Avatar className="w-8 h-8 md:w-10 md:h-10">
              <AvatarFallback className="text-xl font-medium bg-primary text-primary-foreground">
                {testimonial.name.charAt(0)}
              </AvatarFallback>
            </Avatar>
            <div>
              <p className="text-lg font-semibold">{testimonial.name}</p>
              <p className="text-sm text-gray-500">{testimonial.designation}</p>
            </div>
          </div>
          <div className="flex items-center gap-1">
            <StarIcon className="w-5 h-5 fill-muted-foreground stroke-muted-foreground" />
            <StarIcon className="w-5 h-5 fill-muted-foreground stroke-muted-foreground" />
            <StarIcon className="w-5 h-5 fill-muted-foreground stroke-muted-foreground" />
            <StarIcon className="w-5 h-5 fill-muted-foreground stroke-muted-foreground" />
            <StarIcon className="w-5 h-5 fill-muted-foreground stroke-muted-foreground" />
          </div>
        </div>
        <p className="mt-6 text-lg sm:text-2xl lg:text-[1.75rem] xl:text-3xl leading-normal lg:!leading-normal font-semibold tracking-tight">
          &quot;{testimonial.testimonial}&quot;
        </p>
        <div className="flex sm:hidden md:flex mt-6 items-center gap-4">
          <Avatar>
            <AvatarFallback className="text-xl font-medium bg-primary text-primary-foreground">
              {testimonial.name.charAt(0)}
            </AvatarFallback>
          </Avatar>
          <div>
            <p className="text-lg font-semibold">{testimonial.name}</p>
            <p className="text-sm text-gray-500">{testimonial.designation}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
);

export default Testimonial;
