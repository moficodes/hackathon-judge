interface User {
  name: string;
  initials: string;
  role: string;
}

interface HeaderProps {
  user?: User;
  status?: string;
}

const defaultUser: User = {
  name: "Judge User",
  initials: "JD",
  role: "Technical Track"
};

const defaultStatus = "Active Judge";

export default function Header({ user = defaultUser, status = defaultStatus }: HeaderProps) {
  return (
    <header className="h-[64px] border-b border-slate-200 bg-white flex items-center justify-between px-lg z-20">
      {/* Left: Search */}
      <div className="flex items-center gap-sm bg-slate-50 border border-slate-300 px-md py-xs rounded-md w-[320px] group focus-within:ring-2 focus-within:ring-secondary/20 focus-within:border-secondary transition-all">
        <span className="text-slate-400">🔍</span>
        <input 
          type="text" 
          placeholder="Search projects, teams, or hackathons..." 
          aria-label="Search"
          className="bg-transparent border-none outline-none text-sm w-full text-slate-900 placeholder:text-slate-500"
        />
        <kbd className="hidden sm:inline-flex items-center gap-1 px-1.5 font-sans text-[10px] font-medium text-slate-400 bg-white border border-slate-300 rounded shadow-xs">
          ⌘K
        </kbd>
      </div>

      {/* Right: Utilities & User */}
      <div className="flex items-center gap-xl">
        {/* Status Badge */}
        <div className="flex items-center gap-xs bg-secondary-container/10 px-sm py-[2px] rounded-full border border-secondary-container/20">
          <div className="w-2 h-2 rounded-full bg-secondary"></div>
          <span className="text-[10px] font-bold text-secondary tracking-widest uppercase">
            {status}
          </span>
        </div>

        {/* Icons */}
        <div className="flex items-center gap-md text-slate-600">
          <button aria-label="Notifications" className="hover:text-secondary transition-colors cursor-pointer text-xl">🔔</button>
          <button aria-label="Settings" className="hover:text-secondary transition-colors cursor-pointer text-xl">⚙️</button>
        </div>

        {/* User Profile */}
        <div className="flex items-center gap-md border-l border-slate-200 pl-xl cursor-pointer group">
          <div className="w-8 h-8 rounded-full bg-slate-900 text-white flex items-center justify-center text-xs font-bold ring-2 ring-transparent group-hover:ring-secondary/20 transition-all">
            {user.initials}
          </div>
          <div className="hidden md:flex flex-col">
            <span className="text-sm font-semibold text-slate-900 leading-tight">{user.name}</span>
            <span className="text-[11px] text-slate-500 leading-tight">{user.role}</span>
          </div>
          <span className="text-[10px] text-slate-400 group-hover:text-slate-600 transition-colors">▼</span>
        </div>
      </div>
    </header>
  );
}
