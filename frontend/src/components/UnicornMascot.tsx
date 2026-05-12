import { useState } from 'react';
import calmUnicorn from '../assets/unicorn_calm.svg';
import neighUnicorn from '../assets/unicorn_neigh.svg';

export default function UnicornMascot() {
  const [isNeighing, setIsNeighing] = useState(false);

  return (
    <div 
      className="flex items-center gap-3 cursor-pointer group"
      onClick={() => setIsNeighing(!isNeighing)}
      title="Click me!"
    >
      <img 
        src={isNeighing ? neighUnicorn : calmUnicorn} 
        alt="Friendly Unicorn Mascot" 
        className="w-12 h-12 transition-transform group-hover:scale-110"
      />
      <span className="font-bold text-lg select-none">
        Hackathon Judge
      </span>
    </div>
  );
}
