package loader

import (
	"fmt"
	"sync"
	"time"

	dfplayer "github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/dop251/goja"
)

// scoreboardManager gestiona qué scoreboard está asignado a cada jugador (por UUID).
// Es seguro para uso concurrente. Cada ScriptPlugin tiene su propio manager.
type scoreboardManager struct {
	// assigned: UUID → *scoreboard.Scoreboard actualmente visible para ese jugador.
	assigned sync.Map
}

// newScoreboardManager crea un manager vacío.
func newScoreboardManager() *scoreboardManager {
	return &scoreboardManager{}
}

// set asigna un scoreboard a un jugador y lo envía.
func (m *scoreboardManager) set(p *dfplayer.Player, board *scoreboard.Scoreboard) {
	m.assigned.Store(p.UUID().String(), board)
	p.SendScoreboard(board)
}

// remove quita el scoreboard asignado a un jugador y lo elimina de la pantalla.
func (m *scoreboardManager) remove(p *dfplayer.Player) {
	m.assigned.Delete(p.UUID().String())
	p.RemoveScoreboard()
}

// get retorna el scoreboard asignado a un jugador, o nil si no tiene ninguno.
func (m *scoreboardManager) get(uuid string) *scoreboard.Scoreboard {
	if v, ok := m.assigned.Load(uuid); ok {
		return v.(*scoreboard.Scoreboard)
	}
	return nil
}

// getAssignedPlayers retorna todos los jugadores con scoreboard asignado como UUIDs.
func (m *scoreboardManager) getAssignedUUIDs() []string {
	var uuids []string
	m.assigned.Range(func(k, _ interface{}) bool {
		uuids = append(uuids, k.(string))
		return true
	})
	return uuids
}

// liveScoreboard es un scoreboard que se actualiza automáticamente en intervalos.
// Mantiene una lista de jugadores a los que debe reenviar el scoreboard en cada tick.
// El ticker corre en su propia goroutine y todas las llamadas al VM se sincronizan
// con el mutex del plugin (Goja no es thread-safe).
type liveScoreboard struct {
	board    *scoreboard.Scoreboard
	title    string
	players  sync.Map // UUID → *dfplayer.Player
	ticker   *time.Ticker
	stopCh   chan struct{}
	vm       *goja.Runtime
	plugin   *ScriptPlugin
	pluginName string
}

