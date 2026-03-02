package api

import (
	"net"
	"time"

	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	dfplayer "github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/nzxsww/dragonfly-script-api/script/api/event/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/plugin"
)

type dragonflyHandler struct {
	srv       *server.Server
	pluginMgr *plugin.Manager
}

func newDragonflyHandler(srv *server.Server, mgr *plugin.Manager) *dragonflyHandler {
	return &dragonflyHandler{
		srv:       srv,
		pluginMgr: mgr,
	}
}

func (h *dragonflyHandler) HandleMove(ctx *dfplayer.Context, newPos mgl64.Vec3, newRot cube.Rotation) {
	p := ctx.Val()

	e := player.NewPlayerMoveEvent(p, p.Position(), newPos)
	h.pluginMgr.CallEvent(e)

	if e.IsCancelled() {
		ctx.Cancel()
	}
}

func (h *dragonflyHandler) HandleChat(ctx *dfplayer.Context, message *string) {
	p := ctx.Val()

	e := player.NewPlayerChatEvent(p, *message)
	h.pluginMgr.CallEvent(e)

	if e.IsCancelled() {
		ctx.Cancel()
		return
	}

	// Cancelamos el mensaje original de Dragonfly (que usa formato <Nombre> mensaje)
	// y lo enviamos manualmente al chat global con el formato modificado por los plugins.
	// Esto evita el doble prefijo: "<Nombre> [Nombre] mensaje"
	ctx.Cancel()
	chat.Global.WriteString(e.GetMessage()) //nolint:errcheck
}

func (h *dragonflyHandler) HandleBlockBreak(ctx *dfplayer.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	p := ctx.Val()

	e := player.NewBlockBreakEvent(p, pos)
	h.pluginMgr.CallEvent(e)

	if e.IsCancelled() {
		ctx.Cancel()
	}
}

func (h *dragonflyHandler) HandleBlockPlace(ctx *dfplayer.Context, pos cube.Pos, b world.Block) {
	p := ctx.Val()

	e := player.NewBlockPlaceEvent(p, pos, item.Stack{})
	h.pluginMgr.CallEvent(e)

	if e.IsCancelled() {
		ctx.Cancel()
	}
}

func (h *dragonflyHandler) HandleQuit(p *dfplayer.Player) {
	inv.CloseContainer(p)
	e := player.NewPlayerQuitEvent(p)
	h.pluginMgr.CallEvent(e)
}

func (h *dragonflyHandler) HandleItemUse(ctx *dfplayer.Context) {
	p := ctx.Val()
	e := player.NewPlayerItemUseEvent(p)
	h.pluginMgr.CallEvent(e)
	if e.IsCancelled() {
		ctx.Cancel()
	}
}

func (h *dragonflyHandler) HandleItemUseOnBlock(ctx *dfplayer.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
}
func (h *dragonflyHandler) HandleItemUseOnEntity(ctx *dfplayer.Context, e world.Entity) {}
func (h *dragonflyHandler) HandleItemRelease(ctx *dfplayer.Context, it item.Stack, dur time.Duration) {
}
func (h *dragonflyHandler) HandleItemConsume(ctx *dfplayer.Context, it item.Stack) {}

func (h *dragonflyHandler) HandleAttackEntity(ctx *dfplayer.Context, entity world.Entity, force, height *float64, critical *bool) {
	p := ctx.Val()
	e := player.NewPlayerAttackEntityEvent(p, entity, *force, *critical)
	h.pluginMgr.CallEvent(e)
	if e.IsCancelled() {
		ctx.Cancel()
	}
}

func (h *dragonflyHandler) HandleStartBreak(ctx *dfplayer.Context, pos cube.Pos)               {}
func (h *dragonflyHandler) HandleBlockPick(ctx *dfplayer.Context, pos cube.Pos, b world.Block) {}
func (h *dragonflyHandler) HandleHeldSlotChange(ctx *dfplayer.Context, from, to int)           {}

func (h *dragonflyHandler) HandleItemDrop(ctx *dfplayer.Context, s item.Stack) {
	p := ctx.Val()
	e := player.NewPlayerItemDropEvent(p, s)
	h.pluginMgr.CallEvent(e)
	if e.IsCancelled() {
		ctx.Cancel()
	}
}

func (h *dragonflyHandler) HandleItemPickup(ctx *dfplayer.Context, i *item.Stack) {
	p := ctx.Val()
	e := player.NewPlayerItemPickupEvent(p, *i)
	h.pluginMgr.CallEvent(e)
	if e.IsCancelled() {
		ctx.Cancel()
	}
}

func (h *dragonflyHandler) HandleJump(p *dfplayer.Player) {
	e := player.NewPlayerJumpEvent(p)
	h.pluginMgr.CallEvent(e)
}

