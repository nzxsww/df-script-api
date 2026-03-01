package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

type PlayerChatEvent struct {
	player    *player.Player
	message   string
	cancelled bool
}

func NewPlayerChatEvent(p *player.Player, message string) *PlayerChatEvent {
	return &PlayerChatEvent{
		player:  p,
		message: message,
	}
}

func (e *PlayerChatEvent) GetEventName() string {
	return "PlayerChatEvent"
}

func (e *PlayerChatEvent) GetHandlers() *event.HandlerList {
	return PlayerChatEventHandlers
}

func (e *PlayerChatEvent) IsCancelled() bool {
	return e.cancelled
}

func (e *PlayerChatEvent) SetCancelled(cancelled bool) {
	e.cancelled = cancelled
}

func (e *PlayerChatEvent) GetPlayer() *player.Player {
	return e.player
}

func (e *PlayerChatEvent) GetMessage() string {
	return e.message
}

func (e *PlayerChatEvent) SetMessage(message string) {
	e.message = message
}

var PlayerChatEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerChatEvent{}))

func init() {
	event.RegisterEvent("PlayerChatEvent", PlayerChatEventHandlers)
}
