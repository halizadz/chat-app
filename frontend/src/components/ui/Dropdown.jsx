import React, { useState, useRef, useEffect } from 'react';

const Dropdown = ({ trigger, children, align = 'left' }) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const alignmentClasses = {
    left: 'left-0',
    right: 'right-0',
    center: 'left-1/2 -translate-x-1/2',
  };

  return (
    <div className="relative" ref={dropdownRef}>
      <div onClick={() => setIsOpen(!isOpen)}>
        {trigger}
      </div>

      {isOpen && (
        <div className={`
          absolute top-full mt-2 ${alignmentClasses[align]}
          bg-white rounded-lg shadow-lg border border-gray-200
          min-w-[200px] py-2 z-50
          animate-slide-down
        `}>
          {children}
        </div>
      )}
    </div>
  );
};

export const DropdownItem = ({ onClick, icon: Icon, children, danger = false }) => {
  return (
    <button
      onClick={onClick}
      className={`
        w-full px-4 py-2 text-left flex items-center space-x-3
        transition-colors
        ${danger 
          ? 'text-red-600 hover:bg-red-50' 
          : 'text-gray-700 hover:bg-gray-50'
        }
      `}
    >
      {Icon && <Icon className="w-4 h-4" />}
      <span className="text-sm">{children}</span>
    </button>
  );
};

export default Dropdown;