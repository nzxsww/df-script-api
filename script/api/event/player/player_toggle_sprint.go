package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerToggleSprintEvent se dispara cuando un jugador activa o desactiva el sprint.
type PlayerToggleSprintEvent struct {
	player    *player.Player
	sprinting bool
	cancelled bool
}

func NewPlayerToggleSprintEvent(p *player.Player, sprinting bool) *PlayerToggleSprintEvent {
	return &PlayerToggleSprintEvent{player: p, sprinting: sprinting}
}

func (e *PlayerToggleSprintEvent) GetEventName() string           { return "PlayerToggleSprintEvent" }
func (e *PlayerToggleSprintEvent) GetHandlers() *event.HandlerList { return PlayerToggleSprintEventHandlers }
func (e *PlayerToggleSprintEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerToggleSprintEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerToggleSprintEvent) GetPlayer() *player.Player      { return e.player }
func (e *PlayerToggleSprintEvent) IsSprinting() bool              { return e.sprinting }

var PlayerToggleSprintEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerToggleSprintEvent{}))

func init() { event.RegisterEvent("PlayerToggleSprintEvent", PlayerToggleSprintEventHandlers) }
