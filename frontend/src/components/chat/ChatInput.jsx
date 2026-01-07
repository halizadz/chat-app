import React, { useState, useRef, useEffect } from 'react';
import { Send, Paperclip, Smile, X } from 'lucide-react';
import Button from '../ui/Button';
import { fileAPI } from '../../services/api';
import toast from 'react-hot-toast';
import { MESSAGE_TYPES } from '../../utils/constants';

const ChatInput = ({ onSendMessage, onTyping }) => {
  const [message, setMessage] = useState('');
  const [isTyping, setIsTyping] = useState(false);
  const [selectedFile, setSelectedFile] = useState(null);
  const [isUploading, setIsUploading] = useState(false);
  const fileInputRef = useRef(null);
  const typingTimeoutRef = useRef(null);

  useEffect(() => {
    return () => {
      if (typingTimeoutRef.current) {
        clearTimeout(typingTimeoutRef.current);
      }
    };
  }, []);

  const handleTyping = () => {
    if (!isTyping) {
      setIsTyping(true);
      onTyping(true);
    }

    if (typingTimeoutRef.current) {
      clearTimeout(typingTimeoutRef.current);
    }

    typingTimeoutRef.current = setTimeout(() => {
      setIsTyping(false);
      onTyping(false);
    }, 2000);
  };

  const handleChange = (e) => {
    setMessage(e.target.value);
    handleTyping();
  };

  const handleFileSelect = (e) => {
    const file = e.target.files?.[0];
    if (file) {
      if (file.size > 10 * 1024 * 1024) {
        toast.error('File size must be less than 10MB');
        return;
      }
      setSelectedFile(file);
    }
  };

  const handleRemoveFile = () => {
    setSelectedFile(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!message.trim() && !selectedFile) return;

    try {
      if (selectedFile) {
        setIsUploading(true);
        const response = await fileAPI.uploadFile(selectedFile);
        const { url, file_name, file_size } = response.data;

        onSendMessage({
          type: MESSAGE_TYPES.FILE,
          content: message.trim() || 'Sent a file',
          file_url: url,
          file_name: file_name,
          file_size: file_size,
        });

        handleRemoveFile();
      } else {
        onSendMessage({
          type: MESSAGE_TYPES.TEXT,
          content: message.trim(),
        });
      }

      setMessage('');
      setIsTyping(false);
      onTyping(false);
    } catch (error) {
      console.error('Error sending message:', error);
      toast.error('Failed to send message');
    } finally {
      setIsUploading(false);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  return (
    <div className="border-t border-gray-200 bg-white px-6 py-4">
      {selectedFile && (
        <div className="mb-3 animate-slide-down">
          <div className="inline-flex items-center space-x-2 bg-primary-50 text-primary-700 px-3 py-2 rounded-lg">
            <Paperclip className="w-4 h-4" />
            <span className="text-sm font-medium truncate max-w-xs">
              {selectedFile.name}
            </span>
            <span className="text-xs text-primary-600">
              ({(selectedFile.size / 1024).toFixed(2)} KB)
            </span>
            <button
              onClick={handleRemoveFile}
              className="ml-2 hover:bg-primary-100 rounded-full p-1 transition-colors"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="flex items-end space-x-3">
        <div className="flex space-x-2">
          <input
            ref={fileInputRef}
            type="file"
            onChange={handleFileSelect}
            className="hidden"
            accept="image/*,.pdf,.doc,.docx,.txt"
          />
          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            className="p-2.5 text-gray-500 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors"
            title="Attach file"
          >
            <Paperclip className="w-5 h-5" />
          </button>
          
          <button
            type="button"
            className="p-2.5 text-gray-500 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors"
            title="Add emoji"
          >
            <Smile className="w-5 h-5" />
          </button>
        </div>

        <div className="flex-1 relative">
          <textarea
            value={message}
            onChange={handleChange}
            onKeyPress={handleKeyPress}
            placeholder="Type a message..."
            rows={1}
            className="w-full px-4 py-2.5 pr-12 rounded-lg border border-gray-300 focus:border-primary-500 focus:ring-2 focus:ring-primary-500 focus:outline-none resize-none custom-scrollbar"
            style={{ maxHeight: '120px' }}
          />
        </div>

        <Button
          type="submit"
          variant="primary"
          className="px-4 py-2.5"
          disabled={(!message.trim() && !selectedFile) || isUploading}
          isLoading={isUploading}
        >
          <Send className="w-5 h-5" />
        </Button>
      </form>
    </div>
  );
};

export default ChatInput;