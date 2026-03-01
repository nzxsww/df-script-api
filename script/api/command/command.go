// Package command implementa el sistema de comandos para plugins JavaScript.
// Permite registrar comandos desde JS que son ejecutados por jugadores en el servidor.
package command

import (
	"fmt"
	"sync"

	"github.com/df-mc/dragonfly/server/cmd"
	dfitem "github.com/df-mc/dragonfly/server/item"
	dfplayer "github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/dop251/goja"
	"github.com/go-gl/mathgl/mgl64"
)

// jsCommandCallback es el tipo del callback JS que recibe el jugador y los args.
type jsCommandCallback func(player map[string]interface{}, args []string)

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

	// Llamar el callback JS
	callback(playerObj, args)
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
		"sendToast":        func(title, msg string) { p.SendToast(title, msg) },
		"sendJukeboxPopup": func(msg string) { p.SendJukeboxPopup(msg) },
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
	}
}

// MakeJSCallback convierte un callable de Goja en un jsCommandCallback de Go.
// El callback JS recibe: (player, args) donde args es un array JS.
func MakeJSCallback(vm *goja.Runtime, callback goja.Callable, pluginName string) jsCommandCallback {
	return func(playerMap map[string]interface{}, args []string) {
		jsPlayer := vm.ToValue(playerMap)

		// Convertir []string a array JS
		jsArgs := vm.NewArray()
		for i, arg := range args {
			jsArgs.Set(fmt.Sprintf("%d", i), vm.ToValue(arg))
		}

		if _, err := callback(goja.Undefined(), jsPlayer, jsArgs); err != nil {
			fmt.Printf("[%s] Error en handler de comando: %v\n", pluginName, err)
		}
	}
}
