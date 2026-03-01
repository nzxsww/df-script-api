package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/df-mc/dragonfly/server"
	dfitem "github.com/df-mc/dragonfly/server/item"
	dfplayer "github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/dop251/goja"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/nzxsww/dragonfly-script-api/script/api/command"
	"github.com/nzxsww/dragonfly-script-api/script/api/config"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
	evplayer "github.com/nzxsww/dragonfly-script-api/script/api/event/player"
	"github.com/nzxsww/dragonfly-script-api/script/api/plugin"
	"gopkg.in/yaml.v3"
)

// Descriptor contiene los metadatos del plugin leídos desde plugin.yml.
type Descriptor struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Author      string `yaml:"author"`
	Description string `yaml:"description"`
	Main        string `yaml:"main"`
	APIVersion  string `yaml:"api-version"`
}

// playerWrapper expone una API limpia del jugador al JavaScript.
// Evita que el JS acceda directamente a métodos internos de Go.
type playerWrapper struct {
	p *dfplayer.Player
}

func newPlayerWrapper(p *dfplayer.Player) map[string]interface{} {
	w := &playerWrapper{p: p}
	return map[string]interface{}{
		// Identidad
		"getName":    w.getName,
		"getUUID":    w.getUUID,
		"getXUID":    w.getXUID,
		"getNameTag": w.getNameTag,
		"setNameTag": w.setNameTag,
		// Mensajes
		"sendMessage":      w.sendMessage,
		"sendPopup":        w.sendPopup,
		"sendTip":          w.sendTip,
		"sendToast":        w.sendToast,
		"sendJukeboxPopup": w.sendJukeboxPopup,
		// Conexión
		"disconnect": w.disconnect,
		"transfer":   w.transfer,
		"getLatency": w.getLatency,
		// Posición y movimiento
		"getX":        w.getX,
		"getY":        w.getY,
		"getZ":        w.getZ,
		"teleport":    w.teleport,
		"setVelocity": w.setVelocity,
		// Estado físico
		"getHealth":    w.getHealth,
		"getMaxHealth": w.getMaxHealth,
		"setMaxHealth": w.setMaxHealth,
		"getFoodLevel": w.getFoodLevel,
		"setFoodLevel": w.setFoodLevel,
		"isOnGround":   w.isOnGround,
		"isSneaking":   w.isSneaking,
		"isSprinting":  w.isSprinting,
		"isFlying":     w.isFlying,
		"isSwimming":   w.isSwimming,
		"isDead":       w.isDead,
		"isImmobile":   w.isImmobile,
		// Experiencia
		"getExperience":      w.getExperience,
		"getExperienceLevel": w.getExperienceLevel,
		"addExperience":      w.addExperience,
		"setExperienceLevel": w.setExperienceLevel,
		// Modo de juego
		"setGameMode": w.setGameMode,
		"getGameMode": w.getGameMode,
		// Vuelo
		"startFlying": w.startFlying,
		"stopFlying":  w.stopFlying,
		// Efectos visuales
		"setInvisible": w.setInvisible,
		"setVisible":   w.setVisible,
		"isInvisible":  w.isInvisible,
		// Velocidad
		"getSpeed": w.getSpeed,
		"setSpeed": w.setSpeed,
		// Comandos
		"executeCommand": w.executeCommand,
		// Inventario
		"giveItem":       w.giveItem,
		"clearInventory": w.clearInventory,
		"getItemCount":   w.getItemCount,
		// Sonidos
		"playSound": w.playSound,
	}
}

// --- Identidad ---
func (w *playerWrapper) getName() string    { return w.p.Name() }
func (w *playerWrapper) getUUID() string    { return w.p.UUID().String() }
func (w *playerWrapper) getXUID() string    { return w.p.XUID() }
func (w *playerWrapper) getNameTag() string { return w.p.NameTag() }
func (w *playerWrapper) setNameTag(name string) { w.p.SetNameTag(name) }

