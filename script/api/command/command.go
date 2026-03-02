// Package command implementa el sistema de comandos para plugins JavaScript.
// Permite registrar comandos desde JS que son ejecutados por jugadores en el servidor.
package command

import (
	"fmt"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	dfitem "github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"strings"
	dfplayer "github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/dop251/goja"
	"github.com/go-gl/mathgl/mgl64"
)

// jsCommandCallback es el tipo del callback JS que recibe el jugador, args y el tx activo.
type jsCommandCallback func(player map[string]interface{}, args []string, tx *world.Tx)

// registeredCommands guarda los callbacks JS por nombre de comando.
var (
	registeredCommands   = make(map[string]jsCommandCallback)
	registeredCommandsMu sync.RWMutex
)

// jsRunnable es el struct que Dragonfly usa para ejecutar un comando.
// El campo Args captura todos los argumentos como una sola string (Varargs).
type jsRunnable struct {
	Args        cmd.Varargs `cmd:"args,-"`
	commandName string
}

// Run es llamado por Dragonfly cuando el jugador ejecuta el comando.
// Busca el callback JS registrado y lo llama con el jugador y los argumentos.
func (r jsRunnable) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*dfplayer.Player)
	if !ok {
		o.Error("Este comando solo puede ser ejecutado por jugadores.")
		return
	}

	registeredCommandsMu.RLock()
	callback, exists := registeredCommands[r.commandName]
	registeredCommandsMu.RUnlock()

	if !exists {
		o.Errorf("El comando '%s' no tiene handler registrado.", r.commandName)
		return
	}

	// Parsear los argumentos del Varargs
	args := parseArgs(string(r.Args))

	// Construir el wrapper del jugador para exponerlo al JS
	playerObj := buildPlayerMap(p)

	// Llamar el callback JS con el tx activo
	callback(playerObj, args, tx)
}

// Register registra un nuevo comando en Dragonfly y guarda su callback JS.
// name: nombre del comando (sin /)
// description: descripción que aparece en el autocompletado
// aliases: nombres alternativos (puede ser nil)
// callback: función JS que recibe (player, args)
func Register(name, description string, aliases []string, callback jsCommandCallback) {
	registeredCommandsMu.Lock()
	registeredCommands[name] = callback
	registeredCommandsMu.Unlock()

	runnable := jsRunnable{commandName: name}
	command := cmd.New(name, description, aliases, runnable)
	cmd.Register(command)

	fmt.Printf("[Commands] Comando registrado: /%s\n", name)
}

// parseArgs divide el string de Varargs en una lista de argumentos.
func parseArgs(raw string) []string {
	if raw == "" {
		return []string{}
	}
	args := []string{}
	current := ""
	inQuote := false
	for _, ch := range raw {
		switch {
		case ch == '"':
			inQuote = !inQuote
		case ch == ' ' && !inQuote:
			if current != "" {
				args = append(args, current)
				current = ""
			}
		default:
			current += string(ch)
		}
	}
	if current != "" {
		args = append(args, current)
	}
	return args
}

