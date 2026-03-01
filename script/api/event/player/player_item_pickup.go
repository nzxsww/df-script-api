package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerItemPickupEvent se dispara cuando un jugador recoge un item del suelo.
type PlayerItemPickupEvent struct {
	player    *player.Player
	item      item.Stack
	cancelled bool
}

func NewPlayerItemPickupEvent(p *player.Player, it item.Stack) *PlayerItemPickupEvent {
	return &PlayerItemPickupEvent{player: p, item: it}
}

func (e *PlayerItemPickupEvent) GetEventName() string           { return "PlayerItemPickupEvent" }
func (e *PlayerItemPickupEvent) GetHandlers() *event.HandlerList { return PlayerItemPickupEventHandlers }
func (e *PlayerItemPickupEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerItemPickupEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerItemPickupEvent) GetPlayer() *player.Player      { return e.player }
func (e *PlayerItemPickupEvent) GetItemCount() int              { return e.item.Count() }

var PlayerItemPickupEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerItemPickupEvent{}))

func init() { event.RegisterEvent("PlayerItemPickupEvent", PlayerItemPickupEventHandlers) }