// --- Mensajes ---
func (w *playerWrapper) sendMessage(msg string)      { w.p.Message(msg) }
func (w *playerWrapper) sendPopup(msg string)        { w.p.SendPopup(msg) }
func (w *playerWrapper) sendTip(msg string)          { w.p.SendTip(msg) }
func (w *playerWrapper) sendToast(title, msg string) { w.p.SendToast(title, msg) }
func (w *playerWrapper) sendJukeboxPopup(msg string) { w.p.SendJukeboxPopup(msg) }

// --- Conexión ---
func (w *playerWrapper) disconnect(msg string) { w.p.Disconnect(msg) }
func (w *playerWrapper) transfer(address string) {
	if err := w.p.Transfer(address); err != nil {
		fmt.Printf("[playerWrapper] transfer error: %v\n", err)
	}
}
func (w *playerWrapper) getLatency() int64 { return w.p.Latency().Milliseconds() }

// --- Posición y movimiento ---
func (w *playerWrapper) getX() float64 { return w.p.Position().X() }
func (w *playerWrapper) getY() float64 { return w.p.Position().Y() }
func (w *playerWrapper) getZ() float64 { return w.p.Position().Z() }
func (w *playerWrapper) teleport(x, y, z float64) {
	w.p.Teleport(mgl64.Vec3{x, y, z})
}
func (w *playerWrapper) setVelocity(x, y, z float64) {
	w.p.SetVelocity(mgl64.Vec3{x, y, z})
}

// --- Estado físico ---
func (w *playerWrapper) getHealth() float64      { return w.p.Health() }
func (w *playerWrapper) getMaxHealth() float64   { return w.p.MaxHealth() }
func (w *playerWrapper) setMaxHealth(h float64)  { w.p.SetMaxHealth(h) }
func (w *playerWrapper) getFoodLevel() int       { return w.p.Food() }
func (w *playerWrapper) setFoodLevel(f int)      { w.p.SetFood(f) }
func (w *playerWrapper) isOnGround() bool        { return w.p.OnGround() }
func (w *playerWrapper) isSneaking() bool        { return w.p.Sneaking() }
func (w *playerWrapper) isSprinting() bool       { return w.p.Sprinting() }
func (w *playerWrapper) isFlying() bool          { return w.p.Flying() }
func (w *playerWrapper) isSwimming() bool        { return w.p.Swimming() }
func (w *playerWrapper) isDead() bool            { return w.p.Dead() }
func (w *playerWrapper) isImmobile() bool        { return w.p.Immobile() }

// --- Experiencia ---
func (w *playerWrapper) getExperience() int          { return w.p.Experience() }
func (w *playerWrapper) getExperienceLevel() int     { return w.p.ExperienceLevel() }
func (w *playerWrapper) addExperience(amount int)    { w.p.AddExperience(amount) }
func (w *playerWrapper) setExperienceLevel(lvl int)  { w.p.SetExperienceLevel(lvl) }