// buildPlayerMap construye el mapa de métodos del jugador para exponerlo al JS.
// Es el mismo wrapper que usa el loader, replicado aquí para no crear dependencias circulares.
// IMPORTANTE: debe mantenerse sincronizado con newPlayerWrapper en loader/loader.go.
func buildPlayerMap(p *dfplayer.Player) map[string]interface{} {
	return map[string]interface{}{
		// Identidad
		"getName":    func() string { return p.Name() },
		"getUUID":    func() string { return p.UUID().String() },
		"getXUID":    func() string { return p.XUID() },
		"getNameTag": func() string { return p.NameTag() },
		"setNameTag": func(name string) { p.SetNameTag(name) },
		// Mensajes
		"sendMessage":      func(msg string) { p.Message(msg) },
		"sendPopup":        func(msg string) { p.SendPopup(msg) },
		"sendTip":          func(msg string) { p.SendTip(msg) },
		"sendToast":        func(t, msg string) { p.SendToast(t, msg) },
		"sendJukeboxPopup": func(msg string) { p.SendJukeboxPopup(msg) },
		"sendTitle": func(text, subtitle string) {
			t := title.New(text).WithSubtitle(subtitle)
			p.SendTitle(t)
		},
		// Conexión
		"disconnect": func(msg string) { p.Disconnect(msg) },
		"transfer": func(address string) {
			if err := p.Transfer(address); err != nil {
				fmt.Printf("[Commands] transfer error: %v\n", err)
			}
		},
		"getLatency": func() int64 { return p.Latency().Milliseconds() },
		// Posición y movimiento
		"getX":        func() float64 { return p.Position().X() },
		"getY":        func() float64 { return p.Position().Y() },
		"getZ":        func() float64 { return p.Position().Z() },
		"teleport":    func(x, y, z float64) { p.Teleport(mgl64.Vec3{x, y, z}) },
		"setVelocity": func(x, y, z float64) { p.SetVelocity(mgl64.Vec3{x, y, z}) },
		// Estado físico
		"getHealth":    func() float64 { return p.Health() },
		"getMaxHealth": func() float64 { return p.MaxHealth() },
		"setMaxHealth": func(h float64) { p.SetMaxHealth(h) },
		"getFoodLevel": func() int { return p.Food() },
		"setFoodLevel": func(f int) { p.SetFood(f) },
		"isOnGround":   func() bool { return p.OnGround() },
		"isSneaking":   func() bool { return p.Sneaking() },
		"isSprinting":  func() bool { return p.Sprinting() },
		"isFlying":     func() bool { return p.Flying() },
		"isSwimming":   func() bool { return p.Swimming() },
		"isDead":       func() bool { return p.Dead() },
		"isImmobile":   func() bool { return p.Immobile() },
		// Experiencia
		"getExperience":      func() int { return p.Experience() },
		"getExperienceLevel": func() int { return p.ExperienceLevel() },
		"addExperience":      func(amount int) { p.AddExperience(amount) },
		"setExperienceLevel": func(lvl int) { p.SetExperienceLevel(lvl) },
		// Modo de juego
		"setGameMode": func(mode string) {
			switch mode {
			case "survival":
				p.SetGameMode(world.GameModeSurvival)
			case "creative":
				p.SetGameMode(world.GameModeCreative)
			case "adventure":
				p.SetGameMode(world.GameModeAdventure)
			case "spectator":
				p.SetGameMode(world.GameModeSpectator)
			default:
				fmt.Printf("[Commands] setGameMode: modo desconocido '%s'\n", mode)
			}
		},
		"getGameMode": func() string {
			switch p.GameMode() {
			case world.GameModeSurvival:
				return "survival"
			case world.GameModeCreative:
				return "creative"
			case world.GameModeAdventure:
				return "adventure"
			case world.GameModeSpectator:
				return "spectator"
			default:
				return "unknown"
			}
		},
		// Vuelo
		"startFlying": func() { p.StartFlying() },
		"stopFlying":  func() { p.StopFlying() },
		// Efectos visuales
		"setInvisible": func() { p.SetInvisible() },
		"setVisible":   func() { p.SetVisible() },
		"isInvisible":  func() bool { return p.Invisible() },
		// Velocidad
		"getSpeed": func() float64 { return p.Speed() },
		"setSpeed": func(s float64) { p.SetSpeed(s) },
		// Inventario
		"giveItem": func(itemName string, count int) bool {
			it, ok := world.ItemByName(itemName, 0)
			if !ok {
				fmt.Printf("[Commands] giveItem: item desconocido '%s'\n", itemName)
				return false
			}
			stack := dfitem.NewStack(it, count)
			_, err := p.Inventory().AddItem(stack)
			if err != nil {
				fmt.Printf("[Commands] giveItem: inventario lleno para '%s'\n", itemName)
				return false
			}
			return true
		},
		"clearInventory": func() { p.Inventory().Clear() },
		"getItemCount": func(itemName string) int {
			it, ok := world.ItemByName(itemName, 0)
			if !ok {
				return 0
			}
			searchStack := dfitem.NewStack(it, 1)
			total := 0
			for _, s := range p.Inventory().Items() {
				if s.Comparable(searchStack) {
					total += s.Count()
				}
			}
			return total
		},
		// Sonidos
		"playSound": func(soundName string) {
			var s world.Sound
			switch soundName {
			case "click":
				s = sound.Click{}
			case "levelup":
				s = sound.LevelUp{}
			case "pop":
				s = sound.Pop{}
			case "burp":
				s = sound.Burp{}
			case "door_open":
				s = sound.DoorOpen{}
			case "door_close":
				s = sound.DoorClose{}
			case "chest_open":
				s = sound.ChestOpen{}
			case "chest_close":
				s = sound.ChestClose{}
			case "anvil_land":
				s = sound.AnvilLand{}
			case "bow_shoot":
				s = sound.BowShoot{}
			case "deny":
				s = sound.Deny{}
			case "arrow_hit":
				s = sound.ArrowHit{}
			default:
				fmt.Printf("[Commands] playSound: sonido desconocido '%s'\n", soundName)
				return
			}
			p.PlaySound(s)
		},
		// Inventario del jugador
		"getInventory": func() interface{} {
			return buildInventoryMap(p.Inventory(), "player")
		},
		// Comandos
		"executeCommand": func(cmd string) { p.ExecuteCommand(cmd) },
		// Efectos de poción
		"addEffect": func(name string, level int, seconds int) {
			t, ok := effectTypeByName(name)
			if !ok {
				fmt.Printf("[Commands] addEffect: efecto desconocido '%s'\n", name)
				return
			}
			p.AddEffect(effect.New(t, level, time.Duration(seconds)*time.Second))
		},
		"removeEffect": func(name string) {
			t, ok := effectTypeByName(name)
			if !ok {
				fmt.Printf("[Commands] removeEffect: efecto desconocido '%s'\n", name)
				return
			}
			p.RemoveEffect(t)
		},
		"clearEffects": func() {
			for _, e := range p.Effects() {
				p.RemoveEffect(e.Type())
			}
		},
		// Armadura
		"setArmour": func(slot int, itemName string) {
			it, ok := world.ItemByName(itemName, 0)
			if !ok {
				fmt.Printf("[Commands] setArmour: item desconocido '%s'\n", itemName)
				return
			}
			stack := dfitem.NewStack(it, 1)
			switch slot {
			case 0:
				p.Armour().SetHelmet(stack)
			case 1:
				p.Armour().SetChestplate(stack)
			case 2:
				p.Armour().SetLeggings(stack)
			case 3:
				p.Armour().SetBoots(stack)
			default:
				fmt.Printf("[Commands] setArmour: slot inválido %d\n", slot)
			}
		},
		"getArmour": func(slot int) string {
			var stack dfitem.Stack
			switch slot {
			case 0:
				stack = p.Armour().Helmet()
			case 1:
				stack = p.Armour().Chestplate()
			case 2:
				stack = p.Armour().Leggings()
			case 3:
				stack = p.Armour().Boots()
			default:
				return ""
			}
			if stack.Empty() {
				return ""
			}
			name, _ := stack.Item().EncodeItem()
			return name
		},
		"clearArmour": func() {
			p.Armour().Set(dfitem.Stack{}, dfitem.Stack{}, dfitem.Stack{}, dfitem.Stack{})
		},
		// Scoreboard
		"sendScoreboard": func(sbVal goja.Value) {
			if sbVal == nil || goja.IsNull(sbVal) || goja.IsUndefined(sbVal) {
				fmt.Printf("[Commands] sendScoreboard: scoreboard nulo\n")
				return
			}
			obj, ok := sbVal.(*goja.Object)
			if !ok {
				fmt.Printf("[Commands] sendScoreboard: argumento no es un objeto\n")
				return
			}
			raw := obj.Get("_board")
			if raw == nil {
				fmt.Printf("[Commands] sendScoreboard: objeto no tiene _board (¿usaste scoreboard.create()?)\n")
				return
			}
			board, ok := raw.Export().(*scoreboard.Scoreboard)
			if !ok {
				fmt.Printf("[Commands] sendScoreboard: _board no es un *scoreboard.Scoreboard\n")
				return
			}
			p.SendScoreboard(board)
		},
		"removeScoreboard": func() {
			p.RemoveScoreboard()
		},
		// Referencia interna al jugador Go — usada por scoreboardManager y liveScoreboard.
		"_player": p,
	}
}

