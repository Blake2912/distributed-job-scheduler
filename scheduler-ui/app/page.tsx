import { Features } from "@/components/landing_page/features";
import { Footer } from "@/components/landing_page/footer";
import { Hero } from "@/components/landing_page/hero";
import { HowItWorks } from "@/components/landing_page/howItWorks";

export default function Home() {
  return (
    <div className="min-h-screen">
      <Hero />
      <Features />
      <HowItWorks />
      <Footer />
    </div>
  );
}
