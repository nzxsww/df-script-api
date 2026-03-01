package player_test

import (
	"testing"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
	evplayer "github.com/nzxsww/dragonfly-script-api/script/api/event/player"
)

// Todos los tests usan player=nil porque solo testean el comportamiento
// del evento (cancelación, getters, setters, nombre, handlers).
// No se llama ningún método sobre el jugador.

// --- PlayerJoinEvent ---

func TestPlayerJoinEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerJoinEvent(nil)
	if e.GetEventName() != "PlayerJoinEvent" {
		t.Errorf("expected 'PlayerJoinEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerJoinEvent_Handlers(t *testing.T) {
	e := evplayer.NewPlayerJoinEvent(nil)
	if e.GetHandlers() == nil {
		t.Error("GetHandlers() no debe retornar nil")
	}
	if e.GetHandlers() != evplayer.PlayerJoinEventHandlers {
		t.Error("GetHandlers() debe retornar la HandlerList global del evento")
	}
}

func TestPlayerJoinEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerJoinEvent(nil)
	if e.IsCancelled() {
		t.Error("evento nuevo no debe estar cancelado")
	}
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
	e.SetCancelled(false)
	if e.IsCancelled() {
		t.Error("SetCancelled(false) no funcionó")
	}
}

func TestPlayerJoinEvent_JoinMessage(t *testing.T) {
	e := evplayer.NewPlayerJoinEvent(nil)
	e.SetJoinMessage("Hola mundo")
	if e.GetJoinMessage() != "Hola mundo" {
		t.Errorf("expected 'Hola mundo', got '%s'", e.GetJoinMessage())
	}
}

func TestPlayerJoinEvent_RegisteredGlobally(t *testing.T) {
	hl := event.GetEventHandlerList("PlayerJoinEvent")
	if hl == nil {
		t.Error("PlayerJoinEvent no está registrado en el registro global de eventos")
	}
}

// --- PlayerQuitEvent ---

func TestPlayerQuitEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerQuitEvent(nil)
	if e.GetEventName() != "PlayerQuitEvent" {
		t.Errorf("expected 'PlayerQuitEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerQuitEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerQuitEvent(nil)
	if e.IsCancelled() {
		t.Error("evento nuevo no debe estar cancelado")
	}
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerQuitEvent_QuitMessage(t *testing.T) {
	e := evplayer.NewPlayerQuitEvent(nil)
	e.SetQuitMessage("Adiós!")
	if e.GetQuitMessage() != "Adiós!" {
		t.Errorf("expected 'Adiós!', got '%s'", e.GetQuitMessage())
	}
}

func TestPlayerQuitEvent_RegisteredGlobally(t *testing.T) {
	hl := event.GetEventHandlerList("PlayerQuitEvent")
	if hl == nil {
		t.Error("PlayerQuitEvent no está registrado en el registro global de eventos")
	}
}

// --- PlayerChatEvent ---

func TestPlayerChatEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerChatEvent(nil, "hola")
	if e.GetEventName() != "PlayerChatEvent" {
		t.Errorf("expected 'PlayerChatEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerChatEvent_MessageGetSet(t *testing.T) {
	e := evplayer.NewPlayerChatEvent(nil, "mensaje original")
	if e.GetMessage() != "mensaje original" {
		t.Errorf("expected 'mensaje original', got '%s'", e.GetMessage())
	}
	e.SetMessage("mensaje modificado")
	if e.GetMessage() != "mensaje modificado" {
		t.Errorf("expected 'mensaje modificado', got '%s'", e.GetMessage())
	}
}

func TestPlayerChatEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerChatEvent(nil, "test")
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerChatEvent_RegisteredGlobally(t *testing.T) {
	hl := event.GetEventHandlerList("PlayerChatEvent")
	if hl == nil {
		t.Error("PlayerChatEvent no está registrado en el registro global de eventos")
	}
}

// --- PlayerMoveEvent ---

func TestPlayerMoveEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerMoveEvent(nil, mgl64.Vec3{0, 64, 0}, mgl64.Vec3{1, 64, 1})
	if e.GetEventName() != "PlayerMoveEvent" {
		t.Errorf("expected 'PlayerMoveEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerMoveEvent_Positions(t *testing.T) {
	from := mgl64.Vec3{10, 64, 20}
	to := mgl64.Vec3{11, 65, 21}
	e := evplayer.NewPlayerMoveEvent(nil, from, to)

	if e.GetFrom() != from {
		t.Errorf("GetFrom() incorrecto: expected %v, got %v", from, e.GetFrom())
	}
	if e.GetTo() != to {
		t.Errorf("GetTo() incorrecto: expected %v, got %v", to, e.GetTo())
	}
}

func TestPlayerMoveEvent_SetTo(t *testing.T) {
	e := evplayer.NewPlayerMoveEvent(nil, mgl64.Vec3{0, 64, 0}, mgl64.Vec3{1, 64, 1})
	newTo := mgl64.Vec3{99, 100, 99}
	e.SetTo(newTo)
	if e.GetTo() != newTo {
		t.Errorf("SetTo() incorrecto: expected %v, got %v", newTo, e.GetTo())
	}
}

func TestPlayerMoveEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerMoveEvent(nil, mgl64.Vec3{}, mgl64.Vec3{})
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerMoveEvent_RegisteredGlobally(t *testing.T) {
	hl := event.GetEventHandlerList("PlayerMoveEvent")
	if hl == nil {
		t.Error("PlayerMoveEvent no está registrado en el registro global de eventos")
	}
}

// --- BlockBreakEvent ---

func TestBlockBreakEvent_Name(t *testing.T) {
	e := evplayer.NewBlockBreakEvent(nil, cube.Pos{5, 64, 10})
	if e.GetEventName() != "BlockBreakEvent" {
		t.Errorf("expected 'BlockBreakEvent', got '%s'", e.GetEventName())
	}
}

func TestBlockBreakEvent_BlockPos(t *testing.T) {
	pos := cube.Pos{5, 70, 10}
	e := evplayer.NewBlockBreakEvent(nil, pos)
	got := e.GetBlock()
	if got != pos {
		t.Errorf("GetBlock() incorrecto: expected %v, got %v", pos, got)
	}
}

func TestBlockBreakEvent_BlockPos_XYZ(t *testing.T) {
	pos := cube.Pos{3, 64, 7}
	e := evplayer.NewBlockBreakEvent(nil, pos)
	if e.GetBlock().X() != 3 {
		t.Errorf("X() incorrecto: expected 3, got %d", e.GetBlock().X())
	}
	if e.GetBlock().Y() != 64 {
		t.Errorf("Y() incorrecto: expected 64, got %d", e.GetBlock().Y())
	}
	if e.GetBlock().Z() != 7 {
		t.Errorf("Z() incorrecto: expected 7, got %d", e.GetBlock().Z())
	}
}

func TestBlockBreakEvent_Cancelled(t *testing.T) {
	e := evplayer.NewBlockBreakEvent(nil, cube.Pos{})
	if e.IsCancelled() {
		t.Error("evento nuevo no debe estar cancelado")
	}
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestBlockBreakEvent_RegisteredGlobally(t *testing.T) {
	hl := event.GetEventHandlerList("BlockBreakEvent")
	if hl == nil {
		t.Error("BlockBreakEvent no está registrado en el registro global de eventos")
	}
}

// --- BlockPlaceEvent ---

func TestBlockPlaceEvent_Name(t *testing.T) {
	e := evplayer.NewBlockPlaceEvent(nil, cube.Pos{0, 64, 0}, item.Stack{})
	if e.GetEventName() != "BlockPlaceEvent" {
		t.Errorf("expected 'BlockPlaceEvent', got '%s'", e.GetEventName())
	}
}

func TestBlockPlaceEvent_BlockPos(t *testing.T) {
	pos := cube.Pos{1, 65, 2}
	e := evplayer.NewBlockPlaceEvent(nil, pos, item.Stack{})
	if e.GetBlock() != pos {
		t.Errorf("GetBlock() incorrecto: expected %v, got %v", pos, e.GetBlock())
	}
}

func TestBlockPlaceEvent_Cancelled(t *testing.T) {
	e := evplayer.NewBlockPlaceEvent(nil, cube.Pos{}, item.Stack{})
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestBlockPlaceEvent_RegisteredGlobally(t *testing.T) {
	hl := event.GetEventHandlerList("BlockPlaceEvent")
	if hl == nil {
		t.Error("BlockPlaceEvent no está registrado en el registro global de eventos")
	}
}
