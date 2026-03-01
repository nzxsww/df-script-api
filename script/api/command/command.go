// Package command implementa el sistema de comandos para plugins JavaScript.
// Permite registrar comandos desde JS que son ejecutados por jugadores en el servidor.
package command

import (
	"fmt"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity/effect"
	dfitem "github.com/df-mc/dragonfly/server/item"
	dfplayer "github.com/df-mc/dragonfly/server/player"
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

		// Exponer server object con acceso al tx activo (sin deadlock)
		serverObj := buildServerMap(vm, srv, tx)
		vm.Set("server", serverObj)

		if _, err := callback(goja.Undefined(), jsPlayer, jsArgs); err != nil {
			fmt.Printf("[%s] Error en handler de comando: %v\n", pluginName, err)
		}
	}
}

// buildServerMap construye el objeto server usando el tx activo del comando.
// Esto evita el deadlock que causaría ExecWorld() al abrir una nueva transacción
// cuando ya hay una activa (la del comando).
func buildServerMap(vm *goja.Runtime, srv *server.Server, tx *world.Tx) map[string]interface{} {
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
