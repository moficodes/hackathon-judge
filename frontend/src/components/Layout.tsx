import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';

export default function Layout() {
  return (
    <div className="grid grid-cols-[240px_1fr] min-h-screen">
      <Sidebar />
      <div className="bg-background">
        <main className="max-w-[1280px] mx-auto p-xl">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