// --- Modo de juego ---
// setGameMode acepta: "survival", "creative", "adventure", "spectator"
func (w *playerWrapper) setGameMode(mode string) {
	switch mode {
	case "survival":
		w.p.SetGameMode(world.GameModeSurvival)
	case "creative":
		w.p.SetGameMode(world.GameModeCreative)
	case "adventure":
		w.p.SetGameMode(world.GameModeAdventure)
	case "spectator":
		w.p.SetGameMode(world.GameModeSpectator)
	default:
		fmt.Printf("[playerWrapper] setGameMode: modo desconocido '%s'\n", mode)
	}
}
func (w *playerWrapper) getGameMode() string {
	switch w.p.GameMode() {
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
}

// --- Vuelo ---
func (w *playerWrapper) startFlying() { w.p.StartFlying() }
func (w *playerWrapper) stopFlying()  { w.p.StopFlying() }

// --- Efectos visuales ---
func (w *playerWrapper) setInvisible() { w.p.SetInvisible() }
func (w *playerWrapper) setVisible()   { w.p.SetVisible() }
func (w *playerWrapper) isInvisible() bool { return w.p.Invisible() }

// --- Velocidad ---
func (w *playerWrapper) getSpeed() float64    { return w.p.Speed() }
func (w *playerWrapper) setSpeed(s float64)   { w.p.SetSpeed(s) }

// --- Comandos ---
func (w *playerWrapper) executeCommand(cmd string) { w.p.ExecuteCommand(cmd) }

// --- Inventario ---

// giveItem agrega un item al inventario del jugador por nombre de item de Minecraft.
// nombre: nombre del item (ej: "minecraft:diamond", "minecraft:stone")
// cantidad: cuántos items agregar (1-64)
// Retorna true si se agregó correctamente, false si el item no existe o el inventario está lleno.
func (w *playerWrapper) giveItem(itemName string, count int) bool {
	it, ok := world.ItemByName(itemName, 0)
	if !ok {
		fmt.Printf("[playerWrapper] giveItem: item desconocido '%s'\n", itemName)
		return false
	}
	stack := dfitem.NewStack(it, count)
	_, err := w.p.Inventory().AddItem(stack)
	if err != nil {
		fmt.Printf("[playerWrapper] giveItem: inventario lleno para '%s'\n", itemName)
		return false
	}
	return true
}

// clearInventory limpia completamente el inventario del jugador.
func (w *playerWrapper) clearInventory() {
	w.p.Inventory().Clear()
}

// getItemCount retorna la cantidad total de un item específico en el inventario.
func (w *playerWrapper) getItemCount(itemName string) int {
	it, ok := world.ItemByName(itemName, 0)
	if !ok {
		return 0
	}
	searchStack := dfitem.NewStack(it, 1)
	total := 0
	for _, s := range w.p.Inventory().Items() {
		if s.Comparable(searchStack) {
			total += s.Count()
		}
	}
	return total
}

// --- Sonidos ---

// playSound reproduce un sonido al jugador en su posición actual.
// Sonidos disponibles: "click", "levelup", "pop", "burp", "orb", "door_open",
// "door_close", "chest_open", "chest_close", "anvil_land", "bow_shoot"
func (w *playerWrapper) playSound(soundName string) {
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
		fmt.Printf("[playerWrapper] playSound: sonido desconocido '%s'\n", soundName)
		return
	}
	w.p.PlaySound(s)
}

// ScriptPlugin representa un plugin cargado desde JavaScript.
type ScriptPlugin struct {
	name       string
	version    string
	author     string
	dataFolder string
	srv        *server.Server
	pluginMgr  *plugin.Manager
	enabled    bool
	vm         *goja.Runtime
	module     *goja.Object
	cfg        *config.Config
}

func newScriptPlugin(desc *Descriptor, dataFolder string, srv *server.Server, mgr *plugin.Manager) *ScriptPlugin {
	return &ScriptPlugin{
		name:       desc.Name,
		version:    desc.Version,
		author:     desc.Author,
		dataFolder: dataFolder,
		srv:        srv,
		pluginMgr:  mgr,
	}
}

func (p *ScriptPlugin) OnEnable() {
	p.enabled = true
	if p.module != nil {
		if enable, ok := goja.AssertFunction(p.module.Get("onEnable")); ok {
			if _, err := enable(p.module); err != nil {
				fmt.Printf("[%s] Error en onEnable: %v\n", p.name, err)
			}
		}
	}
}

func (p *ScriptPlugin) OnDisable() {
	p.enabled = false
	if p.module != nil {
		if disable, ok := goja.AssertFunction(p.module.Get("onDisable")); ok {
			if _, err := disable(p.module); err != nil {
				fmt.Printf("[%s] Error en onDisable: %v\n", p.name, err)
			}
		}
	}
}

func (p *ScriptPlugin) GetName() string           { return p.name }
func (p *ScriptPlugin) GetDataFolder() string     { return p.dataFolder }
func (p *ScriptPlugin) GetServer() *server.Server { return p.srv }
func (p *ScriptPlugin) IsEnabled() bool           { return p.enabled }

func (p *ScriptPlugin) GetConfig() *config.Config {
	if p.cfg == nil {
		p.cfg = config.New(p.dataFolder+"/config.yml", p)
		p.cfg.SetDefaults(map[string]interface{}{})
		p.cfg.Load()
	}
	return p.cfg
}

// Loader carga plugins JavaScript desde un directorio.
type Loader struct {
	pluginDir string
	srv       *server.Server
	pluginMgr *plugin.Manager
	plugins   []*ScriptPlugin
}

