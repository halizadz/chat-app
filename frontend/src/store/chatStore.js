import { create } from 'zustand';
import { roomAPI } from '../services/api';
import toast from 'react-hot-toast';

export const useChatStore = create((set) => ({
  rooms: [], // ← Pastikan ini array kosong, bukan null
  currentRoom: null,
  messages: [],
  members: [],
  isLoading: false,
  typingUsers: new Set(),

  fetchRooms: async () => {
    set({ isLoading: true });
    try {
      const response = await roomAPI.getUserRooms();
      set({ rooms: response.data || [], isLoading: false }); // ← Fallback ke array kosong
    } catch (error) {
      console.error('Error fetching rooms:', error);
      set({ isLoading: false, rooms: [] }); // ← Set array kosong jika error
    }
  },

  setCurrentRoom: async (room) => {
    set({ currentRoom: room, messages: [], isLoading: true });
    try {
      const [messagesRes, membersRes] = await Promise.all([
        roomAPI.getRoomMessages(room.id),
        roomAPI.getRoomMembers(room.id),
      ]);
      set({ 
        messages: messagesRes.data || [], 
        members: membersRes.data || [],
        isLoading: false 
      });
    } catch (error) {
      console.error('Error loading room data:', error);
      set({ isLoading: false });
    }
  },

  createRoom: async (roomData) => {
    try {
      const response = await roomAPI.createRoom(roomData);
      set((state) => ({ rooms: [...state.rooms, response.data] }));
      toast.success('Room created successfully!');
      return response.data;
    } catch {
      toast.error('Failed to create room');
      return null;
    }
  },

  addMessage: (message) => {
    set((state) => ({
      messages: [...state.messages, message],
    }));
  },

  addTypingUser: (userId, username) => {
    set((state) => {
      const newTypingUsers = new Set(state.typingUsers);
      newTypingUsers.add(JSON.stringify({ userId, username }));
      return { typingUsers: newTypingUsers };
    });
  },

  removeTypingUser: (userId) => {
    set((state) => {
      const newTypingUsers = new Set(
        Array.from(state.typingUsers).filter(
          (user) => JSON.parse(user).userId !== userId
        )
      );
      return { typingUsers: newTypingUsers };
    });
  },

  // Cleanup function
  reset: () => {
    set({
      rooms: [],
      currentRoom: null,
      messages: [],
      members: [],
      isLoading: false,
      typingUsers: new Set(),
    });
  },
}));