package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerAttackEntityEvent se dispara cuando un jugador ataca a una entidad.
// El evento puede cancelarse para impedir el ataque.
type PlayerAttackEntityEvent struct {
	player    *player.Player
	entity    world.Entity
	force     float64
	critical  bool
	cancelled bool
}

func NewPlayerAttackEntityEvent(p *player.Player, entity world.Entity, force float64, critical bool) *PlayerAttackEntityEvent {
	return &PlayerAttackEntityEvent{
		player:   p,
		entity:   entity,
		force:    force,
		critical: critical,
	}
}

func (e *PlayerAttackEntityEvent) GetEventName() string           { return "PlayerAttackEntityEvent" }
func (e *PlayerAttackEntityEvent) GetHandlers() *event.HandlerList { return PlayerAttackEntityEventHandlers }
func (e *PlayerAttackEntityEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerAttackEntityEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerAttackEntityEvent) GetPlayer() *player.Player      { return e.player }
func (e *PlayerAttackEntityEvent) GetForce() float64              { return e.force }
func (e *PlayerAttackEntityEvent) IsCritical() bool               { return e.critical }

var PlayerAttackEntityEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerAttackEntityEvent{}))

func init() { event.RegisterEvent("PlayerAttackEntityEvent", PlayerAttackEntityEventHandlers) }