func NewLoader(pluginDir string, srv *server.Server, mgr *plugin.Manager) *Loader {
	return &Loader{
		pluginDir: pluginDir,
		srv:       srv,
		pluginMgr: mgr,
	}
}

// LoadAll escanea el directorio de plugins y carga todos los que encuentre.
func (l *Loader) LoadAll() ([]plugin.Plugin, error) {
	if err := os.MkdirAll(l.pluginDir, 0755); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(l.pluginDir)
	if err != nil {
		return nil, err
	}

	var loaded []plugin.Plugin

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		p, err := l.loadFromDir(filepath.Join(l.pluginDir, entry.Name()))
		if err != nil {
			fmt.Printf("[Loader] Error cargando plugin '%s': %v\n", entry.Name(), err)
			continue
		}
		if p != nil {
			loaded = append(loaded, p)
		}
	}

	return loaded, nil
}

func (l *Loader) loadFromDir(dir string) (plugin.Plugin, error) {
	pluginYml := filepath.Join(dir, "plugin.yml")
	if _, err := os.Stat(pluginYml); os.IsNotExist(err) {
		return nil, fmt.Errorf("plugin.yml no encontrado en %s", dir)
	}

	desc, err := l.loadDescriptor(pluginYml)
	if err != nil {
		return nil, err
	}

	// Determinar el archivo de script principal
	scriptFile := filepath.Join(dir, desc.Main)
	if desc.Main == "" {
		scriptFile = filepath.Join(dir, "index.js")
	}
	if _, err := os.Stat(scriptFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("archivo principal '%s' no encontrado", scriptFile)
	}

	fmt.Printf("[Loader] Cargando plugin: %s v%s (autor: %s)\n", desc.Name, desc.Version, desc.Author)

	p := newScriptPlugin(desc, dir, l.srv, l.pluginMgr)

	if err := l.loadScript(p, scriptFile); err != nil {
		return nil, fmt.Errorf("error ejecutando script: %v", err)
	}

	l.pluginMgr.AddPlugin(p)
	l.plugins = append(l.plugins, p)

	return p, nil
}

func (l *Loader) loadDescriptor(path string) (*Descriptor, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var desc Descriptor
	if err := yaml.Unmarshal(data, &desc); err != nil {
		return nil, err
	}

	if desc.Name == "" {
		return nil, fmt.Errorf("plugin.yml debe tener un campo 'name'")
	}

	return &desc, nil
}

func (l *Loader) loadScript(p *ScriptPlugin, scriptFile string) error {
	vm := goja.New()
	p.vm = vm

	// Exponer la API al script
	l.registerConsole(vm, p)
	l.registerTimers(vm, vm)
	l.registerPluginInfo(vm, p)
	l.registerEvents(vm, p)
	l.registerCommands(vm, p)
	l.registerConfig(vm, p)

	// Leer y ejecutar el script
	scriptContent, err := os.ReadFile(scriptFile)
	if err != nil {
		return err
	}

	if _, err = vm.RunString(string(scriptContent)); err != nil {
		return fmt.Errorf("error de script: %v", err)
	}

	// Obtener el objeto module si fue exportado
	if mod := vm.Get("module"); mod != nil && !goja.IsUndefined(mod) && !goja.IsNull(mod) {
		p.module = mod.ToObject(vm)
	}

	return nil
}

// registerConsole expone console.log/warn/error/info al JavaScript.
func (l *Loader) registerConsole(vm *goja.Runtime, p *ScriptPlugin) {
	prefix := fmt.Sprintf("[%s]", p.name)
	vm.Set("console", map[string]interface{}{
		"log": func(args ...interface{}) {
			fmt.Println(append([]interface{}{prefix}, args...)...)
		},
		"info": func(args ...interface{}) {
			fmt.Println(append([]interface{}{prefix + " [INFO]"}, args...)...)
		},
		"warn": func(args ...interface{}) {
			fmt.Println(append([]interface{}{prefix + " [WARN]"}, args...)...)
		},
		"error": func(args ...interface{}) {
			fmt.Println(append([]interface{}{prefix + " [ERROR]"}, args...)...)
		},
	})
}

