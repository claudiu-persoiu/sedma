package game

import (
	"errors"
	"fmt"
	"math/rand"
)

// Hub games hub
type Hub struct {
	games    map[string]*game
	Messages chan string
	users    map[*Client]string
}

// NewHub create new Hub
func NewHub() *Hub {
	hub := &Hub{
		Messages: make(chan string),
		games:    make(map[string]*game),
		users:    make(map[*Client]string)}
	return hub
}

// Start start a new game
func (h *Hub) Start(client *Client) error {
	key, ok := h.users[client]
	if !ok {
		g := newGame()
		key = h.registerGame(g)
		if err := h.join(key, client); err != nil {
			return err
		}
	}
	fmt.Println("Starting game: " + key)
	client.Send <- &message{Action: "start", Data: key}

	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (h *Hub) registerGame(game *game) string {
	key := randSeq(7)
	h.games[key] = game

	return key
}

func (h *Hub) join(gameKey string, client *Client) error {
	g, ok := h.games[gameKey]

	if ok == false {
		return errors.New("invalid")
	}

	err := g.AddPlayer(client)

	if err == nil {
		h.users[client] = gameKey
	}

	return err
}

func (h *Hub) begin(client *Client) error {
	g, err := getGameFromClient(h, client)
	if err != nil {
		return err
	}

	return g.Start(client)
}

func (h *Hub) play(client *Client, cardIndex int) error {
	g, err := getGameFromClient(h, client)
	if err != nil {
		return err
	}

	return g.play(client, cardIndex)
}

func (h *Hub) fetchHand(client *Client) error {
	g, err := getGameFromClient(h, client)
	if err != nil {
		return err
	}

	return g.fetchHand(client)
}

func getGameFromClient(h *Hub, c *Client) (*game, error) {
	gKey, ok := h.users[c]

	if ok == false {
		return nil, errors.New("invalid user in hub")
	}

	g, ok := h.games[gKey]

	if ok == false {
		return nil, errors.New("invalid user game key")
	}

	return g, nil
}