// registerScoreboardAPI expone el objeto global `scoreboard` al JavaScript.
//
// Uso en JS:
//
//	// Scoreboard simple:
//	var sb = scoreboard.create("§aMi Servidor");
//	sb.setLine(0, "§7Jugadores: 5");
//	player.sendScoreboard(sb);
//
//	// ScoreboardManager (gestión por jugador):
//	var mgr = scoreboard.getManager();
//	mgr.set(player, sb);           // asigna y envía al jugador
//	mgr.remove(player);            // quita el scoreboard del jugador
//	mgr.get(player);               // retorna el scoreboard activo del jugador o null
//	mgr.hasScoreboard(player);     // true si el jugador tiene scoreboard asignado
//
//	// Live scoreboard (auto-actualizable):
//	var live = scoreboard.createLive("§aTítulo", function(sb) {
//	    sb.setLine(0, "§7Jugadores: " + server.getPlayerCount());
//	}, 1000); // intervalo en ms
//	live.addPlayer(player);        // agrega jugador al live (lo verá actualizado)
//	live.removePlayer(player);     // saca al jugador del live
//	live.stop();                   // detiene el auto-update
func (l *Loader) registerScoreboardAPI(vm *goja.Runtime, p *ScriptPlugin) {
	mgr := newScoreboardManager()

	vm.Set("scoreboard", map[string]interface{}{
		// create(title) — crea un nuevo scoreboard estático.
		"create": func(title string) map[string]interface{} {
			board := scoreboard.New(title)
			return newScoreboardWrapper(board)
		},

		// getManager() — retorna el ScoreboardManager del plugin.
		// Permite asignar/desasignar scoreboards por jugador y consultar el estado.
		"getManager": func() map[string]interface{} {
			return newScoreboardManagerWrapper(mgr, vm, p)
		},

		// createLive(title, fn, intervalMs) — crea un scoreboard que se auto-actualiza.
		// fn recibe (sb, playerData) y se llama una vez por jugador registrado.
		// Goja no es thread-safe, así que el callback se ejecuta bajo vm lock.
		// NO usar server.getPlayers() dentro del callback — usar el playerData recibido.
		"createLive": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) < 3 {
				fmt.Printf("[%s] scoreboard.createLive requiere 3 argumentos: título, función, intervaloMs\n", p.name)
				return goja.Null()
			}
			title := call.Argument(0).String()
			intervalMs := call.Argument(2).ToInteger()
			if intervalMs < 50 {
				intervalMs = 50
			}

			live := &liveScoreboard{
				title:      title,
				ticker:     time.NewTicker(time.Duration(intervalMs) * time.Millisecond),
				stopCh:     make(chan struct{}),
				vm:         vm,
				plugin:     p,
				pluginName: p.name,
			}

			// tickCh recibe los datos del jugador desde la goroutine del ticker
			// y los procesa en el event loop del VM (mismo goroutine que creó el VM).
			// La goroutine del ticker solo recoge los jugadores y notifica — no toca el VM.
			// El VM llama al callback JS en su propia goroutine a través del ticker channel.
			// Como setInterval, corremos el callback en la goroutine del ticker (mismo patrón).
			updateFn, fnOk := goja.AssertFunction(call.Argument(1))

			go func() {
				for {
					select {
					case <-live.stopCh:
						live.ticker.Stop()
						return
					case <-live.ticker.C:
						// Recolectar jugadores actuales en Go puro
						type playerData struct {
							pl   *dfplayer.Player
							name string
							hp   float64
							maxHp float64
							food  int
							x, y, z float64
							gm   string
							flying bool
							expLvl int
						}
						var players []playerData
						live.players.Range(func(_, v interface{}) bool {
							pl := v.(*dfplayer.Player)
							pos := pl.Position()
							players = append(players, playerData{
								pl:     pl,
								name:   pl.Name(),
								hp:     pl.Health(),
								maxHp:  pl.MaxHealth(),
								food:   pl.Food(),
								x:      pos[0],
								y:      pos[1],
								z:      pos[2],
								gm:     gmName(pl),
								flying: pl.Flying(),
								expLvl: pl.ExperienceLevel(),
							})
							return true
						})

						// Si no hay jugadores, no llamar al VM
						if len(players) == 0 {
							continue
						}

						// Llamar al callback JS solo si se proporcionó una función
						if !fnOk {
							// Sin callback: enviar scoreboard vacío con título
							for _, pd := range players {
								board := scoreboard.New(live.title)
								pd.pl.SendScoreboard(board)
							}
							continue
						}

						// Llamar al callback JS con un playerData wrapper liviano
						// El callback se ejecuta bajo vm lock para evitar race conditions de Goja.
						for _, pd := range players {
							freshBoard := scoreboard.New(live.title)
							boardPtr := &freshBoard

							// Construir objeto JS liviano con solo los datos ya leídos en Go
							playerDataMap := map[string]interface{}{
								"getName":            func() string { return pd.name },
								"getHealth":          func() float64 { return pd.hp },
								"getMaxHealth":       func() float64 { return pd.maxHp },
								"getFoodLevel":       func() int { return pd.food },
								"getX":               func() float64 { return pd.x },
								"getY":               func() float64 { return pd.y },
								"getZ":               func() float64 { return pd.z },
								"getGameMode":        func() string { return pd.gm },
								"isFlying":           func() bool { return pd.flying },
								"getExperienceLevel": func() int { return pd.expLvl },
							}

							var callErr error
							live.plugin.withVMLock(func() {
								sbVal := vm.ToValue(newScoreboardWrapperFromPtr(boardPtr))
								plVal := vm.ToValue(playerDataMap)
								_, callErr = updateFn(goja.Undefined(), sbVal, plVal)
							})
							if callErr != nil {
								fmt.Printf("[%s] scoreboard.createLive: error en callback: %v\n", live.pluginName, callErr)
								continue
							}

							pd.pl.SendScoreboard(*boardPtr)
						}
					}
				}
			}()

			return vm.ToValue(newLiveScoreboardWrapper(live, vm, p))
		},
	})
}

