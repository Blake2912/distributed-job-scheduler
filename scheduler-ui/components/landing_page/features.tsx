import {
  BarChart,
  BoltOutlined,
  LockClock,
  Settings,
  ShieldOutlined,
} from "@mui/icons-material";

const features = [
  {
    icon: LockClock,
    title: "Flexible Scheduling",
    description:
      "Set up cron jobs, recurring tasks, or one-time schedules with an easy-to-use interface.",
  },
  {
    icon: BoltOutlined,
    title: "Lightning Fast",
    description:
      "Optimized performance ensures your jobs execute exactly when they should, every time.",
  },
  {
    icon: ShieldOutlined,
    title: "Reliable",
    description: "Automatic retries and failure handling built-in.",
  },
  {
    icon: BarChart,
    title: "Real-time Monitoring",
    description:
      "Track job status, execution history, and performance metrics at a glance.",
  },
  {
    icon: Settings,
    title: "Easy Integration",
    description:
      "Connect with your existing tools and workflows via our REST API or webhooks.",
  },
];

export function Features() {
  return (
    <div className="bg-white py-24 sm:py-32">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="mb-4 text-base font-semibold text-blue-600">
            Everything you need
          </h2>
          <p className="text-4xl font-bold tracking-tight text-gray-900">
            Powerful features for modern teams
          </p>
          <p className="mt-4 text-lg text-gray-600">
            All the tools you need to automate and manage your scheduled tasks
            efficiently.
          </p>
        </div>
        <div className="mx-auto mt-16 max-w-7xl">
          <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
            {features.map((feature) => (
              <div
                key={feature.title}
                className="relative rounded-2xl border border-gray-200 bg-white p-8 shadow-sm transition hover:shadow-md"
              >
                <div className="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-lg bg-blue-100">
                  <feature.icon className="h-6 w-6 text-blue-600" />
                </div>
                <h3 className="mb-2 text-xl font-semibold text-gray-900">
                  {feature.title}
                </h3>
                <p className="text-gray-600">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
