package service

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type WsServer struct {

	// subscriberMessageBuffer controls the max number
	// of messages that can be queued for a subscriber
	// before it is kicked.
	//
	// Defaults to 16.
	subscriberMessageBuffer int

	// publishLimiter controls the rate limit applied to the publish endpoint.
	//
	// Defaults to one publish every 100ms with a burst of 8.
	publishLimiter *rate.Limiter

	// logf controls where logs are sent.
	// Defaults to log.Printf.
	logf func(f string, v ...interface{})

	subscribersMu sync.Mutex
	rooms         map[string]*Room
}

// subscriber represents a subscriber.
// Messages are sent on the msgs channel and if the client
// cannot keep up with the messages, closeSlow is called.
type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}

type Room struct {
	subscribers map[*subscriber]struct{}
}

func NewWsServer() *WsServer {
	wsServer := &WsServer{
		subscriberMessageBuffer: 16,
		logf:                    log.Printf,
		rooms:                   make(map[string]*Room),
		publishLimiter:          rate.NewLimiter(rate.Every(time.Millisecond*5), 50),
	}
	return wsServer
}

func (ws *WsServer) Subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request, voteCode string) error {
	var mu sync.Mutex
	var c *websocket.Conn
	var closed bool
	s := &subscriber{
		msgs: make(chan []byte, ws.subscriberMessageBuffer),
		closeSlow: func() {
			mu.Lock()
			defer mu.Unlock()
			closed = true
			if c != nil {
				c.Close(websocket.StatusPolicyViolation, "connection too slow to keep up with messages")
			}
		},
	}
	ws.addSubscriber(voteCode, s)
	defer ws.deleteSubscriber(voteCode, s)

	c2, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	mu.Lock()
	if closed {
		mu.Unlock()
		return net.ErrClosed
	}
	c = c2
	mu.Unlock()
	defer c.CloseNow()

	ctx = c.CloseRead(ctx)

	log.Println("Client connected to room: ", voteCode)

	for {
		select {
		case msg := <-s.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (ws *WsServer) Publish(voteCode string, msg []byte) {
	ws.subscribersMu.Lock()
	defer ws.subscribersMu.Unlock()

	ws.publishLimiter.Wait(context.Background())
	room, ok := ws.rooms[voteCode]

	if !ok {
		log.Println("No vote with given code!")
		return
	}

	for s := range room.subscribers {
		select {
		case s.msgs <- msg:
		default:
			go s.closeSlow()
		}
	}
}

func (ws *WsServer) addSubscriber(voteCode string, s *subscriber) {
	ws.subscribersMu.Lock()
	room, ok := ws.rooms[voteCode]

	if !ok {
		room = &Room{
			subscribers: make(map[*subscriber]struct{}),
		}
		ws.rooms[voteCode] = room
	}

	room.subscribers[s] = struct{}{}
	ws.subscribersMu.Unlock()
}

func (ws *WsServer) deleteSubscriber(voteCode string, s *subscriber) {
	ws.subscribersMu.Lock()

	room, ok := ws.rooms[voteCode]
	if !ok {
		log.Println("No vote with given code!")
		ws.subscribersMu.Unlock()
		return
	}

	delete(room.subscribers, s)
	if len(room.subscribers) == 0 {
		delete(ws.rooms, voteCode)
	}

	ws.subscribersMu.Unlock()
}

func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return c.Write(ctx, websocket.MessageText, msg)
}