// newScoreboardManagerWrapper construye el objeto JS que expone el ScoreboardManager.
func newScoreboardManagerWrapper(mgr *scoreboardManager, vm *goja.Runtime, p *ScriptPlugin) map[string]interface{} {
	return map[string]interface{}{
		// set(player, sb) — asigna un scoreboard al jugador y lo envía inmediatamente.
		"set": func(playerObj goja.Value, sbVal goja.Value) bool {
			pl := extractPlayerFromJS(playerObj)
			if pl == nil {
				fmt.Printf("[%s] scoreboardManager.set: jugador inválido\n", p.name)
				return false
			}
			board := extractBoardFromJS(sbVal)
			if board == nil {
				fmt.Printf("[%s] scoreboardManager.set: scoreboard inválido (¿usaste scoreboard.create()?)\n", p.name)
				return false
			}
			mgr.set(pl, board)
			return true
		},

		// remove(player) — quita el scoreboard asignado al jugador.
		"remove": func(playerObj goja.Value) {
			pl := extractPlayerFromJS(playerObj)
			if pl == nil {
				fmt.Printf("[%s] scoreboardManager.remove: jugador inválido\n", p.name)
				return
			}
			mgr.remove(pl)
		},

		// get(player) — retorna el scoreboardWrapper activo del jugador, o null si no tiene.
		"get": func(playerObj goja.Value) goja.Value {
			pl := extractPlayerFromJS(playerObj)
			if pl == nil {
				return goja.Null()
			}
			board := mgr.get(pl.UUID().String())
			if board == nil {
				return goja.Null()
			}
			return vm.ToValue(newScoreboardWrapper(board))
		},

		// hasScoreboard(player) — retorna true si el jugador tiene un scoreboard asignado.
		"hasScoreboard": func(playerObj goja.Value) bool {
			pl := extractPlayerFromJS(playerObj)
			if pl == nil {
				return false
			}
			return mgr.get(pl.UUID().String()) != nil
		},

		// getAssignedCount() — retorna la cantidad de jugadores con scoreboard asignado.
		"getAssignedCount": func() int {
			return len(mgr.getAssignedUUIDs())
		},

		// clearAll() — quita el scoreboard de TODOS los jugadores gestionados.
		// Útil en onDisable() para limpiar scoreboards al cerrar el servidor.
		"clearAll": func() {
			mgr.assigned.Range(func(k, v interface{}) bool {
				// Solo eliminar del mapa — no podemos llamar RemoveScoreboard sin el jugador
				mgr.assigned.Delete(k)
				return true
			})
		},
	}
}

// newLiveScoreboardWrapper construye el objeto JS para un liveScoreboard.
func newLiveScoreboardWrapper(live *liveScoreboard, vm *goja.Runtime, p *ScriptPlugin) map[string]interface{} {
	return map[string]interface{}{
		// addPlayer(player) — agrega un jugador al live scoreboard.
		// El jugador comenzará a recibir actualizaciones automáticas.
		"addPlayer": func(playerObj goja.Value) bool {
			pl := extractPlayerFromJS(playerObj)
			if pl == nil {
				fmt.Printf("[%s] live.addPlayer: jugador inválido\n", p.name)
				return false
			}
			live.players.Store(pl.UUID().String(), pl)
			return true
		},

		// removePlayer(player) — saca al jugador del live scoreboard y le quita el scoreboard.
		"removePlayer": func(playerObj goja.Value) bool {
			pl := extractPlayerFromJS(playerObj)
			if pl == nil {
				fmt.Printf("[%s] live.removePlayer: jugador inválido\n", p.name)
				return false
			}
			live.players.Delete(pl.UUID().String())
			pl.RemoveScoreboard()
			return true
		},

		// hasPlayer(player) — retorna true si el jugador está en el live scoreboard.
		"hasPlayer": func(playerObj goja.Value) bool {
			pl := extractPlayerFromJS(playerObj)
			if pl == nil {
				return false
			}
			_, ok := live.players.Load(pl.UUID().String())
			return ok
		},

		// getPlayerCount() — retorna la cantidad de jugadores en el live scoreboard.
		"getPlayerCount": func() int {
			count := 0
			live.players.Range(func(_, _ interface{}) bool {
				count++
				return true
			})
			return count
		},

		// stop() — detiene el auto-update y quita el scoreboard de todos los jugadores.
		"stop": func() {
			select {
			case live.stopCh <- struct{}{}:
			default:
			}
			live.players.Range(func(_, v interface{}) bool {
				pl := v.(*dfplayer.Player)
				pl.RemoveScoreboard()
				live.players.Delete(pl.UUID().String())
				return true
			})
		},

		// clearPlayers() — saca a todos los jugadores del live sin detenerlo.
		"clearPlayers": func() {
			live.players.Range(func(k, v interface{}) bool {
				pl := v.(*dfplayer.Player)
				pl.RemoveScoreboard()
				live.players.Delete(k)
				return true
			})
		},
	}
}

