import React from "react";
import { X, Crown, UserPlus } from "lucide-react";
import Avatar from "../ui/Avatar";
import Button from "../ui/Button";

const MembersModal = ({ isOpen, onClose, members, room, onAddMember }) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 animate-fade-in">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-lg mx-4 max-h-[80vh] flex flex-col animate-slide-up">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <div>
            <h2 className="text-xl font-bold text-gray-900">Room Members</h2>
            <p className="text-sm text-gray-500 mt-1">
              {members.length} {members.length === 1 ? "member" : "members"}
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X className="w-5 h-5 text-gray-500" />
          </button>
        </div>

        {/* Members List */}
        <div className="flex-1 overflow-y-auto p-6 custom-scrollbar">
          <div className="space-y-3">
            {members.map((member) => (
              <div
                key={member.id}
                className="flex items-center space-x-3 p-3 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <Avatar
                  alt={member.username}
                  src={member.avatar_url}
                  size="md"
                  status={member.status}
                />
                <div className="flex-1 min-w-0">
                  <div className="flex items-center space-x-2">
                    <h4 className="font-semibold text-gray-900 truncate">
                      {member.username}
                    </h4>
                    {member.id === room?.created_by && (
                      <Crown className="w-4 h-4 text-yellow-500" />
                    )}
                  </div>
                  <p className="text-sm text-gray-500 truncate">
                    {member.email}
                  </p>
                </div>
                <div>
                  <span
                    className={`
                    px-2 py-1 text-xs font-medium rounded-full
                    ${
                      member.status === "online"
                        ? "bg-green-100 text-green-700"
                        : "bg-gray-100 text-gray-600"
                    }
                  `}
                  >
                    {member.status}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Footer */}
        {room?.type === "group" && (
          <div className="p-6 border-t border-gray-200">
            <Button
              variant="outline"
              className="w-full"
              icon={UserPlus}
              onClick={onAddMember}
            >
              Add Members
            </Button>
          </div>
        )}
      </div>
    </div>
  );
};

export default MembersModal;
