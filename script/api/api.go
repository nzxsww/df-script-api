package api

import (
	"net"
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	dfplayer "github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/nzxsww/dragonfly-script-api/script/api/event/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/plugin"
)

type PlayerHandler interface {
	HandleMove(ctx *dfplayer.Context, newPos mgl64.Vec3, newRot cube.Rotation)
	HandleChat(ctx *dfplayer.Context, message *string)
	HandleBlockBreak(ctx *dfplayer.Context, pos cube.Pos, drops *[]item.Stack, xp *int)
	HandleBlockPlace(ctx *dfplayer.Context, pos cube.Pos, b world.Block)
	HandleQuit(p *dfplayer.Player)
	HandleItemUse(ctx *dfplayer.Context)
	HandleItemUseOnBlock(ctx *dfplayer.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3)
	HandleItemUseOnEntity(ctx *dfplayer.Context, e world.Entity)
	HandleItemRelease(ctx *dfplayer.Context, item item.Stack, dur time.Duration)
	HandleItemConsume(ctx *dfplayer.Context, item item.Stack)
	HandleAttackEntity(ctx *dfplayer.Context, e world.Entity, force, height *float64, critical *bool)
	HandleStartBreak(ctx *dfplayer.Context, pos cube.Pos)
	HandleBlockPick(ctx *dfplayer.Context, pos cube.Pos, b world.Block)
	HandleHeldSlotChange(ctx *dfplayer.Context, from, to int)
	HandleItemDrop(ctx *dfplayer.Context, s item.Stack)
	HandleItemPickup(ctx *dfplayer.Context, i *item.Stack)
	HandleJump(p *dfplayer.Player)
	HandleTeleport(ctx *dfplayer.Context, pos mgl64.Vec3)
	HandleChangeWorld(p *dfplayer.Player, before, after *world.World)
	HandleToggleSprint(ctx *dfplayer.Context, after bool)
	HandleToggleSneak(ctx *dfplayer.Context, after bool)
	HandleFoodLoss(ctx *dfplayer.Context, from int, to *int)
	HandleHeal(ctx *dfplayer.Context, health *float64, src world.HealingSource)
	HandleHurt(ctx *dfplayer.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource)
	HandleDeath(p *dfplayer.Player, src world.DamageSource, keepInv *bool)
	HandleRespawn(p *dfplayer.Player, pos *mgl64.Vec3, w **world.World)
	HandleSkinChange(ctx *dfplayer.Context, s *skin.Skin)
	HandleFireExtinguish(ctx *dfplayer.Context, pos cube.Pos)
	HandleExperienceGain(ctx *dfplayer.Context, amount *int)
	HandlePunchAir(ctx *dfplayer.Context)
	HandleSignEdit(ctx *dfplayer.Context, pos cube.Pos, frontSide bool, oldText, newText string)
	HandleSleep(ctx *dfplayer.Context, sendReminder *bool)
	HandleLecternPageTurn(ctx *dfplayer.Context, pos cube.Pos, oldPage int, newPage *int)
	HandleItemDamage(ctx *dfplayer.Context, i item.Stack, damage *int)
	HandleTransfer(ctx *dfplayer.Context, addr *net.UDPAddr)
	HandleCommandExecution(ctx *dfplayer.Context, command cmd.Command, args []string)
	HandleDiagnostics(p *dfplayer.Player, d session.Diagnostics)
}

type API struct {
	server    *server.Server
	pluginMgr *plugin.Manager
	handler   *dragonflyHandler
}

func New(srv *server.Server) *API {
	return &API{
		server:    srv,
		pluginMgr: plugin.NewManager(srv),
	}
}

func NewHandler(srv *server.Server, mgr *plugin.Manager) PlayerHandler {
	return newDragonflyHandler(srv, mgr)
}

func (a *API) Start() {
	a.handler = newDragonflyHandler(a.server, a.pluginMgr)
}

func (a *API) AcceptPlayer(p *dfplayer.Player) {
	p.Handle(a.handler)

	e := player.NewPlayerJoinEvent(p)
	a.pluginMgr.CallEvent(e)
}

func (a *API) Server() *server.Server {
	return a.server
}

func (a *API) PluginManager() *plugin.Manager {
	return a.pluginMgr
}
