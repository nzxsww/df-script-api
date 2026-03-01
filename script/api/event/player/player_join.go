package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

type PlayerJoinEvent struct {
	player      *player.Player
	joinMessage string
	cancelled   bool
}

func NewPlayerJoinEvent(p *player.Player) *PlayerJoinEvent {
	msg := ""
	if p != nil {
		msg = p.Name() + " se unió al servidor"
	}
	return &PlayerJoinEvent{
		player:      p,
		joinMessage: msg,
	}
}

func (e *PlayerJoinEvent) GetEventName() string {
	return "PlayerJoinEvent"
}

func (e *PlayerJoinEvent) GetHandlers() *event.HandlerList {
	return PlayerJoinEventHandlers
}

func (e *PlayerJoinEvent) IsCancelled() bool {
	return e.cancelled
}

func (e *PlayerJoinEvent) SetCancelled(cancelled bool) {
	e.cancelled = cancelled
}

func (e *PlayerJoinEvent) GetPlayer() *player.Player {
	return e.player
}

func (e *PlayerJoinEvent) GetJoinMessage() string {
	return e.joinMessage
}

func (e *PlayerJoinEvent) SetJoinMessage(message string) {
	e.joinMessage = message
}

var PlayerJoinEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerJoinEvent{}))

func init() {
	event.RegisterEvent("PlayerJoinEvent", PlayerJoinEventHandlers)
}
