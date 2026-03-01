package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerFoodLossEvent se dispara cuando un jugador pierde hambre.
// El nuevo nivel de hambre puede modificarse o cancelarse.
type PlayerFoodLossEvent struct {
	player    *player.Player
	from      int
	to        int
	cancelled bool
}

func NewPlayerFoodLossEvent(p *player.Player, from, to int) *PlayerFoodLossEvent {
	return &PlayerFoodLossEvent{player: p, from: from, to: to}
}

func (e *PlayerFoodLossEvent) GetEventName() string           { return "PlayerFoodLossEvent" }
func (e *PlayerFoodLossEvent) GetHandlers() *event.HandlerList { return PlayerFoodLossEventHandlers }
func (e *PlayerFoodLossEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerFoodLossEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerFoodLossEvent) GetPlayer() *player.Player      { return e.player }
func (e *PlayerFoodLossEvent) GetFrom() int                   { return e.from }
func (e *PlayerFoodLossEvent) GetTo() int                     { return e.to }
func (e *PlayerFoodLossEvent) SetTo(to int)                   { e.to = to }

var PlayerFoodLossEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerFoodLossEvent{}))

func init() { event.RegisterEvent("PlayerFoodLossEvent", PlayerFoodLossEventHandlers) }
