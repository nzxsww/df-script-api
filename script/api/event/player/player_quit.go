package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

type PlayerQuitEvent struct {
	player      *player.Player
	quitMessage string
	cancelled   bool
}

func NewPlayerQuitEvent(p *player.Player) *PlayerQuitEvent {
	msg := ""
	if p != nil {
		msg = p.Name() + " salió del servidor"
	}
	return &PlayerQuitEvent{
		player:      p,
		quitMessage: msg,
	}
}

func (e *PlayerQuitEvent) GetEventName() string {
	return "PlayerQuitEvent"
}

func (e *PlayerQuitEvent) GetHandlers() *event.HandlerList {
	return PlayerQuitEventHandlers
}

func (e *PlayerQuitEvent) IsCancelled() bool {
	return e.cancelled
}

func (e *PlayerQuitEvent) SetCancelled(cancelled bool) {
	e.cancelled = cancelled
}

func (e *PlayerQuitEvent) GetPlayer() *player.Player {
	return e.player
}

func (e *PlayerQuitEvent) GetQuitMessage() string {
	return e.quitMessage
}

func (e *PlayerQuitEvent) SetQuitMessage(message string) {
	e.quitMessage = message
}

var PlayerQuitEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerQuitEvent{}))

func init() {
	event.RegisterEvent("PlayerQuitEvent", PlayerQuitEventHandlers)
}
