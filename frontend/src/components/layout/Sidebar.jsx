import React, { useState } from 'react';
import { MessageSquare, LogOut, Settings, User } from 'lucide-react';
import Avatar from '../ui/Avatar';
import SettingsModal from '../modals/SettingsModal';
import { useAuthStore } from '../../store/authStore';
import { useNavigate } from 'react-router-dom';

const Sidebar = () => {
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();
  const [showSettings, setShowSettings] = useState(false);

  const handleLogout = () => {
    logout();
    navigate('/auth');
  };

  return (
    <>
      <div className="w-16 bg-primary-700 flex flex-col items-center py-4 space-y-4">
        {/* Logo */}
        <div className="w-10 h-10 bg-white rounded-lg flex items-center justify-center mb-4">
          <MessageSquare className="w-6 h-6 text-primary-700" />
        </div>

        {/* Navigation */}
        <nav className="flex-1 flex flex-col items-center space-y-2">
          <button
            className="w-10 h-10 rounded-lg bg-primary-600 flex items-center justify-center hover:bg-primary-500 transition-colors"
            title="Chats"
          >
            <MessageSquare className="w-5 h-5 text-white" />
          </button>

          <button
            onClick={() => setShowSettings(true)}
            className="w-10 h-10 rounded-lg hover:bg-primary-600 flex items-center justify-center transition-colors"
            title="Settings"
          >
            <Settings className="w-5 h-5 text-white" />
          </button>
        </nav>

        {/* User Profile */}
        <div className="flex flex-col items-center space-y-2">
          <button
            className="relative"
            title={user?.username}
          >
            <Avatar 
              alt={user?.username || 'User'} 
              size="md"
              status="online"
            />
          </button>

          <button
            onClick={handleLogout}
            className="w-10 h-10 rounded-lg hover:bg-primary-600 flex items-center justify-center transition-colors"
            title="Logout"
          >
            <LogOut className="w-5 h-5 text-white" />
          </button>
        </div>
      </div>

      <SettingsModal
        isOpen={showSettings}
        onClose={() => setShowSettings(false)}
      />
    </>
  );
};

export default Sidebar;