// effectTypeByName retorna el tipo de efecto de Dragonfly dado su nombre en string.
// IMPORTANTE: debe mantenerse sincronizado con effectTypeByName en loader/loader.go.
func effectTypeByName(name string) (effect.LastingType, bool) {
	switch name {
	case "speed":
		return effect.Speed, true
	case "slowness":
		return effect.Slowness, true
	case "haste":
		return effect.Haste, true
	case "mining_fatigue":
		return effect.MiningFatigue, true
	case "strength":
		return effect.Strength, true
	case "jump_boost":
		return effect.JumpBoost, true
	case "nausea":
		return effect.Nausea, true
	case "regeneration":
		return effect.Regeneration, true
	case "resistance":
		return effect.Resistance, true
	case "fire_resistance":
		return effect.FireResistance, true
	case "water_breathing":
		return effect.WaterBreathing, true
	case "invisibility":
		return effect.Invisibility, true
	case "blindness":
		return effect.Blindness, true
	case "night_vision":
		return effect.NightVision, true
	case "hunger":
		return effect.Hunger, true
	case "weakness":
		return effect.Weakness, true
	case "poison":
		return effect.Poison, true
	case "wither":
		return effect.Wither, true
	case "health_boost":
		return effect.HealthBoost, true
	case "absorption":
		return effect.Absorption, true
	case "saturation":
		return effect.Saturation, true
	case "levitation":
		return effect.Levitation, true
	case "slow_falling":
		return effect.SlowFalling, true
	case "conduit_power":
		return effect.ConduitPower, true
	case "darkness":
		return effect.Darkness, true
	default:
		return nil, false
	}
}

