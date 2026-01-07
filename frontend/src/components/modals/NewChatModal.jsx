import React, { useState, useEffect, useRef } from "react";
import { X, Search, MessageCircle, Loader2 } from "lucide-react";
import Input from "../ui/Input";
import Avatar from "../ui/Avatar";
import { roomAPI, userAPI } from "../../services/api";
import toast from "react-hot-toast";

const NewChatModal = ({ isOpen, onClose, onChatCreated }) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [users, setUsers] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isFetching, setIsFetching] = useState(false);
  const searchTimeoutRef = useRef(null);

  // Fetch users dengan search atau tanpa search
  useEffect(() => {
    const fetchUsers = async () => {
      if (!isOpen) return;

      setIsFetching(true);
      try {
        let response;
        // Deteksi jika query adalah email Gmail atau username Gmail
        // Jika mengandung @gmail.com -> pasti Gmail
        // Jika tidak mengandung @ sama sekali -> bisa jadi username Gmail, coba search Gmail juga
        const hasAtSymbol = searchQuery.includes("@");
        const isGmailSearch =
          searchQuery.includes("@gmail.com") ||
          (!hasAtSymbol && searchQuery.trim() !== "");

        if (searchQuery.trim()) {
          if (isGmailSearch) {
            // Gunakan pencarian Gmail khusus
            // Jika tidak ada @, backend akan handle dengan menambahkan @gmail.com
            response = await userAPI.searchUsersByGmail(searchQuery);
          } else {
            // Gunakan pencarian umum (username atau email apapun)
            response = await userAPI.searchUsers(searchQuery);
          }
        } else {
          // Tanpa search, ambil semua user
          response = await userAPI.getAllUsers();
        }
        setUsers(response.data || []);
      } catch (error) {
        console.error("Error fetching users:", error);
        toast.error("Failed to load users");
        setUsers([]);
      } finally {
        setIsFetching(false);
      }
    };

    // Debounce search untuk menghindari terlalu banyak request
    if (searchTimeoutRef.current) {
      clearTimeout(searchTimeoutRef.current);
    }

    searchTimeoutRef.current = setTimeout(
      () => {
        fetchUsers();
      },
      searchQuery.trim() ? 300 : 0
    ); // Delay 300ms jika ada search query

    return () => {
      if (searchTimeoutRef.current) {
        clearTimeout(searchTimeoutRef.current);
      }
    };
  }, [isOpen, searchQuery]);

  if (!isOpen) return null;

  const handleStartChat = async (userId) => {
    setIsLoading(true);
    try {
      const response = await roomAPI.createPrivateRoom(userId);
      const room = response.data;

      toast.success("Chat started!");
      onChatCreated(room);
      onClose();
      setSearchQuery("");
    } catch (error) {
      toast.error("Failed to start chat");
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 animate-fade-in">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md mx-4 max-h-[80vh] flex flex-col animate-slide-up">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <h2 className="text-xl font-bold text-gray-900">New Chat</h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X className="w-5 h-5 text-gray-500" />
          </button>
        </div>

        {/* Search */}
        <div className="p-6 border-b border-gray-200">
          <Input
            type="text"
            placeholder="Search people..."
            icon={Search}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            autoFocus
          />
        </div>

        {/* Users List */}
        <div className="flex-1 overflow-y-auto p-6 custom-scrollbar">
          {isFetching ? (
            <div className="flex flex-col items-center justify-center py-12">
              <Loader2 className="w-8 h-8 text-primary-600 animate-spin mb-3" />
              <p className="text-gray-500">Loading users...</p>
            </div>
          ) : users.length === 0 ? (
            <div className="text-center py-8">
              <MessageCircle className="w-12 h-12 text-gray-400 mx-auto mb-3" />
              <p className="text-gray-500">
                {searchQuery ? "No users found" : "No users available"}
              </p>
            </div>
          ) : (
            <div className="space-y-2">
              {users.map((user) => (
                <button
                  key={user.id}
                  onClick={() => handleStartChat(user.id)}
                  disabled={isLoading}
                  className="w-full flex items-center space-x-3 p-3 rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <Avatar alt={user.username} size="md" status={user.status} />
                  <div className="flex-1 text-left">
                    <h4 className="font-semibold text-gray-900">
                      {user.username}
                    </h4>
                    <p className="text-sm text-gray-500">{user.email}</p>
                  </div>
                  {isLoading ? (
                    <Loader2 className="w-5 h-5 text-primary-600 animate-spin" />
                  ) : (
                    <MessageCircle className="w-5 h-5 text-gray-400" />
                  )}
                </button>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default NewChatModal;
