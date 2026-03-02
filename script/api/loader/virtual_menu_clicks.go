package loader

import (
	"sync"

	"github.com/bedrock-gophers/intercept/intercept"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// clickInfo guarda el ultimo tipo de click detectado para un jugador.
var (
	menuClickMu   sync.Mutex
	menuClickType = map[uuid.UUID]string{}
)

func setLastClickType(p *player.Player, clickType string) {
	menuClickMu.Lock()
	menuClickType[p.UUID()] = clickType
	menuClickMu.Unlock()
}

func getLastClickType(p *player.Player) string {
	menuClickMu.Lock()
	defer menuClickMu.Unlock()
	if v, ok := menuClickType[p.UUID()]; ok {
		return v
	}
	return "unknown"
}

func init() {
	intercept.Hook(menuClickHandler{})
}

type menuClickHandler struct{}

func (menuClickHandler) HandleClientPacket(ctx *intercept.Context, pk packet.Packet) {
	req, ok := pk.(*packet.ItemStackRequest)
	if !ok {
		return
	}
	// Resolver jugador desde el handle.
	ha := ctx.Val()
	ha.ExecWorld(func(tx *world.Tx, e world.Entity) {
		p := e.(*player.Player)
		for _, data := range req.Requests {
			for _, action := range data.Actions {
				if clickType, ok := clickTypeFromAction(action); ok {
					setLastClickType(p, clickType)
				}
			}
		}
	})
}

func (menuClickHandler) HandleServerPacket(_ *intercept.Context, _ packet.Packet) {}

func clickTypeFromAction(action protocol.StackRequestAction) (string, bool) {
	switch act := action.(type) {
	case *protocol.DropStackRequestAction:
		return "drop", true
	case *protocol.SwapStackRequestAction:
		return "swap", true
	case *protocol.PlaceStackRequestAction:
		// Shift click: mover entre menu (level entity) e inventario/hotbar
		if (act.Source.Container.ContainerID == protocol.ContainerLevelEntity &&
			(act.Destination.Container.ContainerID == protocol.ContainerInventory || act.Destination.Container.ContainerID == protocol.ContainerHotBar)) ||
			((act.Source.Container.ContainerID == protocol.ContainerInventory || act.Source.Container.ContainerID == protocol.ContainerHotBar) &&
				act.Destination.Container.ContainerID == protocol.ContainerLevelEntity) {
			return "shift_click", true
		}
		if act.Count == 1 {
			return "right_click", true
		}
		return "left_click", true
	case *protocol.TakeStackRequestAction:
		// Shift click puede llegar como Take entre contenedores
		if (act.Source.Container.ContainerID == protocol.ContainerLevelEntity &&
			(act.Destination.Container.ContainerID == protocol.ContainerInventory || act.Destination.Container.ContainerID == protocol.ContainerHotBar)) ||
			((act.Source.Container.ContainerID == protocol.ContainerInventory || act.Source.Container.ContainerID == protocol.ContainerHotBar) &&
				act.Destination.Container.ContainerID == protocol.ContainerLevelEntity) {
			return "shift_click", true
		}
		if act.Count == 1 {
			return "right_click", true
		}
		return "left_click", true
	case *protocol.DestroyStackRequestAction:
		return "drop", true
	default:
		_ = act
	}
	return "unknown", false
}
