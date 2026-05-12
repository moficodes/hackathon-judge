import { NavLink } from 'react-router-dom';
import { Home, LayoutDashboard, Info } from 'lucide-react';
import UnicornMascot from './UnicornMascot';

const navItems = [
  { name: 'Home', path: '/', icon: Home },
  { name: 'Dashboard', path: '/dashboard', icon: LayoutDashboard },
  { name: 'About', path: '/about', icon: Info },
];

export default function Sidebar() {
  return (
    <aside className="hidden md:flex bg-primary text-white border-r border-white/10 shadow-sm z-10 flex-col w-[240px]">
      <div className="p-lg border-b border-white/10">
        <UnicornMascot />
      </div>
      <nav className="flex-1 p-md">
        <ul className="flex flex-col gap-sm">
          {navItems.map((item) => (
            <li key={item.path}>
              <NavLink
                to={item.path}
                className={({ isActive }) =>
                  `flex items-center gap-md px-md py-sm rounded transition-colors text-sm font-semibold tracking-wide ${
                    isActive ? 'bg-white/20 text-white' : 'text-white/70 hover:bg-white/10 hover:text-white'
                  }`
                }
              >
                <item.icon className="w-4 h-4" />
                {item.name}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
}
