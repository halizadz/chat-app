import axios from "axios";
import { API_BASE_URL } from "../utils/constants";

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      window.location.href = "/auth";
    }
    return Promise.reject(error);
  }
);

// Auth APIs
export const authAPI = {
  register: (data) => api.post("/auth/register", data),
  login: (data) => api.post("/auth/login", data),
};

// User APIs
export const userAPI = {
  getAllUsers: () => api.get("/users"),
  searchUsers: (query) => api.get(`/users?search=${encodeURIComponent(query)}`),
  searchUsersByGmail: (gmailQuery) =>
    api.get(`/users?gmail=${encodeURIComponent(gmailQuery)}`),
  getUserProfile: () => api.get("/users/me"),
  updateUserProfile: (data) => api.put("/users/me", data),
};

// Room APIs
export const roomAPI = {
  getUserRooms: () => api.get("/rooms"),
  createRoom: (data) => api.post("/rooms", data),
  createPrivateRoom: (userId) =>
    api.post("/rooms/private", { user_id: userId }),
  getRoom: (roomId) => api.get(`/rooms/${roomId}`),
  updateRoom: (roomId, data) => api.put(`/rooms/${roomId}`, data),
  deleteRoom: (roomId) => api.delete(`/rooms/${roomId}`),
  getRoomMessages: (roomId, limit = 50, offset = 0) =>
    api.get(`/rooms/${roomId}/messages?limit=${limit}&offset=${offset}`),
  searchMessages: (roomId, query, limit = 50, offset = 0) =>
    api.get(
      `/rooms/${roomId}/messages/search?q=${encodeURIComponent(
        query
      )}&limit=${limit}&offset=${offset}`
    ),
  markRoomAsRead: (roomId) => api.post(`/rooms/${roomId}/read`),
  getRoomMembers: (roomId) => api.get(`/rooms/${roomId}/members`),
  addMember: (roomId, userId) =>
    api.post(`/rooms/${roomId}/members`, { user_id: userId }),
  removeMember: (roomId, userId) =>
    api.delete(`/rooms/${roomId}/members/${userId}`),
  leaveRoom: (roomId) => api.post(`/rooms/${roomId}/leave`),
};

// Message APIs
export const messageAPI = {
  updateMessage: (messageId, content) =>
    api.put(`/messages/${messageId}`, { content }),
  deleteMessage: (messageId) => api.delete(`/messages/${messageId}`),
};

// File APIs
export const fileAPI = {
  uploadFile: (file) => {
    const formData = new FormData();
    formData.append("file", file);
    return api.post("/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
  },
};

export default api;
