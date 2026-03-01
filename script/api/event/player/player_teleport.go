package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerTeleportEvent se dispara cuando un jugador es teleportado.
// La posición de destino puede modificarse o cancelarse.
type PlayerTeleportEvent struct {
	player    *player.Player
	pos       mgl64.Vec3
	cancelled bool
}

func NewPlayerTeleportEvent(p *player.Player, pos mgl64.Vec3) *PlayerTeleportEvent {
	return &PlayerTeleportEvent{player: p, pos: pos}
}

func (e *PlayerTeleportEvent) GetEventName() string           { return "PlayerTeleportEvent" }
func (e *PlayerTeleportEvent) GetHandlers() *event.HandlerList { return PlayerTeleportEventHandlers }
func (e *PlayerTeleportEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerTeleportEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerTeleportEvent) GetPlayer() *player.Player      { return e.player }
func (e *PlayerTeleportEvent) GetX() float64                  { return e.pos.X() }
func (e *PlayerTeleportEvent) GetY() float64                  { return e.pos.Y() }
func (e *PlayerTeleportEvent) GetZ() float64                  { return e.pos.Z() }
func (e *PlayerTeleportEvent) SetPosition(x, y, z float64)   { e.pos = mgl64.Vec3{x, y, z} }
func (e *PlayerTeleportEvent) GetPosition() mgl64.Vec3        { return e.pos }

var PlayerTeleportEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerTeleportEvent{}))

func init() { event.RegisterEvent("PlayerTeleportEvent", PlayerTeleportEventHandlers) }