func (h *dragonflyHandler) HandleTeleport(ctx *dfplayer.Context, pos mgl64.Vec3) {
	p := ctx.Val()
	e := player.NewPlayerTeleportEvent(p, pos)
	h.pluginMgr.CallEvent(e)
	if e.IsCancelled() {
		ctx.Cancel()
	} else {
		pos = e.GetPosition()
	}
}

func (h *dragonflyHandler) HandleChangeWorld(p *dfplayer.Player, before, after *world.World) {}

func (h *dragonflyHandler) HandleToggleSprint(ctx *dfplayer.Context, after bool) {
	p := ctx.Val()
	e := player.NewPlayerToggleSprintEvent(p, after)
	h.pluginMgr.CallEvent(e)
	if e.IsCancelled() {
		ctx.Cancel()
	}
}

func (h *dragonflyHandler) HandleToggleSneak(ctx *dfplayer.Context, after bool) {
	p := ctx.Val()
	e := player.NewPlayerToggleSneakEvent(p, after)
	h.pluginMgr.CallEvent(e)
	if e.IsCancelled() {
		ctx.Cancel()
	}
}

func (h *dragonflyHandler) HandleFoodLoss(ctx *dfplayer.Context, from int, to *int) {
	p := ctx.Val()
	e := player.NewPlayerFoodLossEvent(p, from, *to)
	h.pluginMgr.CallEvent(e)
	if e.IsCancelled() {
		ctx.Cancel()
	} else {
		*to = e.GetTo()
	}
}

func (h *dragonflyHandler) HandleHeal(ctx *dfplayer.Context, health *float64, src world.HealingSource) {
	p := ctx.Val()
	e := player.NewPlayerHealEvent(p, *health, src)
	h.pluginMgr.CallEvent(e)
	if e.IsCancelled() {
		ctx.Cancel()
	} else {
		*health = e.GetHealth()
	}
}

func (h *dragonflyHandler) HandleHurt(ctx *dfplayer.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	p := ctx.Val()
	e := player.NewPlayerHurtEvent(p, *damage, src)
	h.pluginMgr.CallEvent(e)
	if e.IsCancelled() {
		ctx.Cancel()
	} else {
		*damage = e.GetDamage()
	}
}

func (h *dragonflyHandler) HandleDeath(p *dfplayer.Player, src world.DamageSource, keepInv *bool) {
	e := player.NewPlayerDeathEvent(p, src, *keepInv)
	h.pluginMgr.CallEvent(e)
	*keepInv = e.GetKeepInventory()
}

func (h *dragonflyHandler) HandleRespawn(p *dfplayer.Player, pos *mgl64.Vec3, w **world.World) {
	e := player.NewPlayerRespawnEvent(p, *pos)
	h.pluginMgr.CallEvent(e)
	*pos = e.GetPosition()
}
func (h *dragonflyHandler) HandleSkinChange(ctx *dfplayer.Context, s *skin.Skin)                  {}
func (h *dragonflyHandler) HandleFireExtinguish(ctx *dfplayer.Context, pos cube.Pos)              {}
func (h *dragonflyHandler) HandleExperienceGain(ctx *dfplayer.Context, amount *int) {
	p := ctx.Val()
	e := player.NewPlayerExperienceGainEvent(p, *amount)
	h.pluginMgr.CallEvent(e)
	if e.IsCancelled() {
		ctx.Cancel()
	} else {
		*amount = e.GetAmount()
	}
}
func (h *dragonflyHandler) HandlePunchAir(ctx *dfplayer.Context)                                  {}
func (h *dragonflyHandler) HandleSignEdit(ctx *dfplayer.Context, pos cube.Pos, frontSide bool, oldText, newText string) {
}
func (h *dragonflyHandler) HandleSleep(ctx *dfplayer.Context, sendReminder *bool) {}
func (h *dragonflyHandler) HandleLecternPageTurn(ctx *dfplayer.Context, pos cube.Pos, oldPage int, newPage *int) {
}
func (h *dragonflyHandler) HandleItemDamage(ctx *dfplayer.Context, i item.Stack, damage *int) {}
func (h *dragonflyHandler) HandleTransfer(ctx *dfplayer.Context, addr *net.UDPAddr)           {}
func (h *dragonflyHandler) HandleCommandExecution(ctx *dfplayer.Context, command cmd.Command, args []string) {
	// Los comandos registrados via cmd.Register() son ejecutados automáticamente
	// por Dragonfly antes de llegar aquí. Este handler solo se llama para comandos
	// que NO están registrados en el sistema. Lo dejamos vacío intencionalmente.
}
func (h *dragonflyHandler) HandleDiagnostics(p *dfplayer.Player, d session.Diagnostics) {}