// registerTimers expone setTimeout y setInterval con la firma correcta de JS: fn, delay.
func (l *Loader) registerTimers(vm *goja.Runtime, _ *goja.Runtime) {
	vm.Set("setTimeout", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return goja.Undefined()
		}
		fn := call.Argument(0)
		delay := time.Duration(call.Argument(1).ToInteger()) * time.Millisecond
		go func() {
			time.Sleep(delay)
			if f, ok := goja.AssertFunction(fn); ok {
				if _, err := f(goja.Undefined()); err != nil {
					fmt.Printf("setTimeout error: %v\n", err)
				}
			}
		}()
		return goja.Undefined()
	})

	vm.Set("setInterval", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return goja.Undefined()
		}
		fn := call.Argument(0)
		delay := time.Duration(call.Argument(1).ToInteger()) * time.Millisecond
		ticker := time.NewTicker(delay)
		go func() {
			for range ticker.C {
				if f, ok := goja.AssertFunction(fn); ok {
					if _, err := f(goja.Undefined()); err != nil {
						fmt.Printf("setInterval error: %v\n", err)
						ticker.Stop()
						return
					}
				}
			}
		}()
		// Retorna una función para cancelar el interval (clearInterval)
		return vm.ToValue(func() { ticker.Stop() })
	})

	vm.Set("clearInterval", func(call goja.FunctionCall) goja.Value {
		// setInterval retorna una función Go que detiene el ticker.
		// En JS el plugin puede llamar clearInterval(id) o simplemente id().
		// Aquí lo soportamos llamando la función si es callable.
		fn := call.Argument(0)
		if !goja.IsUndefined(fn) && !goja.IsNull(fn) {
			if f, ok := goja.AssertFunction(fn); ok {
				f(goja.Undefined()) //nolint:errcheck
			}
		}
		return goja.Undefined()
	})
}

// registerPluginInfo expone la variable `plugin` con info del plugin.
func (l *Loader) registerPluginInfo(vm *goja.Runtime, p *ScriptPlugin) {
	vm.Set("plugin", map[string]interface{}{
		"name":       p.name,
		"version":    p.version,
		"author":     p.author,
		"dataFolder": p.dataFolder,
	})
}

// registerConfig expone el objeto `config` para leer/escribir configuración YAML del plugin.
//
// Uso en JS:
//
//	config.set("welcome", "¡Hola!");
//	var msg = config.get("welcome", "Bienvenido"); // segundo arg es el valor por defecto
//	config.save();
func (l *Loader) registerConfig(vm *goja.Runtime, p *ScriptPlugin) {
	cfg := p.GetConfig()
	vm.Set("config", map[string]interface{}{
		// get(key, default) — retorna el valor o el default si no existe
		"get": func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			val := cfg.Get(key)
			if val == nil {
				// retornar el segundo argumento como default si fue provisto
				if len(call.Arguments) >= 2 {
					return call.Argument(1)
				}
				return goja.Null()
			}
			return vm.ToValue(val)
		},
		"set": func(key string, value interface{}) {
			cfg.Set(key, value)
		},
		"getString": func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			v := cfg.GetString(key)
			if v == "" && len(call.Arguments) >= 2 {
				return call.Argument(1)
			}
			return vm.ToValue(v)
		},
		"getInt": func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			v := cfg.GetInt(key)
			if v == 0 && len(call.Arguments) >= 2 {
				return call.Argument(1)
			}
			return vm.ToValue(v)
		},
		"getBool": func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			v := cfg.GetBool(key)
			if !v && len(call.Arguments) >= 2 {
				return call.Argument(1)
			}
			return vm.ToValue(v)
		},
		"getFloat": func(call goja.FunctionCall) goja.Value {
			key := call.Argument(0).String()
			v := cfg.GetFloat(key)
			if v == 0.0 && len(call.Arguments) >= 2 {
				return call.Argument(1)
			}
			return vm.ToValue(v)
		},
		"save": func() {
			if err := cfg.Save(); err != nil {
				fmt.Printf("[%s] Error guardando config: %v\n", p.name, err)
			}
		},
		"reload": func() {
			if err := cfg.Load(); err != nil {
				fmt.Printf("[%s] Error recargando config: %v\n", p.name, err)
			}
		},
		"setDefaults": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) < 1 {
				return goja.Undefined()
			}
			exported := call.Argument(0).Export()
			if defaults, ok := exported.(map[string]interface{}); ok {
				cfg.SetDefaults(defaults)
			}
			return goja.Undefined()
		},
	})
}

