package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerItemDropEvent se dispara cuando un jugador tira un item al suelo.
type PlayerItemDropEvent struct {
	player    *player.Player
	item      item.Stack
	cancelled bool
}

func NewPlayerItemDropEvent(p *player.Player, it item.Stack) *PlayerItemDropEvent {
	return &PlayerItemDropEvent{player: p, item: it}
}

func (e *PlayerItemDropEvent) GetEventName() string           { return "PlayerItemDropEvent" }
func (e *PlayerItemDropEvent) GetHandlers() *event.HandlerList { return PlayerItemDropEventHandlers }
func (e *PlayerItemDropEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerItemDropEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerItemDropEvent) GetPlayer() *player.Player      { return e.player }
func (e *PlayerItemDropEvent) GetItemCount() int              { return e.item.Count() }

var PlayerItemDropEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerItemDropEvent{}))

func init() { event.RegisterEvent("PlayerItemDropEvent", PlayerItemDropEventHandlers) }
