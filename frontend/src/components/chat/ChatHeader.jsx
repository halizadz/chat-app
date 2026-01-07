import React from 'react';
import { MoreVertical, Phone, Video, Users } from 'lucide-react';
import Avatar from '../ui/Avatar';

const ChatHeader = ({ room, members = [], onShowMembers }) => {
  if (!room) {
    return (
      <div className="h-16 border-b border-gray-200 bg-white flex items-center justify-center">
        <p className="text-gray-400">Select a room to start chatting</p>
      </div>
    );
  }

  const onlineMembers = members.filter(m => m.status === 'online').length;

  return (
    <div className="h-16 border-b border-gray-200 bg-white px-6 flex items-center justify-between shadow-sm">
      <div className="flex items-center space-x-4">
        <Avatar 
          alt={room.name} 
          size="md"
          status={room.type === 'private' && onlineMembers > 0 ? 'online' : 'offline'}
        />
        <div>
          <h3 className="font-semibold text-gray-900">{room.name}</h3>
          <p className="text-xs text-gray-500">
            {room.type === 'group' ? (
              <button 
                onClick={onShowMembers}
                className="hover:text-primary-600 transition-colors flex items-center space-x-1"
              >
                <Users className="w-3 h-3" />
                <span>{members.length} members {onlineMembers > 0 && `• ${onlineMembers} online`}</span>
              </button>
            ) : (
              <span>{onlineMembers > 0 ? 'Active now' : 'Offline'}</span>
            )}
          </p>
        </div>
      </div>

      <div className="flex items-center space-x-2">
        <button className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
          <Phone className="w-5 h-5 text-gray-600" />
        </button>
        <button className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
          <Video className="w-5 h-5 text-gray-600" />
        </button>
        <button className="p-2 hover:bg-gray-100 rounded-lg transition-colors">
          <MoreVertical className="w-5 h-5 text-gray-600" />
        </button>
      </div>
    </div>
  );
};

export default ChatHeader;