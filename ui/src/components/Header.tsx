import React from 'react';
import { Menu } from 'lucide-react'; // Using lucide-react for icons

const Header: React.FC = () => {
  return (
    <header className="bg-background text-white max-h-[10vh] p-4 flex items-center justify-between shadow-md">
      {/* Left side: Logo placeholder */}
      <div className="flex items-center space-x-2">
        {/* Placeholder for logo - replace with actual <img> or SVG */}
        <div className="w-8 h-8 bg-red-400 rounded flex items-center justify-center text-sm font-bold">
          L
        </div>
        {/* Application Title */}
        <p className="text-3xl font-bold">Oversee</p>
      </div>

      {/* Right side: Hamburger Menu Button */}
      <button
        className="p-2 rounded hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-300"
        aria-label="Open menu"
      >
        <Menu size={24} />
      </button>
    </header>
  );
};

export default Header;

