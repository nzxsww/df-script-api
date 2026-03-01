package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerExperienceGainEvent se dispara cuando un jugador gana experiencia.
// La cantidad puede modificarse o cancelarse.
type PlayerExperienceGainEvent struct {
	player    *player.Player
	amount    int
	cancelled bool
}

func NewPlayerExperienceGainEvent(p *player.Player, amount int) *PlayerExperienceGainEvent {
	return &PlayerExperienceGainEvent{player: p, amount: amount}
}

func (e *PlayerExperienceGainEvent) GetEventName() string           { return "PlayerExperienceGainEvent" }
func (e *PlayerExperienceGainEvent) GetHandlers() *event.HandlerList { return PlayerExperienceGainEventHandlers }
func (e *PlayerExperienceGainEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerExperienceGainEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerExperienceGainEvent) GetPlayer() *player.Player      { return e.player }
func (e *PlayerExperienceGainEvent) GetAmount() int                 { return e.amount }
func (e *PlayerExperienceGainEvent) SetAmount(a int)                { e.amount = a }

var PlayerExperienceGainEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerExperienceGainEvent{}))

func init() { event.RegisterEvent("PlayerExperienceGainEvent", PlayerExperienceGainEventHandlers) }
