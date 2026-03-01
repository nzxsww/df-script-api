package player_test

import (
	"testing"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
	evplayer "github.com/nzxsww/dragonfly-script-api/script/api/event/player"
)

// emptyStack retorna un item.Stack vacío para usar en tests.
func emptyStack() item.Stack { return item.Stack{} }

// --- PlayerJumpEvent ---

func TestPlayerJumpEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerJumpEvent(nil)
	if e.GetEventName() != "PlayerJumpEvent" {
		t.Errorf("expected 'PlayerJumpEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerJumpEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerJumpEvent(nil)
	if e.IsCancelled() {
		t.Error("evento nuevo no debe estar cancelado")
	}
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerJumpEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerJumpEvent") == nil {
		t.Error("PlayerJumpEvent no está registrado globalmente")
	}
}

// --- PlayerDeathEvent ---

func TestPlayerDeathEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerDeathEvent(nil, nil, false)
	if e.GetEventName() != "PlayerDeathEvent" {
		t.Errorf("expected 'PlayerDeathEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerDeathEvent_KeepInventory(t *testing.T) {
	e := evplayer.NewPlayerDeathEvent(nil, nil, false)
	if e.GetKeepInventory() {
		t.Error("keepInventory debería ser false por defecto")
	}
	e.SetKeepInventory(true)
	if !e.GetKeepInventory() {
		t.Error("SetKeepInventory(true) no funcionó")
	}
}

func TestPlayerDeathEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerDeathEvent(nil, nil, false)
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerDeathEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerDeathEvent") == nil {
		t.Error("PlayerDeathEvent no está registrado globalmente")
	}
}

// --- PlayerRespawnEvent ---

func TestPlayerRespawnEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerRespawnEvent(nil, mgl64.Vec3{})
	if e.GetEventName() != "PlayerRespawnEvent" {
		t.Errorf("expected 'PlayerRespawnEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerRespawnEvent_Position(t *testing.T) {
	pos := mgl64.Vec3{10, 64, 20}
	e := evplayer.NewPlayerRespawnEvent(nil, pos)
	if e.GetX() != 10 || e.GetY() != 64 || e.GetZ() != 20 {
		t.Errorf("posición incorrecta: got %v,%v,%v", e.GetX(), e.GetY(), e.GetZ())
	}
}

func TestPlayerRespawnEvent_SetPosition(t *testing.T) {
	e := evplayer.NewPlayerRespawnEvent(nil, mgl64.Vec3{})
	e.SetPosition(5, 100, 15)
	if e.GetX() != 5 || e.GetY() != 100 || e.GetZ() != 15 {
		t.Errorf("SetPosition no funcionó: got %v,%v,%v", e.GetX(), e.GetY(), e.GetZ())
	}
}

func TestPlayerRespawnEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerRespawnEvent") == nil {
		t.Error("PlayerRespawnEvent no está registrado globalmente")
	}
}

// --- PlayerHurtEvent ---

func TestPlayerHurtEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerHurtEvent(nil, 5.0, nil)
	if e.GetEventName() != "PlayerHurtEvent" {
		t.Errorf("expected 'PlayerHurtEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerHurtEvent_Damage(t *testing.T) {
	e := evplayer.NewPlayerHurtEvent(nil, 5.0, nil)
	if e.GetDamage() != 5.0 {
		t.Errorf("expected 5.0, got %v", e.GetDamage())
	}
	e.SetDamage(2.5)
	if e.GetDamage() != 2.5 {
		t.Errorf("SetDamage no funcionó: got %v", e.GetDamage())
	}
}

func TestPlayerHurtEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerHurtEvent(nil, 1.0, nil)
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerHurtEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerHurtEvent") == nil {
		t.Error("PlayerHurtEvent no está registrado globalmente")
	}
}

// --- PlayerHealEvent ---

func TestPlayerHealEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerHealEvent(nil, 4.0, nil)
	if e.GetEventName() != "PlayerHealEvent" {
		t.Errorf("expected 'PlayerHealEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerHealEvent_Health(t *testing.T) {
	e := evplayer.NewPlayerHealEvent(nil, 4.0, nil)
	if e.GetHealth() != 4.0 {
		t.Errorf("expected 4.0, got %v", e.GetHealth())
	}
	e.SetHealth(2.0)
	if e.GetHealth() != 2.0 {
		t.Errorf("SetHealth no funcionó: got %v", e.GetHealth())
	}
}

func TestPlayerHealEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerHealEvent(nil, 1.0, nil)
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerHealEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerHealEvent") == nil {
		t.Error("PlayerHealEvent no está registrado globalmente")
	}
}

// --- PlayerExperienceGainEvent ---

func TestPlayerExperienceGainEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerExperienceGainEvent(nil, 10)
	if e.GetEventName() != "PlayerExperienceGainEvent" {
		t.Errorf("expected 'PlayerExperienceGainEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerExperienceGainEvent_Amount(t *testing.T) {
	e := evplayer.NewPlayerExperienceGainEvent(nil, 10)
	if e.GetAmount() != 10 {
		t.Errorf("expected 10, got %d", e.GetAmount())
	}
	e.SetAmount(25)
	if e.GetAmount() != 25 {
		t.Errorf("SetAmount no funcionó: got %d", e.GetAmount())
	}
}

func TestPlayerExperienceGainEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerExperienceGainEvent(nil, 5)
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerExperienceGainEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerExperienceGainEvent") == nil {
		t.Error("PlayerExperienceGainEvent no está registrado globalmente")
	}
}

// --- PlayerToggleSprintEvent ---

func TestPlayerToggleSprintEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerToggleSprintEvent(nil, true)
	if e.GetEventName() != "PlayerToggleSprintEvent" {
		t.Errorf("expected 'PlayerToggleSprintEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerToggleSprintEvent_IsSprinting(t *testing.T) {
	e := evplayer.NewPlayerToggleSprintEvent(nil, true)
	if !e.IsSprinting() {
		t.Error("IsSprinting debería ser true")
	}
	e2 := evplayer.NewPlayerToggleSprintEvent(nil, false)
	if e2.IsSprinting() {
		t.Error("IsSprinting debería ser false")
	}
}

func TestPlayerToggleSprintEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerToggleSprintEvent") == nil {
		t.Error("PlayerToggleSprintEvent no está registrado globalmente")
	}
}

// --- PlayerToggleSneakEvent ---

func TestPlayerToggleSneakEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerToggleSneakEvent(nil, true)
	if e.GetEventName() != "PlayerToggleSneakEvent" {
		t.Errorf("expected 'PlayerToggleSneakEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerToggleSneakEvent_IsSneaking(t *testing.T) {
	e := evplayer.NewPlayerToggleSneakEvent(nil, true)
	if !e.IsSneaking() {
		t.Error("IsSneaking debería ser true")
	}
}

func TestPlayerToggleSneakEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerToggleSneakEvent") == nil {
		t.Error("PlayerToggleSneakEvent no está registrado globalmente")
	}
}

// --- PlayerItemDropEvent ---

func TestPlayerItemDropEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerItemDropEvent(nil, emptyStack())
	if e.GetEventName() != "PlayerItemDropEvent" {
		t.Errorf("expected 'PlayerItemDropEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerItemDropEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerItemDropEvent(nil, emptyStack())
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerItemDropEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerItemDropEvent") == nil {
		t.Error("PlayerItemDropEvent no está registrado globalmente")
	}
}

// --- PlayerItemPickupEvent ---

func TestPlayerItemPickupEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerItemPickupEvent(nil, emptyStack())
	if e.GetEventName() != "PlayerItemPickupEvent" {
		t.Errorf("expected 'PlayerItemPickupEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerItemPickupEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerItemPickupEvent(nil, emptyStack())
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerItemPickupEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerItemPickupEvent") == nil {
		t.Error("PlayerItemPickupEvent no está registrado globalmente")
	}
}

