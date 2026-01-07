import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { MessageSquare } from "lucide-react";
import Sidebar from "../components/layout/Sidebar";
import RoomList from "../components/layout/RoomList";
import ChatHeader from "../components/chat/ChatHeader";
import MessageList from "../components/chat/MessageList";
import ChatInput from "../components/chat/ChatInput";
import MembersModal from "../components/modals/MembersModal";
import AddMemberModal from "../components/modals/AddMemberModal";
import { useChatStore } from "../store/chatStore";
import { useAuthStore } from "../store/authStore";
import { useWebSocket } from "../hooks/useWebSocket.js";
import { MESSAGE_TYPES } from "../utils/constants";

const Chat = () => {
  const navigate = useNavigate();
  const { user } = useAuthStore();
  const { currentRoom, members, fetchRooms, setCurrentRoom } = useChatStore();
  const { sendMessage } = useWebSocket(currentRoom?.id);
  const [showMembersModal, setShowMembersModal] = useState(false);
  const [showAddMemberModal, setShowAddMemberModal] = useState(false);

  useEffect(() => {
    if (!user) {
      navigate("/auth");
      return;
    }
    fetchRooms();
  }, [user, navigate, fetchRooms]);

  const handleSelectRoom = (room) => {
    setCurrentRoom(room);
  };

  const handleSendMessage = (messageData) => {
    sendMessage(messageData);
  };

  const handleTyping = (isTyping) => {
    sendMessage({
      type: MESSAGE_TYPES.TYPING,
      content: isTyping ? "start" : "stop",
    });
  };

  return (
    <>
      <div className="h-screen flex bg-gray-100">
        <Sidebar />
        <RoomList onSelectRoom={handleSelectRoom} currentRoom={currentRoom} />

        <div className="flex-1 flex flex-col">
          {currentRoom ? (
            <>
              <ChatHeader
                room={currentRoom}
                members={members}
                onShowMembers={() => setShowMembersModal(true)}
              />
              <MessageList />
              <ChatInput
                onSendMessage={handleSendMessage}
                onTyping={handleTyping}
              />
            </>
          ) : (
            <div className="flex-1 flex items-center justify-center bg-gray-50">
              <div className="text-center">
                <div className="w-24 h-24 bg-gray-200 rounded-full flex items-center justify-center mx-auto mb-6">
                  <MessageSquare className="w-12 h-12 text-gray-400" />
                </div>
                <h2 className="text-2xl font-bold text-gray-900 mb-2">
                  Welcome to ChatApp
                </h2>
                <p className="text-gray-600">
                  Select a conversation from the sidebar to start chatting
                </p>
              </div>
            </div>
          )}
        </div>
      </div>

      <MembersModal
        isOpen={showMembersModal}
        onClose={() => setShowMembersModal(false)}
        members={members}
        room={currentRoom}
        onAddMember={() => {
          setShowMembersModal(false);
          setShowAddMemberModal(true);
        }}
      />

      <AddMemberModal
        isOpen={showAddMemberModal}
        onClose={() => setShowAddMemberModal(false)}
        roomId={currentRoom?.id}
        currentMembers={members}
      />
    </>
  );
};

export default Chat;
