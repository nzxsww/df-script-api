package loader

import (
	"fmt"
	"sync"
	"time"

	dfplayer "github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
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
type liveScoreboard struct {
	board    *scoreboard.Scoreboard
	title    string
	players  sync.Map // UUID → *dfplayer.Player
	ticker   *time.Ticker
	stopCh   chan struct{}
	updateFn goja.Callable
	vm       *goja.Runtime
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
		// fn(sb) es llamado cada intervalMs ms. El scoreboard se recrea limpio en cada tick
		// y se reenvía a todos los jugadores asignados al live scoreboard.
		// Retorna un liveScoreboardWrapper con addPlayer/removePlayer/stop.
		"createLive": func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) < 3 {
				fmt.Printf("[%s] scoreboard.createLive requiere 3 argumentos: título, función, intervaloMs\n", p.name)
				return goja.Null()
			}
			title := call.Argument(0).String()
			updateFn, ok := goja.AssertFunction(call.Argument(1))
			if !ok {
				fmt.Printf("[%s] scoreboard.createLive: el segundo argumento debe ser una función\n", p.name)
				return goja.Null()
			}
			intervalMs := call.Argument(2).ToInteger()
			if intervalMs < 50 {
				intervalMs = 50 // mínimo 50ms para no saturar
			}

			live := &liveScoreboard{
				title:      title,
				ticker:     time.NewTicker(time.Duration(intervalMs) * time.Millisecond),
				stopCh:     make(chan struct{}),
				updateFn:   updateFn,
				vm:         vm,
				pluginName: p.name,
			}

			// Goroutine que hace el tick y reenvía a los jugadores asignados
			go func() {
				for {
					select {
					case <-live.stopCh:
						live.ticker.Stop()
						return
					case <-live.ticker.C:
						// Crear un scoreboard fresco en cada tick
						freshBoard := scoreboard.New(live.title)
						sbWrapper := vm.ToValue(newScoreboardWrapper(freshBoard))

						// Llamar al callback JS con el scoreboard fresco
						if _, err := live.updateFn(goja.Undefined(), sbWrapper); err != nil {
							fmt.Printf("[%s] scoreboard.createLive: error en callback: %v\n", live.pluginName, err)
							continue
						}

						// Reenviar a todos los jugadores asignados
						live.players.Range(func(_, v interface{}) bool {
							pl := v.(*dfplayer.Player)
							pl.SendScoreboard(freshBoard)
							return true
						})
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

// extractBoardFromJS extrae el *scoreboard.Scoreboard de un scoreboardWrapper JS.
func extractBoardFromJS(sbVal goja.Value) *scoreboard.Scoreboard {
	if sbVal == nil || goja.IsNull(sbVal) || goja.IsUndefined(sbVal) {
		return nil
	}
	obj, ok := sbVal.(*goja.Object)
	if !ok {
		return nil
	}
	raw := obj.Get("_board")
	if raw == nil {
		return nil
	}
	board, ok := raw.Export().(*scoreboard.Scoreboard)
	if !ok {
		return nil
	}
	return board
}

// newScoreboardWrapper construye el objeto JS que envuelve un *scoreboard.Scoreboard.
// El scoreboard de Dragonfly es inmutable después de enviarse — siempre reenviar con sendScoreboard.
func newScoreboardWrapper(board *scoreboard.Scoreboard) map[string]interface{} {
	return map[string]interface{}{
		// getTitle() — retorna el título del scoreboard.
		"getTitle": func() string {
			return board.Name()
		},

		// setLine(index, text) — setea el texto de la línea en el índice dado (0-14).
		// Si el índice supera las líneas actuales, rellena con líneas vacías.
		// Panics internamente si index < 0 o index >= 15 — lo capturamos con recover.
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
			board.Set(index, text)
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
			board.Remove(index)
			return true
		},

		// addLine(text) — agrega una línea al final del scoreboard.
		// Retorna false si ya se alcanzaron las 15 líneas máximas.
		"addLine": func(text string) bool {
			_, err := board.WriteString(text)
			if err != nil {
				fmt.Printf("[scoreboardWrapper] addLine: %v\n", err)
				return false
			}
			return true
		},

		// getLines() — retorna todas las líneas actuales como array de strings.
		"getLines": func() []string {
			return board.Lines()
		},

		// getLineCount() — retorna la cantidad de líneas actuales.
		"getLineCount": func() int {
			return len(board.Lines())
		},

		// setDescending() — invierte el orden de las líneas al mostrarse.
		"setDescending": func() {
			board.SetDescending()
		},

		// isDescending() — retorna si el scoreboard está en orden descendente.
		"isDescending": func() bool {
			return board.Descending()
		},

		// removePadding() — elimina el padding automático de espacios en cada línea.
		"removePadding": func() {
			board.RemovePadding()
		},

		// sendTo(player) — método de conveniencia: envía este scoreboard al jugador dado.
		// Equivalente a player.sendScoreboard(sb).
		"sendTo": func(playerObj goja.Value) {
			if playerObj == nil || goja.IsNull(playerObj) || goja.IsUndefined(playerObj) {
				fmt.Printf("[scoreboardWrapper] sendTo: jugador nulo\n")
				return
			}
			if obj, ok := playerObj.(*goja.Object); ok {
				if raw := obj.Get("_player"); raw != nil {
					if pl, ok := raw.Export().(*dfplayer.Player); ok {
						pl.SendScoreboard(board)
					}
				}
			}
		},

		// _board — referencia interna al *scoreboard.Scoreboard para uso desde Go.
		"_board": board,
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
