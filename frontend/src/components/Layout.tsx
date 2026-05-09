import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';

export default function Layout() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-[240px_1fr] h-screen overflow-hidden">
      <Sidebar />
      <div className="bg-background overflow-y-auto">
        <main className="max-w-7xl mx-auto p-xl">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
