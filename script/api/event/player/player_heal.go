package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerHealEvent se dispara cuando un jugador recupera vida.
// La cantidad de curación puede modificarse o cancelarse.
type PlayerHealEvent struct {
	player    *player.Player
	health    float64
	source    world.HealingSource
	cancelled bool
}

func NewPlayerHealEvent(p *player.Player, health float64, src world.HealingSource) *PlayerHealEvent {
	return &PlayerHealEvent{player: p, health: health, source: src}
}

func (e *PlayerHealEvent) GetEventName() string           { return "PlayerHealEvent" }
func (e *PlayerHealEvent) GetHandlers() *event.HandlerList { return PlayerHealEventHandlers }
func (e *PlayerHealEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerHealEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerHealEvent) GetPlayer() *player.Player      { return e.player }
func (e *PlayerHealEvent) GetHealth() float64             { return e.health }
func (e *PlayerHealEvent) SetHealth(h float64)            { e.health = h }

var PlayerHealEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerHealEvent{}))

func init() { event.RegisterEvent("PlayerHealEvent", PlayerHealEventHandlers) }
