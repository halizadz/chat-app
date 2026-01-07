import { useEffect, useRef, useCallback, useState} from "react";
import { WS_BASE_URL } from "../utils/constants";
import { useChatStore } from "../store/chatStore";
import { useAuthStore } from "../store/authStore";

export const useWebSocket = (roomId) => {
  const wsRef = useRef(null);
  const reconnectTimeoutRef = useRef(null);
  const connectRef = useRef(null);
  const [isConnected, setIsConnected] = useState(false);
  const token = useAuthStore((state) => state.token);
  const addMessage = useChatStore((state) => state.addMessage);
  const addTypingUser = useChatStore((state) => state.addTypingUser);
  const removeTypingUser = useChatStore((state) => state.removeTypingUser);

  const connect = useCallback(() => {
    if (!roomId || !token) return;

    // Prevent multiple connections
    if (wsRef.current && (wsRef.current.readyState === WebSocket.OPEN || wsRef.current.readyState === WebSocket.CONNECTING)) {
      console.log("WebSocket already connecting/connected");
      return;
    }

    const ws = new WebSocket(`${WS_BASE_URL}/${roomId}?token=${token}`);

    ws.onopen = () => {
      console.log("WebSocket connected");
      wsRef.current = ws;
      setIsConnected(true); // 2. Set true saat terbuka
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      switch (data.type) {
        case "message":
        case "file":
          addMessage(data);
          break;
        case "typing":
          if (data.is_typing) {
            addTypingUser(data.user_id, data.username);
            setTimeout(() => removeTypingUser(data.user_id), 3000);
          } else {
            removeTypingUser(data.user_id);
          }
          break;
        case "join":
        case "leave":
          addMessage(data);
          break;
        default:
          break;
      }
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      setIsConnected(false); // 3. Set false jika error
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected");
      wsRef.current = null;
      setIsConnected(false); // 4. Set false saat tertutup

      reconnectTimeoutRef.current = setTimeout(() => {
        if (connectRef.current) connectRef.current();
      }, 3000);
    };

    return ws;
  }, [roomId, token, addMessage, addTypingUser, removeTypingUser]);

  useEffect(() => {
    connectRef.current = connect;
  }, [connect]);

  useEffect(() => {
    const ws = connect();

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (ws) {
        ws.close();
      }
    };
  }, [connect]);

 const sendMessage = useCallback((message) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    }
  }, []);

  return {
    sendMessage,
    isConnected,
  };
};
