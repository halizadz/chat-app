import React from 'react';

const Avatar = ({ 
  src, 
  alt, 
  size = 'md', 
  status,
  className = '' 
}) => {
  const sizes = {
    sm: 'w-8 h-8 text-xs',
    md: 'w-10 h-10 text-sm',
    lg: 'w-12 h-12 text-base',
    xl: 'w-16 h-16 text-lg',
  };

  const statusColors = {
    online: 'bg-green-500',
    offline: 'bg-gray-400',
    away: 'bg-yellow-500',
  };

  const getInitials = (name) => {
    if (!name) return '?';
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  };

  return (
    <div className={`relative inline-block ${className}`}>
      <div className={`
        ${sizes[size]} 
        rounded-full overflow-hidden
        bg-gradient-to-br from-primary-400 to-primary-600
        flex items-center justify-center
        text-white font-semibold
        ring-2 ring-white
      `}>
        {src ? (
          <img 
            src={src} 
            alt={alt} 
            className="w-full h-full object-cover"
          />
        ) : (
          <span>{getInitials(alt)}</span>
        )}
      </div>
      
      {status && (
        <span className={`
          absolute bottom-0 right-0 block rounded-full ring-2 ring-white
          ${statusColors[status]}
          ${size === 'sm' ? 'w-2 h-2' : size === 'md' ? 'w-2.5 h-2.5' : 'w-3 h-3'}
        `} />
      )}
    </div>
  );
};

export default Avatar;