// registerCommands expone el objeto `commands` para registrar comandos desde JS.
//
// Uso en JS:
//
//	commands.register("spawn", "Teleporta al spawn", [], function(player, args) {
//	    player.sendMessage("Teleportando al spawn...");
//	});
func (l *Loader) registerCommands(vm *goja.Runtime, p *ScriptPlugin) {
	vm.Set("commands", map[string]interface{}{
		"register": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) < 3 {
				fmt.Printf("[%s] commands.register requiere al menos 3 argumentos: nombre, descripción, callback\n", p.name)
				return goja.Undefined()
			}

			name := call.Argument(0).String()
			description := call.Argument(1).String()

			// El tercer argumento puede ser un array de aliases o el callback directamente
			var aliases []string
			var callbackArg goja.Value

			if len(call.Arguments) >= 4 {
				// commands.register("nombre", "desc", ["alias1", "alias2"], fn)
				if arr, ok := call.Argument(2).Export().([]interface{}); ok {
					for _, a := range arr {
						if s, ok := a.(string); ok {
							aliases = append(aliases, s)
						}
					}
				}
				callbackArg = call.Argument(3)
			} else {
				// commands.register("nombre", "desc", fn)
				callbackArg = call.Argument(2)
			}

			callback, ok := goja.AssertFunction(callbackArg)
			if !ok {
				fmt.Printf("[%s] commands.register: el último argumento debe ser una función\n", p.name)
				return goja.Undefined()
			}

			jsCallback := command.MakeJSCallback(vm, callback, p.name)
			command.Register(name, description, aliases, jsCallback)
			return goja.Undefined()
		},
	})
}

// registerEvents expone el objeto `events` para registrar listeners de eventos.
// Los callbacks JS reciben un objeto con métodos propios de cada evento.
func (l *Loader) registerEvents(vm *goja.Runtime, p *ScriptPlugin) {
	vm.Set("events", map[string]interface{}{
		"on": func(eventName string, callback goja.Callable) {
			l.registerEventHandler(eventName, callback, p)
		},
	})
}

// makeListener construye un RegisteredListener para registrar en una HandlerList.
func makeListener(p *ScriptPlugin, executor event.EventExecutor) event.RegisteredListener {
	return event.RegisteredListener{
		Listener: p,
		Executor: executor,
		Priority: event.PriorityNormal,
		Plugin:   p,
	}
}

