package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerItemUseEvent se dispara cuando un jugador usa un item (click derecho).
// El evento puede cancelarse para impedir el uso.
type PlayerItemUseEvent struct {
	player    *player.Player
	cancelled bool
}

func NewPlayerItemUseEvent(p *player.Player) *PlayerItemUseEvent {
	return &PlayerItemUseEvent{player: p}
}

func (e *PlayerItemUseEvent) GetEventName() string           { return "PlayerItemUseEvent" }
func (e *PlayerItemUseEvent) GetHandlers() *event.HandlerList { return PlayerItemUseEventHandlers }
func (e *PlayerItemUseEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerItemUseEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerItemUseEvent) GetPlayer() *player.Player      { return e.player }

var PlayerItemUseEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerItemUseEvent{}))

func init() { event.RegisterEvent("PlayerItemUseEvent", PlayerItemUseEventHandlers) }
