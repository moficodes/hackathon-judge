import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Header from './Header';

export default function Layout() {
  return (
    <div className="grid grid-cols-[240px_1fr] min-h-screen">
      <Sidebar />
      <div className="flex flex-col min-w-0">
        <Header />
        <main className="flex-1 bg-slate-50 overflow-y-auto">
          <div className="max-w-[1280px] mx-auto p-xl">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}