// extractPlayerFromJS extrae el *dfplayer.Player de un playerWrapper JS.
// El playerWrapper expone "_player" como campo interno para este propósito.
// Retorna nil si el objeto no es un playerWrapper válido.
func extractPlayerFromJS(playerObj goja.Value) *dfplayer.Player {
	if playerObj == nil || goja.IsNull(playerObj) || goja.IsUndefined(playerObj) {
		return nil
	}
	obj, ok := playerObj.(*goja.Object)
	if !ok {
		return nil
	}
	raw := obj.Get("_player")
	if raw == nil {
		return nil
	}
	pl, ok := raw.Export().(*dfplayer.Player)
	if !ok {
		return nil
	}
	return pl
}

// extractBoardFromJS extrae el *scoreboard.Scoreboard actual de un scoreboardWrapper JS.
// Lee "_boardPtr" (**scoreboard.Scoreboard) para obtener siempre el board vigente,
// incluso si setLines() lo reemplazó internamente.
func extractBoardFromJS(sbVal goja.Value) *scoreboard.Scoreboard {
	if sbVal == nil || goja.IsNull(sbVal) || goja.IsUndefined(sbVal) {
		return nil
	}
	obj, ok := sbVal.(*goja.Object)
	if !ok {
		return nil
	}
	raw := obj.Get("_boardPtr")
	if raw == nil {
		return nil
	}
	boardPtr, ok := raw.Export().(**scoreboard.Scoreboard)
	if !ok {
		return nil
	}
	return *boardPtr
}

// newScoreboardWrapper construye el objeto JS que envuelve un *scoreboard.Scoreboard.
// Usa un puntero al puntero (**scoreboard.Scoreboard) para que setLines() pueda
// reemplazar el board completo sin dejar líneas huérfanas.
func newScoreboardWrapper(board *scoreboard.Scoreboard) map[string]interface{} {
	// boardPtr permite que setLines() reemplace el board interno apuntando a uno nuevo.
	boardPtr := &board

	return map[string]interface{}{
		// getTitle() — retorna el título del scoreboard.
		"getTitle": func() string {
			return (*boardPtr).Name()
		},

		// setLine(index, text) — setea el texto de la línea en el índice dado (0-14).
		"setLine": func(index int, text string) bool {
			if index < 0 || index >= 15 {
				fmt.Printf("[scoreboardWrapper] setLine: índice fuera de rango %d (debe ser 0-14)\n", index)
				return false
			}
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[scoreboardWrapper] setLine: error interno: %v\n", r)
				}
			}()
			(*boardPtr).Set(index, text)
			return true
		},

		// removeLine(index) — remueve la línea en el índice dado.
		"removeLine": func(index int) bool {
			if index < 0 || index >= 15 {
				fmt.Printf("[scoreboardWrapper] removeLine: índice fuera de rango %d (debe ser 0-14)\n", index)
				return false
			}
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[scoreboardWrapper] removeLine: error interno: %v\n", r)
				}
			}()
			(*boardPtr).Remove(index)
			return true
		},

		// addLine(text) — agrega una línea al final del scoreboard.
		// Retorna false si ya se alcanzaron las 15 líneas máximas.
		"addLine": func(text string) bool {
			_, err := (*boardPtr).WriteString(text)
			if err != nil {
				fmt.Printf("[scoreboardWrapper] addLine: %v\n", err)
				return false
			}
			return true
		},

		// setLines(array) — reemplaza TODO el contenido del scoreboard con el array dado.
		// Recrea el board internamente para garantizar que no queden líneas huérfanas.
		// Máximo 15 líneas — las que superen el límite se ignoran con un aviso.
		// Ejemplo:
		//   var lines = ["§7Jugadores: 5", "§7Mapa: Lobby"];
		//   if (condicion) lines.push("§aExtra");
		//   sb.setLines(lines);
		"setLines": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return goja.Undefined()
			}
			exported := call.Argument(0).Export()
			var lines []string
			switch v := exported.(type) {
			case []interface{}:
				for _, item := range v {
					lines = append(lines, fmt.Sprint(item))
				}
			case []string:
				lines = v
			default:
				fmt.Printf("[scoreboardWrapper] setLines: el argumento debe ser un array\n")
				return goja.Undefined()
			}
			if len(lines) > 15 {
				fmt.Printf("[scoreboardWrapper] setLines: máximo 15 líneas, se truncará\n")
				lines = lines[:15]
			}
			// Recrear el board con el mismo título para limpiar las líneas anteriores
			title := (*boardPtr).Name()
			fresh := scoreboard.New(title)
			for i, line := range lines {
				fresh.Set(i, line)
			}
			*boardPtr = fresh
			return goja.Undefined()
		},

		// getLines() — retorna todas las líneas actuales como array de strings.
		"getLines": func() []string {
			return (*boardPtr).Lines()
		},

		// getLineCount() — retorna la cantidad de líneas actuales.
		"getLineCount": func() int {
			return len((*boardPtr).Lines())
		},

		// setDescending() — invierte el orden de las líneas al mostrarse.
		"setDescending": func() {
			(*boardPtr).SetDescending()
		},

		// isDescending() — retorna si el scoreboard está en orden descendente.
		"isDescending": func() bool {
			return (*boardPtr).Descending()
		},

		// removePadding() — elimina el padding automático de espacios en cada línea.
		"removePadding": func() {
			(*boardPtr).RemovePadding()
		},

		// sendTo(player) — método de conveniencia: envía este scoreboard al jugador dado.
		"sendTo": func(playerObj goja.Value) {
			if playerObj == nil || goja.IsNull(playerObj) || goja.IsUndefined(playerObj) {
				fmt.Printf("[scoreboardWrapper] sendTo: jugador nulo\n")
				return
			}
			if obj, ok := playerObj.(*goja.Object); ok {
				if raw := obj.Get("_player"); raw != nil {
					if pl, ok := raw.Export().(*dfplayer.Player); ok {
						pl.SendScoreboard(*boardPtr)
					}
				}
			}
		},

		// _boardPtr — referencia interna al **scoreboard.Scoreboard para uso desde Go.
		// Usamos el puntero al puntero para que extractBoardFromJS siempre lea el board actual.
		"_boardPtr": boardPtr,
	}
}