// registerEventHandler conecta un callback JS al sistema de eventos para el evento dado.
func (l *Loader) registerEventHandler(eventName string, callback goja.Callable, p *ScriptPlugin) {
	vm := p.vm

	switch eventName {
	case "PlayerJoin":
		evplayer.PlayerJoinEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerJoinEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value {
				return vm.ToValue(newPlayerWrapper(ev.GetPlayer()))
			}))
			jsEvent.Set("getJoinMessage", vm.ToValue(func() string {
				return ev.GetJoinMessage()
			}))
			jsEvent.Set("setJoinMessage", vm.ToValue(func(msg string) {
				ev.SetJoinMessage(msg)
			}))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool {
				return ev.IsCancelled()
			}))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) {
				ev.SetCancelled(c)
			}))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerJoin: %v\n", p.name, err)
			}
		}))

	case "PlayerQuit":
		evplayer.PlayerQuitEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerQuitEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value {
				return vm.ToValue(newPlayerWrapper(ev.GetPlayer()))
			}))
			jsEvent.Set("getQuitMessage", vm.ToValue(func() string {
				return ev.GetQuitMessage()
			}))
			jsEvent.Set("setQuitMessage", vm.ToValue(func(msg string) {
				ev.SetQuitMessage(msg)
			}))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool {
				return ev.IsCancelled()
			}))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) {
				ev.SetCancelled(c)
			}))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerQuit: %v\n", p.name, err)
			}
		}))

	case "PlayerChat":
		evplayer.PlayerChatEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerChatEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value {
				return vm.ToValue(newPlayerWrapper(ev.GetPlayer()))
			}))
			jsEvent.Set("getMessage", vm.ToValue(func() string {
				return ev.GetMessage()
			}))
			jsEvent.Set("setMessage", vm.ToValue(func(msg string) {
				ev.SetMessage(msg)
			}))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool {
				return ev.IsCancelled()
			}))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) {
				ev.SetCancelled(c)
			}))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerChat: %v\n", p.name, err)
			}
		}))

	case "PlayerMove":
		evplayer.PlayerMoveEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerMoveEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value {
				return vm.ToValue(newPlayerWrapper(ev.GetPlayer()))
			}))
			jsEvent.Set("getFromX", vm.ToValue(func() float64 { return ev.GetFrom().X() }))
			jsEvent.Set("getFromY", vm.ToValue(func() float64 { return ev.GetFrom().Y() }))
			jsEvent.Set("getFromZ", vm.ToValue(func() float64 { return ev.GetFrom().Z() }))
			jsEvent.Set("getToX", vm.ToValue(func() float64 { return ev.GetTo().X() }))
			jsEvent.Set("getToY", vm.ToValue(func() float64 { return ev.GetTo().Y() }))
			jsEvent.Set("getToZ", vm.ToValue(func() float64 { return ev.GetTo().Z() }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool {
				return ev.IsCancelled()
			}))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) {
				ev.SetCancelled(c)
			}))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerMove: %v\n", p.name, err)
			}
		}))

	case "BlockBreak":
		evplayer.BlockBreakEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.BlockBreakEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value {
				return vm.ToValue(newPlayerWrapper(ev.GetPlayer()))
			}))
			jsEvent.Set("getBlockX", vm.ToValue(func() int { return ev.GetBlock().X() }))
			jsEvent.Set("getBlockY", vm.ToValue(func() int { return ev.GetBlock().Y() }))
			jsEvent.Set("getBlockZ", vm.ToValue(func() int { return ev.GetBlock().Z() }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool {
				return ev.IsCancelled()
			}))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) {
				ev.SetCancelled(c)
			}))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler BlockBreak: %v\n", p.name, err)
			}
		}))

	case "BlockPlace":
		evplayer.BlockPlaceEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.BlockPlaceEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value {
				return vm.ToValue(newPlayerWrapper(ev.GetPlayer()))
			}))
			jsEvent.Set("getBlockX", vm.ToValue(func() int { return ev.GetBlock().X() }))
			jsEvent.Set("getBlockY", vm.ToValue(func() int { return ev.GetBlock().Y() }))
			jsEvent.Set("getBlockZ", vm.ToValue(func() int { return ev.GetBlock().Z() }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool {
				return ev.IsCancelled()
			}))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) {
				ev.SetCancelled(c)
			}))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler BlockPlace: %v\n", p.name, err)
			}
		}))

	case "PlayerJump":
		evplayer.PlayerJumpEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerJumpEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerJump: %v\n", p.name, err)
			}
		}))

	case "PlayerDeath":
		evplayer.PlayerDeathEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerDeathEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("getKeepInventory", vm.ToValue(func() bool { return ev.GetKeepInventory() }))
			jsEvent.Set("setKeepInventory", vm.ToValue(func(k bool) { ev.SetKeepInventory(k) }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerDeath: %v\n", p.name, err)
			}
		}))

	case "PlayerRespawn":
		evplayer.PlayerRespawnEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerRespawnEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("getX", vm.ToValue(func() float64 { return ev.GetX() }))
			jsEvent.Set("getY", vm.ToValue(func() float64 { return ev.GetY() }))
			jsEvent.Set("getZ", vm.ToValue(func() float64 { return ev.GetZ() }))
			jsEvent.Set("setPosition", vm.ToValue(func(x, y, z float64) { ev.SetPosition(x, y, z) }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerRespawn: %v\n", p.name, err)
			}
		}))

	case "PlayerHurt":
		evplayer.PlayerHurtEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerHurtEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("getDamage", vm.ToValue(func() float64 { return ev.GetDamage() }))
			jsEvent.Set("setDamage", vm.ToValue(func(d float64) { ev.SetDamage(d) }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerHurt: %v\n", p.name, err)
			}
		}))

	case "PlayerHeal":
		evplayer.PlayerHealEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerHealEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("getHealth", vm.ToValue(func() float64 { return ev.GetHealth() }))
			jsEvent.Set("setHealth", vm.ToValue(func(h float64) { ev.SetHealth(h) }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerHeal: %v\n", p.name, err)
			}
		}))

	case "PlayerExperienceGain":
		evplayer.PlayerExperienceGainEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerExperienceGainEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("getAmount", vm.ToValue(func() int { return ev.GetAmount() }))
			jsEvent.Set("setAmount", vm.ToValue(func(a int) { ev.SetAmount(a) }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerExperienceGain: %v\n", p.name, err)
			}
		}))

	case "PlayerToggleSprint":
		evplayer.PlayerToggleSprintEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerToggleSprintEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("isSprinting", vm.ToValue(func() bool { return ev.IsSprinting() }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerToggleSprint: %v\n", p.name, err)
			}
		}))

	case "PlayerToggleSneak":
		evplayer.PlayerToggleSneakEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerToggleSneakEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("isSneaking", vm.ToValue(func() bool { return ev.IsSneaking() }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerToggleSneak: %v\n", p.name, err)
			}
		}))

	case "PlayerItemDrop":
		evplayer.PlayerItemDropEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerItemDropEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("getItemCount", vm.ToValue(func() int { return ev.GetItemCount() }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerItemDrop: %v\n", p.name, err)
			}
		}))

	case "PlayerItemPickup":
		evplayer.PlayerItemPickupEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerItemPickupEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("getItemCount", vm.ToValue(func() int { return ev.GetItemCount() }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerItemPickup: %v\n", p.name, err)
			}
		}))

	case "PlayerFoodLoss":
		evplayer.PlayerFoodLossEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerFoodLossEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("getFrom", vm.ToValue(func() int { return ev.GetFrom() }))
			jsEvent.Set("getTo", vm.ToValue(func() int { return ev.GetTo() }))
			jsEvent.Set("setTo", vm.ToValue(func(t int) { ev.SetTo(t) }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerFoodLoss: %v\n", p.name, err)
			}
		}))

	case "PlayerTeleport":
		evplayer.PlayerTeleportEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerTeleportEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("getX", vm.ToValue(func() float64 { return ev.GetX() }))
			jsEvent.Set("getY", vm.ToValue(func() float64 { return ev.GetY() }))
			jsEvent.Set("getZ", vm.ToValue(func() float64 { return ev.GetZ() }))
			jsEvent.Set("setPosition", vm.ToValue(func(x, y, z float64) { ev.SetPosition(x, y, z) }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerTeleport: %v\n", p.name, err)
			}
		}))

	case "PlayerAttackEntity":
		evplayer.PlayerAttackEntityEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerAttackEntityEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("getForce", vm.ToValue(func() float64 { return ev.GetForce() }))
			jsEvent.Set("isCritical", vm.ToValue(func() bool { return ev.IsCritical() }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerAttackEntity: %v\n", p.name, err)
			}
		}))

	case "PlayerItemUse":
		evplayer.PlayerItemUseEventHandlers.Register(makeListener(p, func(e event.Event) {
			ev := e.(*evplayer.PlayerItemUseEvent)
			jsEvent := vm.NewObject()
			jsEvent.Set("getPlayer", vm.ToValue(func() goja.Value { return vm.ToValue(newPlayerWrapper(ev.GetPlayer())) }))
			jsEvent.Set("isCancelled", vm.ToValue(func() bool { return ev.IsCancelled() }))
			jsEvent.Set("setCancelled", vm.ToValue(func(c bool) { ev.SetCancelled(c) }))
			if _, err := callback(goja.Undefined(), jsEvent); err != nil {
				fmt.Printf("[%s] Error en handler PlayerItemUse: %v\n", p.name, err)
			}
		}))

	default:
		fmt.Printf("[%s] Advertencia: evento desconocido '%s'\n", p.name, eventName)
	}
}

// GetPlugins retorna todos los plugins JS cargados.
func (l *Loader) GetPlugins() []*ScriptPlugin {
	return l.plugins
}
