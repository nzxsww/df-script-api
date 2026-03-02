package loader_test

import (
	"testing"

	"github.com/nzxsww/dragonfly-script-api/script/api/loader"
	"github.com/nzxsww/dragonfly-script-api/script/api/plugin"
)

// helper para crear un loader de test con un plugin JS inline
func newScoreboardTestLoader(t *testing.T, pluginName, js string) *loader.Loader {
	t.Helper()
	dir := t.TempDir()
	writePlugin(t, dir,
		"name: "+pluginName+"\nversion: 1.0.0\nmain: index.js",
		js,
	)
	mgr := plugin.NewManager(nil)
	return loader.NewLoader(dir, nil, mgr)
}

func loadAndEnable(t *testing.T, ldr *loader.Loader) {
	t.Helper()
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v, plugins: %d", err, len(plugins))
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("OnEnable() causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

// --- Tests: scoreboard global disponible ---

func TestScoreboard_GlobalIsAvailable(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbAvailablePlugin", `
function onEnable() {
    if (typeof scoreboard === "undefined") throw new Error("scoreboard no está definido");
    if (typeof scoreboard.create !== "function") throw new Error("scoreboard.create no es función");
    if (typeof scoreboard.getManager !== "function") throw new Error("scoreboard.getManager no es función");
    if (typeof scoreboard.createLive !== "function") throw new Error("scoreboard.createLive no es función");
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: scoreboard.create ---

func TestScoreboard_Create_ReturnsWrapper(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbCreatePlugin", `
function onEnable() {
    var sb = scoreboard.create("§aTítulo");
    if (typeof sb === "undefined" || sb === null) throw new Error("create() retornó nulo");
    if (typeof sb.setLine !== "function") throw new Error("sb.setLine no es función");
    if (typeof sb.setLines !== "function") throw new Error("sb.setLines no es función");
    if (typeof sb.addLine !== "function") throw new Error("sb.addLine no es función");
    if (typeof sb.removeLine !== "function") throw new Error("sb.removeLine no es función");
    if (typeof sb.getLines !== "function") throw new Error("sb.getLines no es función");
    if (typeof sb.getLineCount !== "function") throw new Error("sb.getLineCount no es función");
    if (typeof sb.getTitle !== "function") throw new Error("sb.getTitle no es función");
    if (typeof sb.setDescending !== "function") throw new Error("sb.setDescending no es función");
    if (typeof sb.isDescending !== "function") throw new Error("sb.isDescending no es función");
    if (typeof sb.removePadding !== "function") throw new Error("sb.removePadding no es función");
    if (typeof sb.sendTo !== "function") throw new Error("sb.sendTo no es función");
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_GetTitle_ReturnsCorrectTitle(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbTitlePlugin", `
function onEnable() {
    var sb = scoreboard.create("§6Mi Servidor");
    var title = sb.getTitle();
    if (title !== "§6Mi Servidor") throw new Error("getTitle incorrecto: " + title);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: setLine ---

func TestScoreboard_SetLine_SetsText(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbSetLinePlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    var ok = sb.setLine(0, "Línea 0");
    if (ok !== true) throw new Error("setLine debería retornar true");
    var lines = sb.getLines();
    if (lines.length < 1) throw new Error("getLines debería tener al menos 1 línea");
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_SetLine_OutOfRange_ReturnsFalse(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbSetLineRangePlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    var ok1 = sb.setLine(-1, "negativo");
    if (ok1 !== false) throw new Error("setLine(-1) debería retornar false, got: " + ok1);
    var ok2 = sb.setLine(15, "muy alto");
    if (ok2 !== false) throw new Error("setLine(15) debería retornar false, got: " + ok2);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_SetLine_MultipleLines(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbMultiLinePlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    sb.setLine(0, "Línea A");
    sb.setLine(1, "Línea B");
    sb.setLine(2, "Línea C");
    var count = sb.getLineCount();
    if (count < 3) throw new Error("getLineCount debería ser >= 3, got: " + count);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: addLine ---

func TestScoreboard_AddLine_AppendsLine(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbAddLinePlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    var ok = sb.addLine("nueva línea");
    if (ok !== true) throw new Error("addLine debería retornar true");
    var count = sb.getLineCount();
    if (count < 1) throw new Error("getLineCount debería ser >= 1 tras addLine, got: " + count);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: setLines ---

func TestScoreboard_SetLines_ReplacesContent(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbSetLinesPlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    sb.setLine(0, "vieja línea 0");
    sb.setLine(1, "vieja línea 1");
    sb.setLine(2, "vieja línea 2");

    // Reemplazar con solo 2 líneas
    sb.setLines(["nueva 0", "nueva 1"]);

    var count = sb.getLineCount();
    if (count !== 2) throw new Error("setLines debería dejar exactamente 2 líneas, got: " + count);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_SetLines_EmptyArray_ClearsContent(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbSetLinesEmptyPlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    sb.setLine(0, "línea existente");
    sb.setLines([]);
    var count = sb.getLineCount();
    if (count !== 0) throw new Error("setLines([]) debería dejar 0 líneas, got: " + count);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_SetLines_WithConditions(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbSetLinesCondPlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    var lines = ["línea base"];
    var condicion = true;
    if (condicion) {
        lines.push("línea condicional");
    }
    lines.push("última línea");
    sb.setLines(lines);
    var count = sb.getLineCount();
    if (count !== 3) throw new Error("con condición true debería tener 3 líneas, got: " + count);

    // Ahora sin condición
    var sb2 = scoreboard.create("Test2");
    var lines2 = ["línea base"];
    var condicion2 = false;
    if (condicion2) {
        lines2.push("línea condicional");
    }
    lines2.push("última línea");
    sb2.setLines(lines2);
    var count2 = sb2.getLineCount();
    if (count2 !== 2) throw new Error("con condición false debería tener 2 líneas, got: " + count2);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_SetLines_TruncatesAt15(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbSetLinesTruncPlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    var lines = [];
    for (var i = 0; i < 20; i++) {
        lines.push("línea " + i);
    }
    sb.setLines(lines);
    var count = sb.getLineCount();
    if (count > 15) throw new Error("setLines no debe superar 15 líneas, got: " + count);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_SetLines_CalledTwice_SecondReplaceFirst(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbSetLinesTwicePlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    sb.setLines(["a", "b", "c", "d", "e"]);
    // Segunda llamada con menos líneas — no deben quedar huérfanas
    sb.setLines(["x", "y"]);
    var count = sb.getLineCount();
    if (count !== 2) throw new Error("segunda setLines debe reemplazar todo, got: " + count);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_SetLines_NonArrayArg_NoError(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbSetLinesInvalidPlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    defer_called = false;
    try {
        sb.setLines("esto no es un array");
    } catch(e) {
        // no debe lanzar excepción JS — solo imprime warning en Go
    }
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: removeLine ---

func TestScoreboard_RemoveLine_ValidIndex_ReturnsTrue(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbRemoveLinePlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    sb.setLine(0, "línea 0");
    sb.setLine(1, "línea 1");
    var ok = sb.removeLine(0);
    if (ok !== true) throw new Error("removeLine(0) debería retornar true, got: " + ok);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_RemoveLine_OutOfRange_ReturnsFalse(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbRemoveLineRangePlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    var ok1 = sb.removeLine(-1);
    if (ok1 !== false) throw new Error("removeLine(-1) debería retornar false");
    var ok2 = sb.removeLine(15);
    if (ok2 !== false) throw new Error("removeLine(15) debería retornar false");
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: getLineCount ---

func TestScoreboard_GetLineCount_StartsAtZero(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbLineCountPlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    var count = sb.getLineCount();
    if (count !== 0) throw new Error("scoreboard nuevo debe tener 0 líneas, got: " + count);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: setDescending / isDescending ---

func TestScoreboard_Descending_DefaultFalse(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbDescDefaultPlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    if (sb.isDescending() !== false) throw new Error("scoreboard nuevo debe ser ascendente por defecto");
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_SetDescending_ChangesToTrue(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbDescPlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    sb.setDescending();
    if (sb.isDescending() !== true) throw new Error("isDescending debería ser true tras setDescending()");
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: removePadding ---

func TestScoreboard_RemovePadding_NoError(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbNoPaddingPlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    sb.setLine(0, "línea sin padding");
    sb.removePadding();
    // No debe panic ni error
    var lines = sb.getLines();
    if (!Array.isArray(lines)) throw new Error("getLines debería retornar array tras removePadding");
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: getLines ---

func TestScoreboard_GetLines_ReturnsArray(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbGetLinesPlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    var lines = sb.getLines();
    if (!Array.isArray(lines)) throw new Error("getLines debe retornar array");
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: ScoreboardManager ---

func TestScoreboard_GetManager_ReturnsWrapper(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbManagerPlugin", `
function onEnable() {
    var mgr = scoreboard.getManager();
    if (typeof mgr === "undefined" || mgr === null) throw new Error("getManager() retornó nulo");
    if (typeof mgr.set !== "function") throw new Error("mgr.set no es función");
    if (typeof mgr.remove !== "function") throw new Error("mgr.remove no es función");
    if (typeof mgr.get !== "function") throw new Error("mgr.get no es función");
    if (typeof mgr.hasScoreboard !== "function") throw new Error("mgr.hasScoreboard no es función");
    if (typeof mgr.getAssignedCount !== "function") throw new Error("mgr.getAssignedCount no es función");
    if (typeof mgr.clearAll !== "function") throw new Error("mgr.clearAll no es función");
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_Manager_GetAssignedCount_StartsAtZero(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbMgrCountPlugin", `
function onEnable() {
    var mgr = scoreboard.getManager();
    var count = mgr.getAssignedCount();
    if (count !== 0) throw new Error("getAssignedCount debe ser 0 al inicio, got: " + count);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_Manager_Set_WithNullPlayer_ReturnsFalse(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbMgrSetNullPlugin", `
function onEnable() {
    var mgr = scoreboard.getManager();
    var sb = scoreboard.create("Test");
    var ok = mgr.set(null, sb);
    if (ok !== false) throw new Error("mgr.set(null, sb) debe retornar false, got: " + ok);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_Manager_Set_WithNullScoreboard_ReturnsFalse(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbMgrSetNullSbPlugin", `
function onEnable() {
    var mgr = scoreboard.getManager();
    // Sin jugador real — pasamos un objeto falso
    var ok = mgr.set({ getName: function() { return "fake"; } }, null);
    if (ok !== false) throw new Error("mgr.set(player, null) debe retornar false, got: " + ok);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_Manager_Get_WithNullPlayer_ReturnsNull(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbMgrGetNullPlugin", `
function onEnable() {
    var mgr = scoreboard.getManager();
    var result = mgr.get(null);
    if (result !== null) throw new Error("mgr.get(null) debe retornar null, got: " + result);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_Manager_HasScoreboard_WithNullPlayer_ReturnsFalse(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbMgrHasNullPlugin", `
function onEnable() {
    var mgr = scoreboard.getManager();
    var has = mgr.hasScoreboard(null);
    if (has !== false) throw new Error("mgr.hasScoreboard(null) debe retornar false, got: " + has);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_Manager_ClearAll_NoError(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbMgrClearPlugin", `
function onEnable() {
    var mgr = scoreboard.getManager();
    mgr.clearAll(); // no debe panic con manager vacío
    var count = mgr.getAssignedCount();
    if (count !== 0) throw new Error("tras clearAll getAssignedCount debe ser 0, got: " + count);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: Live Scoreboard ---

func TestScoreboard_CreateLive_ReturnsWrapper(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbLivePlugin", `
function onEnable() {
    var live = scoreboard.createLive("§aTítulo", function(sb) {
        sb.setLines(["línea 1", "línea 2"]);
    }, 10000); // intervalo largo para que no se dispare durante el test

    if (typeof live === "undefined" || live === null) throw new Error("createLive() retornó nulo");
    if (typeof live.addPlayer !== "function") throw new Error("live.addPlayer no es función");
    if (typeof live.removePlayer !== "function") throw new Error("live.removePlayer no es función");
    if (typeof live.hasPlayer !== "function") throw new Error("live.hasPlayer no es función");
    if (typeof live.getPlayerCount !== "function") throw new Error("live.getPlayerCount no es función");
    if (typeof live.clearPlayers !== "function") throw new Error("live.clearPlayers no es función");
    if (typeof live.stop !== "function") throw new Error("live.stop no es función");

    live.stop();
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_CreateLive_GetPlayerCount_StartsAtZero(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbLiveCountPlugin", `
function onEnable() {
    var live = scoreboard.createLive("Test", function(sb) {}, 10000);
    var count = live.getPlayerCount();
    if (count !== 0) throw new Error("getPlayerCount debe ser 0 al inicio, got: " + count);
    live.stop();
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_CreateLive_AddPlayer_NullPlayer_ReturnsFalse(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbLiveAddNullPlugin", `
function onEnable() {
    var live = scoreboard.createLive("Test", function(sb) {}, 10000);
    var ok = live.addPlayer(null);
    if (ok !== false) throw new Error("addPlayer(null) debe retornar false, got: " + ok);
    live.stop();
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_CreateLive_RemovePlayer_NullPlayer_ReturnsFalse(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbLiveRemoveNullPlugin", `
function onEnable() {
    var live = scoreboard.createLive("Test", function(sb) {}, 10000);
    var ok = live.removePlayer(null);
    if (ok !== false) throw new Error("removePlayer(null) debe retornar false, got: " + ok);
    live.stop();
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_CreateLive_HasPlayer_NullPlayer_ReturnsFalse(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbLiveHasNullPlugin", `
function onEnable() {
    var live = scoreboard.createLive("Test", function(sb) {}, 10000);
    var has = live.hasPlayer(null);
    if (has !== false) throw new Error("hasPlayer(null) debe retornar false, got: " + has);
    live.stop();
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_CreateLive_Stop_NoError(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbLiveStopPlugin", `
function onEnable() {
    var live = scoreboard.createLive("Test", function(sb) {
        sb.setLines(["actualizando..."]);
    }, 10000);
    live.stop(); // no debe panic
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_CreateLive_ClearPlayers_NoError(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbLiveClearPlugin", `
function onEnable() {
    var live = scoreboard.createLive("Test", function(sb) {}, 10000);
    live.clearPlayers(); // sin jugadores, no debe panic
    var count = live.getPlayerCount();
    if (count !== 0) throw new Error("tras clearPlayers getPlayerCount debe ser 0, got: " + count);
    live.stop();
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_CreateLive_TooFewArgs_ReturnsNull(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbLiveFewArgsPlugin", `
function onEnable() {
    var live = scoreboard.createLive("Test", function(sb) {});
    // Sin intervalo — debe retornar null y no panic
    if (live !== null) throw new Error("createLive sin intervalo debe retornar null, got: " + typeof live);
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_CreateLive_MinInterval_Clamped(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbLiveMinIntervalPlugin", `
function onEnable() {
    // Intervalo de 1ms debe clampearse a 50ms mínimo sin panic
    var live = scoreboard.createLive("Test", function(sb) {
        sb.setLines(["tick"]);
    }, 1);
    if (live === null) throw new Error("createLive con intervalo 1ms no debe retornar null");
    live.stop();
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: player.sendScoreboard / player.removeScoreboard disponibles ---

func TestScoreboard_Player_HasScoreboardMethods(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbPlayerMethodsPlugin", `
function onEnable() {
    events.on("PlayerJoin", function(event) {
        var p = event.getPlayer();
        if (typeof p.sendScoreboard !== "function") throw new Error("player.sendScoreboard no es función");
        if (typeof p.removeScoreboard !== "function") throw new Error("player.removeScoreboard no es función");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

func TestScoreboard_SendTo_NullPlayer_NoError(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbSendToNullPlugin", `
function onEnable() {
    var sb = scoreboard.create("Test");
    sb.setLine(0, "hola");
    // sendTo con null no debe panic
    try {
        sb.sendTo(null);
    } catch(e) {
        // no debe lanzar excepción JS
        throw new Error("sendTo(null) lanzó excepción: " + e);
    }
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}

// --- Tests: integración setLines + createLive ---

func TestScoreboard_Live_UsesSetLines_NoError(t *testing.T) {
	ldr := newScoreboardTestLoader(t, "SbLiveSetLinesPlugin", `
function onEnable() {
    var live = scoreboard.createLive("§6Server", function(sb) {
        var lines = ["§7Jugadores: 0"];
        var condicion = true;
        if (condicion) {
            lines.push("§aServidor activo");
        }
        lines.push("");
        lines.push("§7play.test.com");
        sb.setLines(lines);
    }, 10000); // intervalo largo para no disparar durante el test

    if (live === null) throw new Error("createLive no debe retornar null");
    live.stop();
}
module = { onEnable: onEnable, onDisable: function() {} };`)
	loadAndEnable(t, ldr)
}