// MakeJSCallback convierte un callable de Goja en un jsCommandCallback de Go.
// El callback JS recibe: (player, args) donde args es un array JS.
// srv puede ser nil (en tests), en cuyo caso server.getPlayer() retorna null.
func MakeJSCallback(vm *goja.Runtime, callback goja.Callable, pluginName string, srv *server.Server) jsCommandCallback {
	return func(playerMap map[string]interface{}, args []string, tx *world.Tx) {
		jsPlayer := vm.ToValue(playerMap)

		// Convertir []string a array JS
		jsArgs := vm.NewArray()
		for i, arg := range args {
			jsArgs.Set(fmt.Sprintf("%d", i), vm.ToValue(arg))
		}

		// Exponer server y world con el tx activo (sin deadlock)
		vm.Set("server", BuildServerMapFromTx(vm, srv, tx))
		vm.Set("world", BuildWorldMapFromTx(vm, srv, tx))

		if _, err := callback(goja.Undefined(), jsPlayer, jsArgs); err != nil {
			fmt.Printf("[%s] Error en handler de comando: %v\n", pluginName, err)
		}
	}
}

// inventoryTypeFromBlock retorna el tipo de inventario como string.
// Debe mantenerse sincronizado con inventoryTypeFromBlock en loader/loader.go.
func inventoryTypeFromBlock(b world.Block) string {
	switch b.(type) {
	case block.Chest:
		return "chest"
	case block.Barrel:
		return "barrel"
	case block.Hopper:
		return "hopper"
	case block.Furnace:
		return "furnace"
	case block.BlastFurnace:
		return "blast_furnace"
	case block.Smoker:
		return "smoker"
	case block.BrewingStand:
		return "brewing_stand"
	default:
		return "container"
	}
}

// enchantmentTypeByName retorna un EnchantmentType dado su nombre.
func enchantmentTypeByName(name string) (dfitem.EnchantmentType, bool) {
	switch strings.ToLower(name) {
	case "sharpness":
		return enchantment.Sharpness, true
	case "efficiency":
		return enchantment.Efficiency, true
	case "unbreaking":
		return enchantment.Unbreaking, true
	case "silk_touch":
		return enchantment.SilkTouch, true
	case "power":
		return enchantment.Power, true
	case "punch":
		return enchantment.Punch, true
	case "flame":
		return enchantment.Flame, true
	case "infinity":
		return enchantment.Infinity, true
	case "protection":
		return enchantment.Protection, true
	case "fire_protection":
		return enchantment.FireProtection, true
	case "blast_protection":
		return enchantment.BlastProtection, true
	case "projectile_protection":
		return enchantment.ProjectileProtection, true
	case "feather_falling":
		return enchantment.FeatherFalling, true
	case "thorns":
		return enchantment.Thorns, true
	case "respiration":
		return enchantment.Respiration, true
	case "aqua_affinity":
		return enchantment.AquaAffinity, true
	case "depth_strider":
		return enchantment.DepthStrider, true
	case "mending":
		return enchantment.Mending, true
	case "vanishing":
		return enchantment.CurseOfVanishing, true
	case "fire_aspect":
		return enchantment.FireAspect, true
	case "knockback":
		return enchantment.Knockback, true
	default:
		return nil, false
	}
}

