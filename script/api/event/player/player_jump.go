package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerJumpEvent se dispara cuando un jugador salta.
type PlayerJumpEvent struct {
	player    *player.Player
	cancelled bool
}

func NewPlayerJumpEvent(p *player.Player) *PlayerJumpEvent {
	return &PlayerJumpEvent{player: p}
}

func (e *PlayerJumpEvent) GetEventName() string           { return "PlayerJumpEvent" }
func (e *PlayerJumpEvent) GetHandlers() *event.HandlerList { return PlayerJumpEventHandlers }
func (e *PlayerJumpEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerJumpEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerJumpEvent) GetPlayer() *player.Player      { return e.player }

var PlayerJumpEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerJumpEvent{}))

func init() { event.RegisterEvent("PlayerJumpEvent", PlayerJumpEventHandlers) }
