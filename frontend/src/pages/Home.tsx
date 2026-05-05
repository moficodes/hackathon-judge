import { Link } from 'react-router-dom';

export default function Home() {
  return (
    <div className="p-4">
      <h2 className="text-2xl font-bold text-blue-500 mb-4">Home</h2>
      <p className="mb-4">Welcome to Hackathon Judge.</p>
      <Link to="/about" className="text-blue-600 hover:underline">Go to About</Link>
    </div>
  );
}
