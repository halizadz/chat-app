import React from 'react';

const TypingIndicator = ({ typingUsers }) => {
  if (typingUsers.size === 0) return null;

  const users = Array.from(typingUsers).map(u => JSON.parse(u));
  const names = users.map(u => u.username).join(', ');

  return (
    <div className="px-6 py-2 animate-slide-up">
      <div className="flex items-center space-x-2 text-gray-500">
        <div className="flex space-x-1">
          <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
          <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
          <div className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
        </div>
        <span className="text-sm">
          {names} {users.length === 1 ? 'is' : 'are'} typing...
        </span>
      </div>
    </div>
  );
};

export default TypingIndicator;