import React, { useState, useEffect, useRef } from "react";
import { X, Search, UserPlus, Loader2 } from "lucide-react";
import Button from "../ui/Button";
import Input from "../ui/Input";
import Avatar from "../ui/Avatar";
import { roomAPI, userAPI } from "../../services/api";
import toast from "react-hot-toast";

const AddMemberModal = ({ isOpen, onClose, roomId, currentMembers }) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [users, setUsers] = useState([]);
  const [selectedUsers, setSelectedUsers] = useState([]);
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
          response = await userAPI.getAllUsers();
        }

        // Filter out current members dari hasil
        const fetchedUsers = response.data || [];
        const filtered = fetchedUsers.filter((user) => {
          // Handle both member.id and member.user_id for consistency
          return !currentMembers.some((member) => {
            const memberId = member.id || member.user_id;
            return memberId === user.id;
          });
        });
        setUsers(filtered);
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
    );

    return () => {
      if (searchTimeoutRef.current) {
        clearTimeout(searchTimeoutRef.current);
      }
    };
  }, [isOpen, searchQuery, currentMembers]);

  if (!isOpen) return null;

  const toggleUser = (user) => {
    setSelectedUsers((prev) =>
      prev.includes(user.id)
        ? prev.filter((id) => id !== user.id)
        : [...prev, user.id]
    );
  };

  const handleAddMembers = async () => {
    if (selectedUsers.length === 0) return;

    setIsLoading(true);
    try {
      for (const userId of selectedUsers) {
        await roomAPI.addMember(roomId, userId);
      }
      toast.success(`Added ${selectedUsers.length} member(s) successfully`);
      onClose();
      setSelectedUsers([]);
    } catch {
      toast.error("Failed to add members");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 animate-fade-in">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md mx-4 max-h-[80vh] flex flex-col animate-slide-up">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <h2 className="text-xl font-bold text-gray-900">Add Members</h2>
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
            placeholder="Search users..."
            icon={Search}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
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
              <p className="text-gray-500">No users found</p>
            </div>
          ) : (
            <div className="space-y-2">
              {users.map((user) => (
                <button
                  key={user.id}
                  onClick={() => toggleUser(user)}
                  className={`
                    w-full flex items-center space-x-3 p-3 rounded-lg transition-colors
                    ${
                      selectedUsers.includes(user.id)
                        ? "bg-primary-50 border-2 border-primary-500"
                        : "border-2 border-transparent hover:bg-gray-50"
                    }
                  `}
                >
                  <Avatar alt={user.username} size="md" status={user.status} />
                  <div className="flex-1 text-left">
                    <h4 className="font-semibold text-gray-900">
                      {user.username}
                    </h4>
                    <p className="text-sm text-gray-500">{user.email}</p>
                  </div>
                  {selectedUsers.includes(user.id) && (
                    <div className="w-5 h-5 bg-primary-600 rounded-full flex items-center justify-center">
                      <svg
                        className="w-3 h-3 text-white"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                          clipRule="evenodd"
                        />
                      </svg>
                    </div>
                  )}
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="p-6 border-t border-gray-200">
          <div className="flex space-x-3">
            <Button variant="ghost" onClick={onClose} className="flex-1">
              Cancel
            </Button>
            <Button
              variant="primary"
              onClick={handleAddMembers}
              className="flex-1"
              isLoading={isLoading}
              disabled={selectedUsers.length === 0}
              icon={UserPlus}
            >
              Add {selectedUsers.length > 0 && `(${selectedUsers.length})`}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AddMemberModal;
