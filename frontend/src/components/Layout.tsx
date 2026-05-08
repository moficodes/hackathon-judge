import { Outlet, Link } from 'react-router-dom';

export default function Layout() {
  return (
    <div className="grid grid-cols-[240px_1fr] min-h-screen">
      {/* Sidebar Navigation */}
      <aside className="bg-primary text-on-primary border-r border-outline-variant/20 shadow-sm z-10 flex flex-col">
        <div className="p-lg font-bold text-lg border-b border-outline-variant/20">
          Hackathon Judge
        </div>
        <nav className="flex-1 p-md">
          <ul className="flex flex-col gap-sm">
            <li>
              <Link 
                to="/" 
                className="block px-md py-sm rounded hover:bg-white/10 transition-colors text-sm font-semibold tracking-wide"
              >
                Home
              </Link>
            </li>
            <li>
              <Link 
                to="/dashboard" 
                className="block px-md py-sm rounded hover:bg-white/10 transition-colors text-sm font-semibold tracking-wide"
              >
                Dashboard
              </Link>
            </li>
            <li>
              <Link 
                to="/about" 
                className="block px-md py-sm rounded hover:bg-white/10 transition-colors text-sm font-semibold tracking-wide text-on-primary/70"
              >
                About
              </Link>
            </li>
          </ul>
        </nav>
      </aside>

      {/* Main Content Area */}
      <div className="bg-background">
        <main className="max-w-[1280px] mx-auto p-xl">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
