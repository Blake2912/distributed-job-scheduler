import { ArrowForward, CalendarMonthOutlined } from "@mui/icons-material";

export function Hero() {
  return (
    <div className="relative overflow-hidden bg-linear-to-b from-blue-50 to-white">
      <div className="mx-auto max-w-7xl px-6 py-24 sm:py-32 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <div className="mb-8 inline-flex items-center gap-2 rounded-full bg-blue-100 px-4 py-2 text-sm text-blue-700">
            <CalendarMonthOutlined className="h-4 w-4" />
            Simple. Powerful. Reliable.
          </div>
          <h1 className="mb-6 text-5xl font-bold tracking-tight text-gray-900 sm:text-7xl">
            Schedule Jobs with Confidence
          </h1>
          <p className="mb-10 text-lg leading-8 text-gray-600">
            Automate your workflows with our intuitive job scheduler. Set it
            once, and let it run reliably—no complexity, no hassle.
          </p>
          <div className="flex flex-col gap-4 sm:flex-row sm:justify-center">
            <button className="inline-flex items-center justify-center gap-2 rounded-lg bg-blue-600 px-8 py-3 text-base font-semibold text-white shadow-lg transition hover:bg-blue-700">
              Get Started Free
              <ArrowForward className="h-5 w-5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
