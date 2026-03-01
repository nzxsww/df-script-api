package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerToggleSneakEvent se dispara cuando un jugador activa o desactiva el agachado.
type PlayerToggleSneakEvent struct {
	player    *player.Player
	sneaking  bool
	cancelled bool
}

func NewPlayerToggleSneakEvent(p *player.Player, sneaking bool) *PlayerToggleSneakEvent {
	return &PlayerToggleSneakEvent{player: p, sneaking: sneaking}
}

func (e *PlayerToggleSneakEvent) GetEventName() string           { return "PlayerToggleSneakEvent" }
func (e *PlayerToggleSneakEvent) GetHandlers() *event.HandlerList { return PlayerToggleSneakEventHandlers }
func (e *PlayerToggleSneakEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerToggleSneakEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerToggleSneakEvent) GetPlayer() *player.Player      { return e.player }
func (e *PlayerToggleSneakEvent) IsSneaking() bool               { return e.sneaking }

var PlayerToggleSneakEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerToggleSneakEvent{}))

func init() { event.RegisterEvent("PlayerToggleSneakEvent", PlayerToggleSneakEventHandlers) }
