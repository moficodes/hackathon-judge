import { NavLink } from 'react-router-dom';

const navItems = [
  { name: 'Home', path: '/' },
  { name: 'Dashboard', path: '/dashboard' },
  { name: 'About', path: '/about' },
];

export default function Sidebar() {
  return (
    <aside className="bg-primary text-white border-r border-white/10 shadow-sm z-10 flex flex-col w-[240px]">
      <div className="p-lg font-bold text-lg border-b border-white/10">
        Hackathon Judge
      </div>
      <nav className="flex-1 p-md">
        <ul className="flex flex-col gap-sm">
          {navItems.map((item) => (
            <li key={item.path}>
              <NavLink
                to={item.path}
                className={({ isActive }) =>
                  `block px-md py-sm rounded transition-colors text-sm font-semibold tracking-wide ${
                    isActive ? 'bg-white/20 text-white' : 'text-white/70 hover:bg-white/10 hover:text-white'
                  }`
                }
              >
                {item.name}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
}
