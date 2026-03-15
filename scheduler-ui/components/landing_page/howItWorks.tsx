"use-client";
import {
  CalendarMonthOutlined,
  CheckCircleOutlineOutlined,
  PlayArrowRounded,
  PlusOne,
} from "@mui/icons-material";

const steps = [
  {
    icon: PlusOne,
    title: "Create a Job",
    description: "Define your task with a simple interface or API call.",
    step: "01",
  },
  {
    icon: CalendarMonthOutlined,
    title: "Set the Schedule",
    description: "Choose when and how often your job should run.",
    step: "02",
  },
  {
    icon: PlayArrowRounded,
    title: "Activate & Monitor",
    description: "Turn it on and watch it execute automatically.",
    step: "03",
  },
  {
    icon: CheckCircleOutlineOutlined,
    title: "Get Results",
    description: "Receive notifications and view detailed execution logs.",
    step: "04",
  },
];

export function HowItWorks() {
  return (
    <div className="bg-gray-50 py-24 sm:py-32">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="mb-4 text-base font-semibold text-blue-600">
            How it works
          </h2>
          <p className="text-4xl font-bold tracking-tight text-gray-900">
            Get started in minutes
          </p>
          <p className="mt-4 text-lg text-gray-600">
            Four simple steps to automate your workflows and save hours every
            week.
          </p>
        </div>
        <div className="mx-auto mt-16 max-w-5xl">
          <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-4">
            {steps.map((step) => (
              <div key={step.title} className="relative">
                <div className="flex flex-col items-center text-center">
                  <div className="mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-blue-600 text-2xl font-bold text-white">
                    {step.step}
                  </div>
                  <div className="mb-3 inline-flex h-12 w-12 items-center justify-center rounded-lg bg-white shadow-md">
                    <step.icon className="h-6 w-6 text-blue-600" />
                  </div>
                  <h3 className="mb-2 text-lg font-semibold text-gray-900">
                    {step.title}
                  </h3>
                  <p className="text-sm text-gray-600">{step.description}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
