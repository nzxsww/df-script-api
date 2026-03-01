package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

type PlayerMoveEvent struct {
	player    *player.Player
	from      mgl64.Vec3
	to        mgl64.Vec3
	cancelled bool
}

func NewPlayerMoveEvent(p *player.Player, from, to mgl64.Vec3) *PlayerMoveEvent {
	return &PlayerMoveEvent{
		player: p,
		from:   from,
		to:     to,
	}
}

func (e *PlayerMoveEvent) GetEventName() string {
	return "PlayerMoveEvent"
}

func (e *PlayerMoveEvent) GetHandlers() *event.HandlerList {
	return PlayerMoveEventHandlers
}

func (e *PlayerMoveEvent) IsCancelled() bool {
	return e.cancelled
}

func (e *PlayerMoveEvent) SetCancelled(cancelled bool) {
	e.cancelled = cancelled
}

func (e *PlayerMoveEvent) GetPlayer() *player.Player {
	return e.player
}

func (e *PlayerMoveEvent) GetFrom() mgl64.Vec3 {
	return e.from
}

func (e *PlayerMoveEvent) GetTo() mgl64.Vec3 {
	return e.to
}

func (e *PlayerMoveEvent) SetTo(pos mgl64.Vec3) {
	e.to = pos
}

var PlayerMoveEventHandlers = event.NewHandlerList(reflect.TypeOf(&PlayerMoveEvent{}))

func init() {
	event.RegisterEvent("PlayerMoveEvent", PlayerMoveEventHandlers)
}