func enchantmentName(t dfitem.EnchantmentType) string {
	name := strings.ToLower(t.Name())
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

func extractStackFromJS(value interface{}) (dfitem.Stack, bool) {
	switch v := value.(type) {
	case *goja.Object:
		value = v
	case goja.Value:
		if obj, ok := v.(*goja.Object); ok {
			value = obj
		} else {
			value = v.Export()
		}
	}
	if m, ok := value.(map[string]interface{}); ok {
		if s, ok := m["_stack"].(dfitem.Stack); ok {
			return s, true
		}
	}
	if obj, ok := value.(*goja.Object); ok {
		if raw := obj.Get("_stack"); raw != nil {
			if s, ok := raw.Export().(dfitem.Stack); ok {
				return s, true
			}
		}
	}
	if obj, ok := value.(*goja.Object); ok {
		name := callString(obj, "getName")
		count := callInt(obj, "getCount", 1)
		if name == "" {
			return dfitem.Stack{}, false
		}
		it, ok := world.ItemByName(name, 0)
		if !ok {
			return dfitem.Stack{}, false
		}
		stack := dfitem.NewStack(it, count)
		if display := callString(obj, "getDisplayName"); display != "" {
			stack = stack.WithCustomName(display)
		}
		if lore := callStringSlice(obj, "getLore"); len(lore) > 0 {
			stack = stack.WithLore(lore...)
		}
		if dur := callInt(obj, "getDurability", 0); stack.MaxDurability() != -1 {
			stack = stack.WithDurability(dur)
		}
		if ench := callEnchantments(obj); len(ench) > 0 {
			stack = stack.WithEnchantments(ench...)
		}
		return stack, true
	}
	return dfitem.Stack{}, false
}

func callString(obj *goja.Object, method string) string {
	if fnVal := obj.Get(method); fnVal != nil {
		if fn, ok := goja.AssertFunction(fnVal); ok {
			res, err := fn(obj)
			if err == nil {
				return res.String()
			}
		}
	}
	return ""
}

func callInt(obj *goja.Object, method string, def int) int {
	if fnVal := obj.Get(method); fnVal != nil {
		if fn, ok := goja.AssertFunction(fnVal); ok {
			res, err := fn(obj)
			if err == nil {
				return int(res.ToInteger())
			}
		}
	}
	return def
}

func callStringSlice(obj *goja.Object, method string) []string {
	if fnVal := obj.Get(method); fnVal != nil {
		if fn, ok := goja.AssertFunction(fnVal); ok {
			res, err := fn(obj)
			if err == nil {
				if arr, ok := res.Export().([]interface{}); ok {
					out := make([]string, 0, len(arr))
					for _, v := range arr {
						out = append(out, fmt.Sprint(v))
					}
					return out
				}
			}
		}
	}
	return nil
}

func callEnchantments(obj *goja.Object) []dfitem.Enchantment {
	if fnVal := obj.Get("getEnchantments"); fnVal != nil {
		if fn, ok := goja.AssertFunction(fnVal); ok {
			res, err := fn(obj)
			if err == nil {
				if arr, ok := res.Export().([]interface{}); ok {
					var out []dfitem.Enchantment
					for _, v := range arr {
						m, ok := v.(map[string]interface{})
						if !ok {
							continue
						}
						name, _ := m["name"].(string)
						lvlVal := 0
						switch v := m["level"].(type) {
						case int64:
							lvlVal = int(v)
						case int:
							lvlVal = v
						case float64:
							lvlVal = int(v)
						}
						if et, ok := enchantmentTypeByName(name); ok {
							out = append(out, dfitem.NewEnchantment(et, lvlVal))
						}
					}
					return out
				}
			}
		}
	}
	return nil
}

func newItemWrapper(stack dfitem.Stack) map[string]interface{} {
	return map[string]interface{}{
		"getName": func() string {
			name, _ := stack.Item().EncodeItem()
			return name
		},
		"getCount": func() int { return stack.Count() },
		"setCount": func(count int) map[string]interface{} {
			if count < 1 {
				count = 1
			}
			newStack := dfitem.NewStack(stack.Item(), count)
			newStack = newStack.WithCustomName(stack.CustomName())
			newStack = newStack.WithLore(stack.Lore()...)
			newStack = newStack.WithEnchantments(stack.Enchantments()...)
			newStack = newStack.WithAnvilCost(stack.AnvilCost())
			if stack.Unbreakable() {
				newStack = newStack.AsUnbreakable()
			}
			if stack.MaxDurability() != -1 {
				newStack = newStack.WithDurability(stack.Durability())
			}
			return newItemWrapper(newStack)
		},
		"getDisplayName": func() string { return stack.CustomName() },
		"setDisplayName": func(name string) map[string]interface{} {
			return newItemWrapper(stack.WithCustomName(name))
		},
		"getLore": func() []string { return stack.Lore() },
		"setLore": func(lines []string) map[string]interface{} {
			return newItemWrapper(stack.WithLore(lines...))
		},
		"getDurability": func() int { return stack.Durability() },
		"setDurability": func(d int) map[string]interface{} { return newItemWrapper(stack.WithDurability(d)) },
		"getMaxDurability": func() int { return stack.MaxDurability() },
		"isEnchanted": func() bool { return len(stack.Enchantments()) > 0 },
		"getEnchantments": func() []map[string]interface{} {
			var res []map[string]interface{}
			for _, ench := range stack.Enchantments() {
				name := enchantmentName(ench.Type())
				res = append(res, map[string]interface{}{
					"name":  name,
					"level": ench.Level(),
				})
			}
			return res
		},
		"addEnchantment": func(name string, level int) map[string]interface{} {
			et, ok := enchantmentTypeByName(name)
			if !ok {
				return newItemWrapper(stack)
			}
			return newItemWrapper(stack.WithEnchantments(dfitem.NewEnchantment(et, level)))
		},
		"removeEnchantment": func(name string) map[string]interface{} {
			et, ok := enchantmentTypeByName(name)
			if !ok {
				return newItemWrapper(stack)
			}
			return newItemWrapper(stack.WithoutEnchantments(et))
		},
		"clone": func() map[string]interface{} { return newItemWrapper(stack) },
		"_stack": stack,
	}
}

// buildInventoryMap construye un mapa JS para interactuar con un inventario.
// Debe mantenerse sincronizado con newInventoryWrapper en loader/loader.go.
func buildInventoryMap(inv *inventory.Inventory, invType string) map[string]interface{} {
	return map[string]interface{}{
		"getType": func() string { return invType },
		"getSize": func() int { return inv.Size() },
		"getItem": func(slot int) interface{} {
			stack, err := inv.Item(slot)
			if err != nil || stack.Empty() {
				return nil
			}
			return newItemWrapper(stack)
		},
		"setItem": func(slot int, itemName string, count int) bool {
			it, ok := world.ItemByName(itemName, 0)
			if !ok {
				return false
			}
			stack := dfitem.NewStack(it, count)
			return inv.SetItem(slot, stack) == nil
		},
		"setItemStack": func(slot int, itemObj goja.Value) bool {
			if raw, ok := extractStackFromJS(itemObj); ok {
				return inv.SetItem(slot, raw) == nil
			}
			return false
		},
		"addItem": func(itemName string, count int) int {
			it, ok := world.ItemByName(itemName, 0)
			if !ok {
				return count
			}
			stack := dfitem.NewStack(it, count)
			added, _ := inv.AddItem(stack)
			return count - added
		},
		"addItemStack": func(itemObj goja.Value) int {
			if raw, ok := extractStackFromJS(itemObj); ok {
				added, _ := inv.AddItem(raw)
				return raw.Count() - added
			}
			return 0
		},
		"removeItem": func(itemName string, count int) bool {
			it, ok := world.ItemByName(itemName, 0)
			if !ok {
				return false
			}
			stack := dfitem.NewStack(it, count)
			return inv.RemoveItem(stack) == nil
		},
		"clear": func() { inv.Clear() },
		"contains": func(itemName string) bool {
			it, ok := world.ItemByName(itemName, 0)
			if !ok {
				return false
			}
			return inv.ContainsItem(dfitem.NewStack(it, 1))
		},
		"getItems": func() []interface{} {
			var items []interface{}
			for i := 0; i < inv.Size(); i++ {
				stack, err := inv.Item(i)
				if err != nil || stack.Empty() {
					continue
				}
				items = append(items, map[string]interface{}{
					"slot": i,
					"item": newItemWrapper(stack),
				})
			}
			return items
		},
		"count": func(itemName string) int {
			it, ok := world.ItemByName(itemName, 0)
			if !ok {
				return 0
			}
			total := 0
			for i := 0; i < inv.Size(); i++ {
				stack, err := inv.Item(i)
				if err != nil || stack.Empty() {
					continue
				}
				if stack.Item() == it {
					total += stack.Count()
				}
			}
			return total
		},
		"setContents": func(items []interface{}) bool {
			inv.Clear()
			for _, item := range items {
				m, ok := item.(map[string]interface{})
				if !ok {
					continue
				}
				slot, _ := m["slot"].(int64)
				if rawItem, ok := m["item"]; ok {
					if stack, ok := extractStackFromJS(rawItem); ok {
						_ = inv.SetItem(int(slot), stack)
						continue
					}
				}
				name, _ := m["name"].(string)
				countVal, _ := m["count"].(int64)
				if name == "" {
					continue
				}
				it, ok := world.ItemByName(name, 0)
				if !ok {
					fmt.Printf("[Commands] setContents: item desconocido '%s'\n", name)
					continue
				}
				stack := dfitem.NewStack(it, int(countVal))
				_ = inv.SetItem(int(slot), stack)
			}
			return true
		},
	}
}

// BuildEntityMap construye un mapa JS con métodos para interactuar con una entidad.
// Debe mantenerse sincronizado con newEntityWrapper en loader/loader.go.
func BuildEntityMap(e world.Entity, tx *world.Tx) map[string]interface{} {
	m := map[string]interface{}{
		"getUUID": func() string { return e.H().UUID().String() },
		"getType": func() string { return e.H().Type().EncodeEntity() },
		"getX":    func() float64 { return e.Position().X() },
		"getY":    func() float64 { return e.Position().Y() },
		"getZ":    func() float64 { return e.Position().Z() },
		"remove":  func() { tx.RemoveEntity(e) },
		"teleport": func(x, y, z float64) {
			if mover, ok := e.(interface{ Teleport(mgl64.Vec3) }); ok {
				mover.Teleport(mgl64.Vec3{x, y, z})
			}
		},
		"setVelocity": func(x, y, z float64) {
			if mover, ok := e.(interface{ SetVelocity(mgl64.Vec3) }); ok {
				mover.SetVelocity(mgl64.Vec3{x, y, z})
			}
		},
	}
	if living, ok := e.(entity.Living); ok {
		m["getHealth"] = func() float64 { return living.Health() }
		m["getMaxHealth"] = func() float64 { return living.MaxHealth() }
		m["setMaxHealth"] = func(h float64) { living.SetMaxHealth(h) }
		m["isDead"] = func() bool { return living.Dead() }
		m["hurt"] = func(damage float64) { living.Hurt(damage, entity.AttackDamageSource{}) }
		m["heal"] = func(health float64) { living.Heal(health, entity.FoodHealingSource{}) }
		m["knockBack"] = func(x, y, z, force, height float64) {
			living.KnockBack(mgl64.Vec3{x, y, z}, force, height)
		}
		m["addEffect"] = func(name string, level int, seconds int) {
			t, ok := effectTypeByName(name)
			if !ok {
				return
			}
			living.AddEffect(effect.New(t, level, time.Duration(seconds)*time.Second))
		}
		m["removeEffect"] = func(name string) {
			t, ok := effectTypeByName(name)
			if !ok {
				return
			}
			living.RemoveEffect(t)
		}
		m["clearEffects"] = func() {
			for _, eff := range living.Effects() {
				living.RemoveEffect(eff.Type())
			}
		}
		m["getSpeed"] = func() float64 { return living.Speed() }
	}
	if pl, ok := e.(*dfplayer.Player); ok {
		m["getName"] = func() string { return pl.Name() }
		m["sendMessage"] = func(msg string) { pl.Message(msg) }
		m["sendTitle"] = func(text, subtitle string) {
			t := title.New(text).WithSubtitle(subtitle)
			pl.SendTitle(t)
		}
		m["disconnect"] = func(msg string) { pl.Disconnect(msg) }
	}
	return m
}

// BuildWorldMapFromTx construye el objeto world usando el tx activo.
// Exportado para ser usado desde loader.go en eventos.
func BuildWorldMapFromTx(vm *goja.Runtime, srv *server.Server, tx *world.Tx) map[string]interface{} {
	return map[string]interface{}{
		// --- Inventarios de bloques ---
		"getInventory": func(x, y, z int) interface{} {
			if srv == nil || tx == nil {
				return nil
			}
			b := tx.Block(cube.Pos{x, y, z})
			type container interface {
				Inventory(*world.Tx, cube.Pos) *inventory.Inventory
			}
			if c, ok := b.(container); ok {
				invType := inventoryTypeFromBlock(b)
				return buildInventoryMap(c.Inventory(tx, cube.Pos{x, y, z}), invType)
			}
			return nil
		},
		// --- Bloques ---
		"setBlock": func(x, y, z int, blockName string) {
			if srv == nil || tx == nil {
				return
			}
			b, ok := world.BlockByName(blockName, nil)
			if !ok {
				fmt.Printf("[Commands] world.setBlock: bloque desconocido '%s'\n", blockName)
				return
			}
			tx.SetBlock(cube.Pos{x, y, z}, b, nil)
		},
		"getBlock": func(x, y, z int) string {
			if srv == nil || tx == nil {
				return ""
			}
			b := tx.Block(cube.Pos{x, y, z})
			name, _ := b.EncodeBlock()
			return name
		},
		"getHighestBlock": func(x, z int) int {
			if srv == nil || tx == nil {
				return 0
			}
			return tx.HighestBlock(x, z)
		},
		// --- Entidades ---
		"getEntities": func() []interface{} {
			if srv == nil || tx == nil {
				return []interface{}{}
			}
			var entities []interface{}
			for e := range tx.Entities() {
				entities = append(entities, BuildEntityMap(e, tx))
			}
			return entities
		},
		"getEntitiesInRadius": func(x, y, z, radius float64) []interface{} {
			if srv == nil || tx == nil {
				return []interface{}{}
			}
			center := mgl64.Vec3{x, y, z}
			box := cube.Box(x-radius, y-radius, z-radius, x+radius, y+radius, z+radius)
			var entities []interface{}
			for e := range tx.EntitiesWithin(box) {
				if e.Position().Sub(center).Len() <= radius {
					entities = append(entities, BuildEntityMap(e, tx))
				}
			}
			return entities
		},
		"removeEntityByUUID": func(uuidStr string) bool {
			if srv == nil || tx == nil {
				return false
			}
			for e := range tx.Entities() {
				if e.H().UUID().String() == uuidStr {
					tx.RemoveEntity(e)
					return true
				}
			}
			return false
		},
		// --- Partículas ---
		"spawnParticle": func(x, y, z float64, particleName string) {
			// Las partículas no requieren tx especial, se delegan al loader
		},
		// --- Spawn de entidades ---
		"spawnEntity": func(entityType string, x, y, z float64, opts goja.Value) {
			if srv == nil || tx == nil {
				return
			}
			pos := mgl64.Vec3{x, y, z}
			spawnOpts := world.EntitySpawnOpts{Position: pos}

			var optsMap map[string]goja.Value
			if opts != nil && !goja.IsUndefined(opts) && !goja.IsNull(opts) {
				if obj, ok := opts.(*goja.Object); ok {
					optsMap = make(map[string]goja.Value)
					for _, key := range obj.Keys() {
						optsMap[key] = obj.Get(key)
					}
				}
			}
			getFloat := func(key string, def float64) float64 {
				if v, ok := optsMap[key]; ok {
					return v.ToFloat()
				}
				return def
			}
			getString := func(key string, def string) string {
				if v, ok := optsMap[key]; ok {
					return v.String()
				}
				return def
			}
			getInt := func(key string, def int) int {
				if v, ok := optsMap[key]; ok {
					return int(v.ToInteger())
				}
				return def
			}

			switch entityType {
			case "lightning":
				tx.AddEntity(entity.NewLightning(spawnOpts))
			case "tnt":
				fuse := time.Duration(getFloat("fuse", 4) * float64(time.Second))
				tx.AddEntity(entity.NewTNT(spawnOpts, fuse))
			case "text":
				text := getString("text", "")
				tx.AddEntity(entity.NewText(text, pos))
			case "experience_orb":
				amount := getInt("amount", 1)
				for _, orb := range entity.NewExperienceOrbs(pos, amount) {
					tx.AddEntity(orb)
				}
			case "item":
				itemName := getString("item", "minecraft:stone")
				count := getInt("count", 1)
				it, ok := world.ItemByName(itemName, 0)
				if !ok {
					fmt.Printf("[Commands] world.spawnEntity: item desconocido '%s'\n", itemName)
					return
				}
				stack := dfitem.NewStack(it, count)
				tx.AddEntity(entity.NewItem(spawnOpts, stack))
			default:
				fmt.Printf("[Commands] world.spawnEntity: tipo desconocido '%s'\n", entityType)
			}
		},
		// --- Jugadores (por retrocompatibilidad) ---
		"getPlayers": func() []interface{} {
			if srv == nil || tx == nil {
				return []interface{}{}
			}
			var players []interface{}
			for pl := range srv.Players(tx) {
				players = append(players, buildPlayerMap(pl))
			}
			return players
		},
		"getPlayerCount": func() int {
			if srv == nil {
				return 0
			}
			return srv.PlayerCount()
		},
		"broadcast": func(msg string) {
			if srv == nil || tx == nil {
				return
			}
			for pl := range srv.Players(tx) {
				pl.Message(msg)
			}
		},
	}
}

// BuildServerMapFromTx construye el objeto server usando el tx activo.
// Exportado para ser usado desde loader.go en eventos.
// Esto evita el deadlock que causaría ExecWorld() al abrir una nueva transacción
// cuando ya hay una activa (la del comando).
func BuildServerMapFromTx(vm *goja.Runtime, srv *server.Server, tx *world.Tx) map[string]interface{} {
	return map[string]interface{}{
		"getPlayers": func() []interface{} {
			if srv == nil || tx == nil {
				return []interface{}{}
			}
			var players []interface{}
			for pl := range srv.Players(tx) {
				players = append(players, buildPlayerMap(pl))
			}
			return players
		},
		"getPlayerCount": func() int {
			if srv == nil {
				return 0
			}
			return srv.PlayerCount()
		},
		"getMaxPlayers": func() int {
			if srv == nil {
				return 0
			}
			return srv.MaxPlayerCount()
		},
		"getPlayer": func(name string) goja.Value {
			if srv == nil || tx == nil {
				return goja.Null()
			}
			for pl := range srv.Players(tx) {
				if pl.Name() == name {
					return vm.ToValue(buildPlayerMap(pl))
				}
			}
			return goja.Null()
		},
		"getPlayerByXUID": func(xuid string) goja.Value {
			if srv == nil || tx == nil {
				return goja.Null()
			}
			for pl := range srv.Players(tx) {
				if pl.XUID() == xuid {
					return vm.ToValue(buildPlayerMap(pl))
				}
			}
			return goja.Null()
		},
		"broadcast": func(msg string) {
			if srv == nil || tx == nil {
				return
			}
			for pl := range srv.Players(tx) {
				pl.Message(msg)
			}
		},
		"broadcastTitle": func(text, subtitle string) {
			if srv == nil || tx == nil {
				return
			}
			t := title.New(text).WithSubtitle(subtitle)
			for pl := range srv.Players(tx) {
				pl.SendTitle(t)
			}
		},
		"getName": func() string {
			return ""
		},
		"shutdown": func() {
			if srv == nil {
				return
			}
			go func() {
				if err := srv.Close(); err != nil {
					fmt.Printf("[Commands] Error al cerrar servidor: %v\n", err)
				}
			}()
		},
	}
}
