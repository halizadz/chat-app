import React from 'react';
import { format } from 'date-fns';
import Avatar from '../ui/Avatar';

const RoomItem = ({ room, isActive, onClick, lastMessage }) => {
  const formatTime = (timestamp) => {
    const date = new Date(timestamp);
    const today = new Date();
    
    if (date.toDateString() === today.toDateString()) {
      return format(date, 'HH:mm');
    }
    
    return format(date, 'MMM d');
  };

  return (
    <button
      onClick={onClick}
      className={`
        w-full px-4 py-3 flex items-center space-x-3 hover:bg-gray-100 transition-all duration-200
        ${isActive ? 'bg-primary-50 border-l-4 border-primary-600' : 'border-l-4 border-transparent'}
      `}
    >
      <Avatar 
        alt={room.name} 
        size="md"
        status={room.type === 'private' ? 'online' : undefined}
      />
      
      <div className="flex-1 min-w-0 text-left">
        <div className="flex items-center justify-between mb-1">
          <h4 className={`font-semibold truncate ${isActive ? 'text-primary-900' : 'text-gray-900'}`}>
            {room.name}
          </h4>
          {lastMessage && (
            <span className="text-xs text-gray-500 ml-2 flex-shrink-0">
              {formatTime(lastMessage.created_at)}
            </span>
          )}
        </div>
        
        <div className="flex items-center justify-between">
          <p className="text-sm text-gray-600 truncate">
            {lastMessage ? (
              <>
                {lastMessage.type === 'file' ? '📎 File' : lastMessage.content}
              </>
            ) : (
              <span className="text-gray-400">No messages yet</span>
            )}
          </p>
          {lastMessage?.unread_count > 0 && (
            <span className="ml-2 px-2 py-0.5 bg-primary-600 text-white text-xs font-semibold rounded-full flex-shrink-0">
              {lastMessage.unread_count > 99 ? '99+' : lastMessage.unread_count}
            </span>
          )}
        </div>
      </div>
    </button>
  );
};

export default RoomItem;