// --- PlayerFoodLossEvent ---

func TestPlayerFoodLossEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerFoodLossEvent(nil, 20, 18)
	if e.GetEventName() != "PlayerFoodLossEvent" {
		t.Errorf("expected 'PlayerFoodLossEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerFoodLossEvent_FromTo(t *testing.T) {
	e := evplayer.NewPlayerFoodLossEvent(nil, 20, 18)
	if e.GetFrom() != 20 {
		t.Errorf("GetFrom: expected 20, got %d", e.GetFrom())
	}
	if e.GetTo() != 18 {
		t.Errorf("GetTo: expected 18, got %d", e.GetTo())
	}
	e.SetTo(15)
	if e.GetTo() != 15 {
		t.Errorf("SetTo: expected 15, got %d", e.GetTo())
	}
}

func TestPlayerFoodLossEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerFoodLossEvent(nil, 20, 18)
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerFoodLossEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerFoodLossEvent") == nil {
		t.Error("PlayerFoodLossEvent no está registrado globalmente")
	}
}

// --- PlayerTeleportEvent ---

func TestPlayerTeleportEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerTeleportEvent(nil, mgl64.Vec3{})
	if e.GetEventName() != "PlayerTeleportEvent" {
		t.Errorf("expected 'PlayerTeleportEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerTeleportEvent_Position(t *testing.T) {
	e := evplayer.NewPlayerTeleportEvent(nil, mgl64.Vec3{10, 64, 20})
	if e.GetX() != 10 || e.GetY() != 64 || e.GetZ() != 20 {
		t.Errorf("posición incorrecta: %v,%v,%v", e.GetX(), e.GetY(), e.GetZ())
	}
	e.SetPosition(0, 100, 0)
	if e.GetX() != 0 || e.GetY() != 100 || e.GetZ() != 0 {
		t.Errorf("SetPosition no funcionó: %v,%v,%v", e.GetX(), e.GetY(), e.GetZ())
	}
}

func TestPlayerTeleportEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerTeleportEvent(nil, mgl64.Vec3{})
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerTeleportEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerTeleportEvent") == nil {
		t.Error("PlayerTeleportEvent no está registrado globalmente")
	}
}

// --- PlayerAttackEntityEvent ---

func TestPlayerAttackEntityEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerAttackEntityEvent(nil, nil, 1.0, false)
	if e.GetEventName() != "PlayerAttackEntityEvent" {
		t.Errorf("expected 'PlayerAttackEntityEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerAttackEntityEvent_ForceAndCritical(t *testing.T) {
	e := evplayer.NewPlayerAttackEntityEvent(nil, nil, 2.5, true)
	if e.GetForce() != 2.5 {
		t.Errorf("GetForce: expected 2.5, got %v", e.GetForce())
	}
	if !e.IsCritical() {
		t.Error("IsCritical debería ser true")
	}
}

func TestPlayerAttackEntityEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerAttackEntityEvent(nil, nil, 1.0, false)
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerAttackEntityEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerAttackEntityEvent") == nil {
		t.Error("PlayerAttackEntityEvent no está registrado globalmente")
	}
}

// --- PlayerItemUseEvent ---

func TestPlayerItemUseEvent_Name(t *testing.T) {
	e := evplayer.NewPlayerItemUseEvent(nil)
	if e.GetEventName() != "PlayerItemUseEvent" {
		t.Errorf("expected 'PlayerItemUseEvent', got '%s'", e.GetEventName())
	}
}

func TestPlayerItemUseEvent_Cancelled(t *testing.T) {
	e := evplayer.NewPlayerItemUseEvent(nil)
	e.SetCancelled(true)
	if !e.IsCancelled() {
		t.Error("SetCancelled(true) no funcionó")
	}
}

func TestPlayerItemUseEvent_RegisteredGlobally(t *testing.T) {
	if event.GetEventHandlerList("PlayerItemUseEvent") == nil {
		t.Error("PlayerItemUseEvent no está registrado globalmente")
	}
}
