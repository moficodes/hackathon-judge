import { Link } from 'react-router-dom';

export default function About() {
  return (
    <div>
      <h2>About</h2>
      <p>This is a tool for judging hackathons.</p>
      <Link to="/">Go back Home</Link>
    </div>
  );
}
