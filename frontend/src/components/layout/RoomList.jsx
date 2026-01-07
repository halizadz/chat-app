import React, { useState } from 'react';
import { Plus, Search, MessageSquare, Users, MessageCircle } from 'lucide-react';
import RoomItem from './RoomItem';
import Button from '../ui/Button';
import Input from '../ui/Input';
import CreateRoomModal from '../modals/CreateRoomModal';
import NewChatModal from '../modals/NewChatModal';
import { useChatStore } from '../../store/chatStore';

const RoomList = ({ onSelectRoom, currentRoom }) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showNewChatModal, setShowNewChatModal] = useState(false);
  const [showMenu, setShowMenu] = useState(false);
  const { rooms, setCurrentRoom } = useChatStore();

  // Fix: Pastikan rooms adalah array
  const roomsArray = Array.isArray(rooms) ? rooms : [];

  const filteredRooms = roomsArray.filter(room =>
    room.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const handleChatCreated = (room) => {
    setCurrentRoom(room);
  };

  return (
    <>
      <div className="w-80 bg-white border-r border-gray-200 flex flex-col">
        {/* Header */}
        <div className="p-4 border-b border-gray-200">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-bold text-gray-900">Messages</h2>
            
            {/* Dropdown Button */}
            <div className="relative">
              <Button
                variant="primary"
                size="sm"
                icon={Plus}
                onClick={() => setShowMenu(!showMenu)}
              >
                New
              </Button>

              {showMenu && (
                <div className="absolute right-0 mt-2 w-56 bg-white rounded-lg shadow-lg border border-gray-200 py-2 z-10 animate-slide-down">
                  <button
                    onClick={() => {
                      setShowNewChatModal(true);
                      setShowMenu(false);
                    }}
                    className="w-full px-4 py-2 text-left flex items-center space-x-3 hover:bg-gray-50 transition-colors"
                  >
                    <MessageCircle className="w-4 h-4 text-gray-600" />
                    <div>
                      <p className="text-sm font-medium text-gray-900">New Chat</p>
                      <p className="text-xs text-gray-500">Start a private conversation</p>
                    </div>
                  </button>

                  <button
                    onClick={() => {
                      setShowCreateModal(true);
                      setShowMenu(false);
                    }}
                    className="w-full px-4 py-2 text-left flex items-center space-x-3 hover:bg-gray-50 transition-colors"
                  >
                    <Users className="w-4 h-4 text-gray-600" />
                    <div>
                      <p className="text-sm font-medium text-gray-900">New Group</p>
                      <p className="text-xs text-gray-500">Create a group chat</p>
                    </div>
                  </button>
                </div>
              )}
            </div>
          </div>

          {/* Search */}
          <Input
            type="text"
            placeholder="Search conversations..."
            icon={Search}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>

        {/* Room List */}
        <div className="flex-1 overflow-y-auto custom-scrollbar">
          {filteredRooms.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-center px-4">
              <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
                <MessageSquare className="w-8 h-8 text-gray-400" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">
                {searchQuery ? 'No results found' : 'No conversations yet'}
              </h3>
              <p className="text-sm text-gray-500 mb-4">
                {searchQuery 
                  ? 'Try searching with different keywords'
                  : 'Start a new conversation to get started'
                }
              </p>
              {!searchQuery && (
                <Button
                  variant="outline"
                  size="sm"
                  icon={Plus}
                  onClick={() => setShowMenu(true)}
                >
                  Start Chat
                </Button>
              )}
            </div>
          ) : (
            filteredRooms.map((room) => (
              <RoomItem
                key={room.id}
                room={room}
                isActive={currentRoom?.id === room.id}
                onClick={() => onSelectRoom(room)}
              />
            ))
          )}
        </div>
      </div>

      {/* Close menu when clicking outside */}
      {showMenu && (
        <div 
          className="fixed inset-0 z-0" 
          onClick={() => setShowMenu(false)}
        />
      )}

      <NewChatModal
        isOpen={showNewChatModal}
        onClose={() => setShowNewChatModal(false)}
        onChatCreated={handleChatCreated}
      />

      <CreateRoomModal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
      />
    </>
  );
};

export default RoomList;