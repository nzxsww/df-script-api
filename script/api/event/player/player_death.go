package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// PlayerDeathEvent se dispara cuando un jugador muere.
// keepInventory puede modificarse para que el jugador conserve su inventario.
type PlayerDeathEvent struct {
	player        *player.Player
	source        world.DamageSource
	keepInventory bool
	cancelled     bool
}

func NewPlayerDeathEvent(p *player.Player, src world.DamageSource, keepInv bool) *PlayerDeathEvent {
	return &PlayerDeathEvent{
		player:        p,
		source:        src,
		keepInventory: keepInv,
	}
}

func (e *PlayerDeathEvent) GetEventName() string           { return "PlayerDeathEvent" }
func (e *PlayerDeathEvent) GetHandlers() *event.HandlerList { return PlayerDeathEventHandlers }
func (e *PlayerDeathEvent) IsCancelled() bool              { return e.cancelled }
func (e *PlayerDeathEvent) SetCancelled(c bool)            { e.cancelled = c }
func (e *PlayerDeathEvent) GetPlayer() *player.Player      { return e.player }
func (e *PlayerDeathEvent) GetKeepInventory() bool         { return e.keepInventory }
func (e *PlayerDeathEvent) SetKeepInventory(keep bool)     { e.keepInventory = keep }

var PlayerDeathEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerDeathEvent{}))

func init() { event.RegisterEvent("PlayerDeathEvent", PlayerDeathEventHandlers) }
