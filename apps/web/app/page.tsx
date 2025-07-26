import FAQ from "@/components/faq";
import Features from "@/components/features";
import Footer from "@/components/footer";
import Hero from "@/components/hero";
import { Navbar } from "@/components/navbar";
import Testimonial from "@/components/testimonial";
import CodeExample from "@/components/code-example";

export default function Home() {
  return (
    <>
      <Navbar />
      <Hero />
      <Features />
      <CodeExample />
      <FAQ />
      <Testimonial />
      <Footer />
    </>
  );
}
