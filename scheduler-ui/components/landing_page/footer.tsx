import { CalendarMonth } from "@mui/icons-material";

export function Footer() {
  return (
    <footer className="bg-gray-900">
      <div className="mx-auto max-w-7xl px-6 py-12 lg:px-8">
        <div className="grid grid-cols-2 gap-12 md:grid-cols-2">
          <div className="col-span-2 md:col-span-1">
            <div className="mb-4 flex items-center gap-2">
              <CalendarMonth className="h-6 w-6 text-blue-500" />
              <span className="text-lg font-bold text-white">JobScheduler</span>
            </div>
            <p className="text-sm text-gray-400">
              Simple and reliable job scheduling for modern teams.
            </p>
          </div>
          <div>
            <h3 className="mb-4 text-sm font-semibold text-white">Product</h3>
            <ul className="space-y-2 text-sm text-gray-400">
              <li>
                <a
                  href="https://github.com/Blake2912/distributed-job-scheduler"
                  className="hover:text-white"
                >
                  GitHub
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/Blake2912/distributed-job-scheduler/releases"
                  className="hover:text-white"
                >
                  Changelog
                </a>
              </li>
            </ul>
          </div>
        </div>
        <div className="mt-12 border-t border-gray-800 pt-8 text-center text-sm text-gray-400">
          <p>© 2026 JobScheduler. All rights reserved.</p>
        </div>
      </div>
    </footer>
  );
}
