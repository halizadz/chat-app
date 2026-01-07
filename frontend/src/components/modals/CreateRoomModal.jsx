import React, { useState } from 'react';
import { X, Globe, Lock } from 'lucide-react';
import Button from '../ui/Button';
import Input from '../ui/Input';
import { useChatStore } from '../../store/chatStore';
import { ROOM_TYPES } from '../../utils/constants';

const CreateRoomModal = ({ isOpen, onClose }) => {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    type: ROOM_TYPES.GROUP,
  });
  const [isLoading, setIsLoading] = useState(false);
  const { createRoom } = useChatStore();

  if (!isOpen) return null;

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);

    const room = await createRoom(formData);
    
    setIsLoading(false);
    
    if (room) {
      onClose();
      setFormData({ name: '', description: '', type: ROOM_TYPES.GROUP });
    }
  };

  const handleChange = (e) => {
    setFormData(prev => ({
      ...prev,
      [e.target.name]: e.target.value
    }));
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 animate-fade-in">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md mx-4 animate-slide-up">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <h2 className="text-xl font-bold text-gray-900">Create New Room</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X className="w-5 h-5 text-gray-500" />
          </button>
        </div>

        {/* Body */}
        <form onSubmit={handleSubmit} className="p-6 space-y-5">
          <Input
            label="Room Name"
            name="name"
            placeholder="General Chat"
            value={formData.name}
            onChange={handleChange}
            required
          />

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1.5">
              Description (Optional)
            </label>
            <textarea
              name="description"
              placeholder="What's this room about?"
              value={formData.description}
              onChange={handleChange}
              rows={3}
              className="w-full px-4 py-2.5 rounded-lg border border-gray-300 focus:border-primary-500 focus:ring-2 focus:ring-primary-500 focus:outline-none resize-none"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-3">
              Room Type
            </label>
            <div className="space-y-2">
              <label className="flex items-center p-3 border-2 border-gray-200 rounded-lg cursor-pointer hover:bg-gray-50 transition-colors">
                <input
                  type="radio"
                  name="type"
                  value={ROOM_TYPES.GROUP}
                  checked={formData.type === ROOM_TYPES.GROUP}
                  onChange={handleChange}
                  className="w-4 h-4 text-primary-600 focus:ring-primary-500"
                />
                <div className="ml-3 flex items-center">
                  <Globe className="w-5 h-5 text-gray-600 mr-2" />
                  <div>
                    <p className="font-medium text-gray-900">Group Room</p>
                    <p className="text-sm text-gray-500">Multiple members can join</p>
                  </div>
                </div>
              </label>

              <label className="flex items-center p-3 border-2 border-gray-200 rounded-lg cursor-pointer hover:bg-gray-50 transition-colors">
                <input
                  type="radio"
                  name="type"
                  value={ROOM_TYPES.PRIVATE}
                  checked={formData.type === ROOM_TYPES.PRIVATE}
                  onChange={handleChange}
                  className="w-4 h-4 text-primary-600 focus:ring-primary-500"
                />
                <div className="ml-3 flex items-center">
                  <Lock className="w-5 h-5 text-gray-600 mr-2" />
                  <div>
                    <p className="font-medium text-gray-900">Private Room</p>
                    <p className="text-sm text-gray-500">One-on-one conversation</p>
                  </div>
                </div>
              </label>
            </div>
          </div>

          {/* Footer */}
          <div className="flex space-x-3 pt-4">
            <Button
              type="button"
              variant="ghost"
              onClick={onClose}
              className="flex-1"
            >
              Cancel
            </Button>
            <Button
              type="submit"
              variant="primary"
              className="flex-1"
              isLoading={isLoading}
            >
              Create Room
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default CreateRoomModal;
