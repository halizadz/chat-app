import React from 'react';
import { format } from 'date-fns';
import { File, CheckCheck } from 'lucide-react';
import Avatar from '../ui/Avatar';
import Input from "../ui/Input";
import { useAuthStore } from '../../store/authStore';

const MessageItem = ({ message }) => {
  const currentUser = useAuthStore((state) => state.user);
  const isSent = message.sender_id === currentUser?.id;
  const isSystem = message.type === 'join' || message.type === 'leave';

  if (isSystem) {
    return (
      <div className="flex justify-center my-4">
        <div className="bg-gray-100 text-gray-600 text-xs px-3 py-1.5 rounded-full">
          {message.content}
        </div>
      </div>
    );
  }

  const formatTime = (timestamp) => {
    return format(new Date(timestamp), 'HH:mm');
  };

  const renderFilePreview = () => {
    if (message.type === 'file' && message.file_url) {
      const isImage = /\.(jpg|jpeg|png|gif|webp)$/i.test(message.file_url);
      const fullUrl = `http://localhost:8080${message.file_url}`;
      
      if (isImage) {
        return (
          <div className="mb-2 rounded-lg overflow-hidden">
            <img 
              src={fullUrl}
              alt={message.file_name}
              className="max-w-xs max-h-64 object-cover cursor-pointer hover:opacity-90 transition-opacity"
              onClick={() => window.open(fullUrl, '_blank')}
            />
          </div>
        );
      } else {
        return (
          <a
            href={fullUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center space-x-2 mb-2 p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors">
            <File className="w-5 h-5 text-gray-600" />
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium text-gray-900 truncate">
                {message.file_name}
              </p>
              <p className="text-xs text-gray-500">
                {(message.file_size / 1024).toFixed(2)} KB
              </p>
            </div>
          </a>
        );
      }
    }
    return null;
  };


  return (
    <div className={`flex items-end space-x-2 mb-4 ${isSent ? 'flex-row-reverse space-x-reverse' : ''} animate-slide-up`}>
      {!isSent && (
        <Avatar 
          alt={message.sender?.username || message.username} 
          size="sm"
          className="flex-shrink-0"
        />
      )}
      
      <div className={`flex flex-col ${isSent ? 'items-end' : 'items-start'} max-w-[70%]`}>
        {!isSent && (
          <span className="text-xs text-gray-500 mb-1 ml-2">
            {message.sender?.username || message.username}
          </span>
        )}
        
        <div className={`message-bubble ${isSent ? 'sent' : 'received'}`}>
          {renderFilePreview()}
          
          {message.content && (
            <p className="whitespace-pre-wrap break-words">
              {message.content}
            </p>
          )}
          
          <div className={`flex items-center justify-end space-x-1 mt-1 ${isSent ? 'text-white/70' : 'text-gray-500'}`}>
            <span className="text-xs">
              {formatTime(message.created_at || message.timestamp)}
            </span>
            {isSent && (
              <CheckCheck className="w-3.5 h-3.5" />
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default MessageItem;