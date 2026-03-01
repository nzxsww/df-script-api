package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerRespawnEvent se dispara cuando un jugador reaparece tras morir.
// La posición de reaparición puede modificarse.
type PlayerRespawnEvent struct {
	player    *player.Player
	pos       mgl64.Vec3
	cancelled bool
}

func NewPlayerRespawnEvent(p *player.Player, pos mgl64.Vec3) *PlayerRespawnEvent {
	return &PlayerRespawnEvent{player: p, pos: pos}
}

func (e *PlayerRespawnEvent) GetEventName() string           { return "PlayerRespawnEvent" }
func (e *PlayerRespawnEvent) GetHandlers() *event.HandlerList { return PlayerRespawnEventHandlers }
func (e *PlayerRespawnEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerRespawnEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerRespawnEvent) GetPlayer() *player.Player      { return e.player }
func (e *PlayerRespawnEvent) GetX() float64                  { return e.pos.X() }
func (e *PlayerRespawnEvent) GetY() float64                  { return e.pos.Y() }
func (e *PlayerRespawnEvent) GetZ() float64                  { return e.pos.Z() }
func (e *PlayerRespawnEvent) SetPosition(x, y, z float64)   { e.pos = mgl64.Vec3{x, y, z} }
func (e *PlayerRespawnEvent) GetPosition() mgl64.Vec3        { return e.pos }

var PlayerRespawnEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerRespawnEvent{}))

func init() { event.RegisterEvent("PlayerRespawnEvent", PlayerRespawnEventHandlers) }
