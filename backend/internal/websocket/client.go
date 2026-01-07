package websocket

import (
    "encoding/json"
    "log"
    "time"

    "github.com/google/uuid"
    "github.com/gorilla/websocket"
)

type Client struct {
    ID       uuid.UUID
    Hub      *Hub
    Conn     *websocket.Conn
    Send     chan []byte
    Rooms    map[uuid.UUID]bool
    Username string
}

func NewClient(hub *Hub, conn *websocket.Conn, userID uuid.UUID, username string) *Client {
    return &Client{
        ID:       userID,
        Hub:      hub,
        Conn:     conn,
        Send:     make(chan []byte, 256),
        Rooms:    make(map[uuid.UUID]bool),
        Username: username,
    }
}

// ReadPump pumps messages from the websocket connection to the hub
func (c *Client) ReadPump() {
    defer func() {
        c.Hub.Unregister <- c
        c.Conn.Close()
    }()

    c.Conn.SetReadDeadline(time.Now().Add(PongWait))
    c.Conn.SetPongHandler(func(string) error {
        c.Conn.SetReadDeadline(time.Now().Add(PongWait))
        return nil
    })

    c.Conn.SetReadLimit(MaxMessageSize)

    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }

        var msg Message
        if err := json.Unmarshal(message, &msg); err != nil {
            log.Printf("error unmarshaling message: %v", err)
            continue
        }

        msg.SenderID = c.ID
        msg.Username = c.Username
        msg.Timestamp = time.Now()

        c.Hub.Broadcast <- &msg
    }
}

// WritePump pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
    ticker := time.NewTicker(PingPeriod)
    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.Send:
            c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
            if !ok {
                c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.Conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            // Add queued messages to the current websocket message
            n := len(c.Send)
            for i := 0; i < n; i++ {
                w.Write([]byte{'\n'})
                w.Write(<-c.Send)
            }

            if err := w.Close(); err != nil {
                return
            }

        case <-ticker.C:
            c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}