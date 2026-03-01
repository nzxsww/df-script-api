package player

import (
	"reflect"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

type BlockBreakEvent struct {
	player    *player.Player
	block     cube.Pos
	cancelled bool
}

func NewBlockBreakEvent(p *player.Player, pos cube.Pos) *BlockBreakEvent {
	return &BlockBreakEvent{
		player: p,
		block:  pos,
	}
}

func (e *BlockBreakEvent) GetEventName() string {
	return "BlockBreakEvent"
}

func (e *BlockBreakEvent) GetHandlers() *event.HandlerList {
	return BlockBreakEventHandlers
}

func (e *BlockBreakEvent) IsCancelled() bool {
	return e.cancelled
}

func (e *BlockBreakEvent) SetCancelled(cancelled bool) {
	e.cancelled = cancelled
}

func (e *BlockBreakEvent) GetPlayer() *player.Player {
	return e.player
}

func (e *BlockBreakEvent) GetBlock() cube.Pos {
	return e.block
}

var BlockBreakEventHandlers = event.NewHandlerList(reflect.TypeOf(&BlockBreakEvent{}))

func init() {
	event.RegisterEvent("BlockBreakEvent", BlockBreakEventHandlers)
}

type BlockPlaceEvent struct {
	player    *player.Player
	block     cube.Pos
	blockType item.Stack
	cancelled bool
}

func NewBlockPlaceEvent(p *player.Player, pos cube.Pos, b item.Stack) *BlockPlaceEvent {
	return &BlockPlaceEvent{
		player:    p,
		block:     pos,
		blockType: b,
	}
}

func (e *BlockPlaceEvent) GetEventName() string {
	return "BlockPlaceEvent"
}

func (e *BlockPlaceEvent) GetHandlers() *event.HandlerList {
	return BlockPlaceEventHandlers
}

func (e *BlockPlaceEvent) IsCancelled() bool {
	return e.cancelled
}

func (e *BlockPlaceEvent) SetCancelled(cancelled bool) {
	e.cancelled = cancelled
}

func (e *BlockPlaceEvent) GetPlayer() *player.Player {
	return e.player
}

func (e *BlockPlaceEvent) GetBlock() cube.Pos {
	return e.block
}

func (e *BlockPlaceEvent) GetBlockStack() item.Stack {
	return e.blockType
}

var BlockPlaceEventHandlers = event.NewHandlerList(reflect.TypeOf(&BlockPlaceEvent{}))

func init() {
	event.RegisterEvent("BlockPlaceEvent", BlockPlaceEventHandlers)
}