// gmName retorna el nombre del gamemode del jugador como string JS.
// Llama pl.GameMode() desde Go puro sin necesidad de Tx — seguro desde goroutines.
func gmName(pl *dfplayer.Player) string {
	switch pl.GameMode() {
	case world.GameModeSurvival:
		return "survival"
	case world.GameModeCreative:
		return "creative"
	case world.GameModeAdventure:
		return "adventure"
	case world.GameModeSpectator:
		return "spectator"
	default:
		return "survival"
	}
}

// newScoreboardWrapperFromPtr construye el objeto JS que envuelve un **scoreboard.Scoreboard.
// Igual a newScoreboardWrapper pero acepta un puntero ya existente — usado en createLive
// para que setLines() pueda actualizar el board original y detectar cambios.
func newScoreboardWrapperFromPtr(boardPtr **scoreboard.Scoreboard) map[string]interface{} {
	return map[string]interface{}{
		"getTitle": func() string { return (*boardPtr).Name() },
		"setLine": func(index int, text string) bool {
			if index < 0 || index >= 15 {
				return false
			}
			defer func() { recover() }()
			(*boardPtr).Set(index, text)
			return true
		},
		"removeLine": func(index int) bool {
			if index < 0 || index >= 15 {
				return false
			}
			defer func() { recover() }()
			(*boardPtr).Remove(index)
			return true
		},
		"addLine": func(text string) bool {
			_, err := (*boardPtr).WriteString(text)
			return err == nil
		},
		"setLines": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return goja.Undefined()
			}
			exported := call.Argument(0).Export()
			var lines []string
			switch v := exported.(type) {
			case []interface{}:
				for _, item := range v {
					lines = append(lines, fmt.Sprint(item))
				}
			case []string:
				lines = v
			default:
				return goja.Undefined()
			}
			if len(lines) > 15 {
				lines = lines[:15]
			}
			title := (*boardPtr).Name()
			fresh := scoreboard.New(title)
			for i, line := range lines {
				fresh.Set(i, line)
			}
			*boardPtr = fresh
			return goja.Undefined()
		},
		"getLines":      func() []string { return (*boardPtr).Lines() },
		"getLineCount":  func() int { return len((*boardPtr).Lines()) },
		"setDescending": func() { (*boardPtr).SetDescending() },
		"isDescending":  func() bool { return (*boardPtr).Descending() },
		"removePadding": func() { (*boardPtr).RemovePadding() },
		"sendTo": func(playerObj goja.Value) {
			pl := extractPlayerFromJS(playerObj)
			if pl != nil {
				pl.SendScoreboard(*boardPtr)
			}
		},
		"_boardPtr": boardPtr,
	}
}

// sendScoreboardToPlayer envía un scoreboardWrapper JS al jugador dado.
// Extrae el *scoreboard.Scoreboard del objeto JS usando extractBoardFromJS.
func sendScoreboardToPlayer(p *dfplayer.Player, sbVal goja.Value) {
	board := extractBoardFromJS(sbVal)
	if board == nil {
		fmt.Printf("[scoreboard] sendScoreboard: scoreboard inválido (¿usaste scoreboard.create()?)\n")
		return
	}
	p.SendScoreboard(board)
}
