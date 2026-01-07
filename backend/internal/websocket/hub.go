package websocket

import (
    "encoding/json"
    "log"
    "sync"
    "time"

    "github.com/google/uuid"
)

type Hub struct {
    // Registered clients
    Clients map[uuid.UUID]*Client

    // Client rooms mapping
    Rooms map[uuid.UUID]map[uuid.UUID]*Client

    // Inbound messages from clients
    Broadcast chan *Message

    // Register requests from clients
    Register chan *Client

    // Unregister requests from clients
    Unregister chan *Client

    // Typing indicators
    Typing chan *TypingIndicator

    mu sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        Clients:    make(map[uuid.UUID]*Client),
        Rooms:      make(map[uuid.UUID]map[uuid.UUID]*Client),
        Broadcast:  make(chan *Message, 256),
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
        Typing:     make(chan *TypingIndicator),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.Register:
            h.mu.Lock()
            h.Clients[client.ID] = client
            log.Printf("Client registered: %s (%s)", client.Username, client.ID)
            h.mu.Unlock()

        case client := <-h.Unregister:
            h.mu.Lock()
            if _, ok := h.Clients[client.ID]; ok {
                // Remove from all rooms
                for roomID := range client.Rooms {
                    h.leaveRoom(client, roomID)
                }
                
                delete(h.Clients, client.ID)
                close(client.Send)
                log.Printf("Client unregistered: %s (%s)", client.Username, client.ID)
            }
            h.mu.Unlock()

        case message := <-h.Broadcast:
            h.mu.RLock()
            if message.Type == "message" || message.Type == "file" {
                // Send to all clients in the room
                if room, ok := h.Rooms[message.RoomID]; ok {
                    messageBytes, err := json.Marshal(message)
                    if err != nil {
                        log.Printf("error marshaling message: %v", err)
                        h.mu.RUnlock()
                        continue
                    }

                    for clientID, client := range room {
                        select {
                        case client.Send <- messageBytes:
                        default:
                            close(client.Send)
                            delete(h.Clients, clientID)
                            delete(room, clientID)
                        }
                    }
                }
            }
            h.mu.RUnlock()

        case typing := <-h.Typing:
            h.mu.RLock()
            if room, ok := h.Rooms[typing.RoomID]; ok {
                typingBytes, err := json.Marshal(map[string]interface{}{
                    "type":      "typing",
                    "room_id":   typing.RoomID,
                    "user_id":   typing.UserID,
                    "username":  typing.Username,
                    "is_typing": typing.IsTyping,
                })
                if err != nil {
                    log.Printf("error marshaling typing indicator: %v", err)
                    h.mu.RUnlock()
                    continue
                }

                for clientID, client := range room {
                    // Don't send typing indicator to the typer
                    if clientID != typing.UserID {
                        select {
                        case client.Send <- typingBytes:
                        default:
                            close(client.Send)
                            delete(h.Clients, clientID)
                            delete(room, clientID)
                        }
                    }
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (h *Hub) JoinRoom(client *Client, roomID uuid.UUID) {
    h.mu.Lock()
    defer h.mu.Unlock()

    if h.Rooms[roomID] == nil {
        h.Rooms[roomID] = make(map[uuid.UUID]*Client)
    }

    h.Rooms[roomID][client.ID] = client
    client.Rooms[roomID] = true

    log.Printf("Client %s joined room %s", client.Username, roomID)

    // Broadcast join message
    joinMsg := &Message{
        Type:      "join",
        RoomID:    roomID,
        SenderID:  client.ID,
        Username:  client.Username,
        Content:   client.Username + " joined the room",
        Timestamp: time.Now(),
    }
    h.Broadcast <- joinMsg
}

func (h *Hub) leaveRoom(client *Client, roomID uuid.UUID) {
    if room, ok := h.Rooms[roomID]; ok {
        delete(room, client.ID)
        delete(client.Rooms, roomID)

        if len(room) == 0 {
            delete(h.Rooms, roomID)
        }

        log.Printf("Client %s left room %s", client.Username, roomID)

        // Broadcast leave message
        leaveMsg := &Message{
            Type:      "leave",
            RoomID:    roomID,
            SenderID:  client.ID,
            Username:  client.Username,
            Content:   client.Username + " left the room",
            Timestamp: time.Now(),
        }
        h.Broadcast <- leaveMsg
    }
}

func (h *Hub) LeaveRoom(client *Client, roomID uuid.UUID) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.leaveRoom(client, roomID)
}