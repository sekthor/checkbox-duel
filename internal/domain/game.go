package domain

import (
	"time"

	"github.com/nrednav/cuid2"
)

// game is the hub (hub-spoke)
type Game struct {
	ID        string
	CreatedAt time.Time
	Started   time.Time
	Finished  time.Time

	broadcast  chan []byte
	register   chan *Player
	unregister chan *Player
	players    map[*Player]bool
}

func NewGame() *Game {
	return &Game{
		ID:        cuid2.Generate(),
		CreatedAt: time.Now(),

		broadcast:  make(chan []byte),
		register:   make(chan *Player),
		unregister: make(chan *Player),
		players:    make(map[*Player]bool),
	}
}

func (g *Game) Run() {
	for {
		select {
		case player := <-g.register:
			g.players[player] = true

		case player := <-g.unregister:
			if _, ok := g.players[player]; ok {
				delete(g.players, player)
				close(player.send)
			}

		case message := <-g.broadcast:
			for player := range g.players {
				select {
				case player.send <- message:
				default:
					close(player.send)
					delete(g.players, player)
				}
			}
		}
	}
}
