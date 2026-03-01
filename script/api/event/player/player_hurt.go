package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerHurtEvent se dispara cuando un jugador recibe daño.
// El daño puede modificarse o el evento puede cancelarse para ignorar el daño.
type PlayerHurtEvent struct {
	player    *player.Player
	damage    float64
	source    world.DamageSource
	cancelled bool
}

func NewPlayerHurtEvent(p *player.Player, damage float64, src world.DamageSource) *PlayerHurtEvent {
	return &PlayerHurtEvent{player: p, damage: damage, source: src}
}

func (e *PlayerHurtEvent) GetEventName() string           { return "PlayerHurtEvent" }
func (e *PlayerHurtEvent) GetHandlers() *event.HandlerList { return PlayerHurtEventHandlers }
func (e *PlayerHurtEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerHurtEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerHurtEvent) GetPlayer() *player.Player      { return e.player }
func (e *PlayerHurtEvent) GetDamage() float64             { return e.damage }
func (e *PlayerHurtEvent) SetDamage(d float64)            { e.damage = d }

var PlayerHurtEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerHurtEvent{}))

func init() { event.RegisterEvent("PlayerHurtEvent", PlayerHurtEventHandlers) }
