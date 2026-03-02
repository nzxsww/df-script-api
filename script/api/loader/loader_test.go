package loader_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nzxsww/dragonfly-script-api/script/api/loader"
	"github.com/nzxsww/dragonfly-script-api/script/api/plugin"
)

// time se usa en TestLoader_SetTimeout_Executes
var _ = time.Millisecond

// helpers para crear plugins temporales en disco

func writePlugin(t *testing.T, dir, pluginYml, indexJS string) string {
	t.Helper()
	pluginDir := filepath.Join(dir, "test-plugin")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatalf("no se pudo crear directorio de plugin: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pluginDir, "plugin.yml"), []byte(pluginYml), 0644); err != nil {
		t.Fatalf("no se pudo escribir plugin.yml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pluginDir, "index.js"), []byte(indexJS), 0644); err != nil {
		t.Fatalf("no se pudo escribir index.js: %v", err)
	}
	return dir
}

func newTestLoader(t *testing.T, pluginDir string) *loader.Loader {
	t.Helper()
	mgr := plugin.NewManager(nil)
	return loader.NewLoader(pluginDir, nil, mgr)
}

// --- Tests de carga básica ---

func TestLoader_LoadAll_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	ldr := newTestLoader(t, dir)

	plugins, err := ldr.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll() en directorio vacío devolvió error: %v", err)
	}
	if len(plugins) != 0 {
		t.Errorf("expected 0 plugins, got %d", len(plugins))
	}
}

func TestLoader_LoadAll_BasicPlugin(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: TestPlugin
version: 1.0.0
author: Test
main: index.js`,
		`console.log("plugin cargado");
module = { onEnable: function() {}, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll() devolvió error: %v", err)
	}
	if len(plugins) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(plugins))
	}
	if plugins[0].GetName() != "TestPlugin" {
		t.Errorf("expected 'TestPlugin', got '%s'", plugins[0].GetName())
	}
}

func TestLoader_LoadAll_MissingPluginYml(t *testing.T) {
	dir := t.TempDir()
	// Crear carpeta sin plugin.yml
	os.MkdirAll(filepath.Join(dir, "bad-plugin"), 0755)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	// No debe devolver error global, simplemente saltea el plugin inválido
	if err != nil {
		t.Fatalf("LoadAll() no debe fallar por plugins inválidos: %v", err)
	}
	if len(plugins) != 0 {
		t.Errorf("expected 0 plugins, got %d", len(plugins))
	}
}

func TestLoader_LoadAll_MissingMainScript(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "bad-plugin")
	os.MkdirAll(pluginDir, 0755)
	os.WriteFile(filepath.Join(pluginDir, "plugin.yml"), []byte(`name: BadPlugin
version: 1.0.0
main: noexiste.js`), 0644)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll() no debe fallar globalmente: %v", err)
	}
	if len(plugins) != 0 {
		t.Errorf("expected 0 plugins (plugin inválido), got %d", len(plugins))
	}
}

func TestLoader_LoadAll_InvalidJS(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: BadJS
version: 1.0.0
main: index.js`,
		`esto no es javascript válido !!!@#$`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll() no debe fallar globalmente por JS inválido: %v", err)
	}
	if len(plugins) != 0 {
		t.Errorf("expected 0 plugins (JS inválido), got %d", len(plugins))
	}
}

func TestLoader_LoadAll_MissingNameInYml(t *testing.T) {
	dir := t.TempDir()
	pluginDir := filepath.Join(dir, "no-name")
	os.MkdirAll(pluginDir, 0755)
	os.WriteFile(filepath.Join(pluginDir, "plugin.yml"), []byte(`version: 1.0.0
main: index.js`), 0644)
	os.WriteFile(filepath.Join(pluginDir, "index.js"), []byte(`module = {};`), 0644)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll() no debe fallar globalmente: %v", err)
	}
	if len(plugins) != 0 {
		t.Errorf("expected 0 plugins (sin nombre), got %d", len(plugins))
	}
}

// --- Tests de OnEnable / OnDisable ---

func TestLoader_OnEnable_Called(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: EnablePlugin
version: 1.0.0
main: index.js`,
		`var enabled = false;
function onEnable() { enabled = true; }
function onDisable() {}
module = { onEnable: onEnable, onDisable: onDisable };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v, plugins: %d", err, len(plugins))
	}

	// onEnable no debe panic
	plugins[0].OnEnable()
}

func TestLoader_OnDisable_Called(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: DisablePlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {}
function onDisable() {}
module = { onEnable: onEnable, onDisable: onDisable };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
	plugins[0].OnDisable() // no debe panic
}

func TestLoader_OnEnable_WithError_DoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ErrorPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() { throw new Error("error intencional"); }
module = { onEnable: onEnable };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}

	// No debe panic aunque el JS lance una excepción
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("OnEnable() causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

// --- Tests de plugin.yml campos ---

func TestLoader_PluginMetadata(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: MetaPlugin
version: 2.5.1
author: TestAuthor
description: Un plugin de prueba
main: index.js`,
		`module = { onEnable: function() {}, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}

	p := plugins[0]
	if p.GetName() != "MetaPlugin" {
		t.Errorf("name: expected 'MetaPlugin', got '%s'", p.GetName())
	}
}

// --- Tests de console ---

func TestLoader_Console_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ConsolePlugin
version: 1.0.0
main: index.js`,
		`console.log("log test");
console.info("info test");
console.warn("warn test");
console.error("error test");
module = { onEnable: function() {}, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
}

// --- Tests de timers ---

func TestLoader_SetTimeout_Executes(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: TimerPlugin
version: 1.0.0
main: index.js`,
		`var fired = false;
setTimeout(function() { fired = true; }, 10);
module = { onEnable: function() {}, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	// Esperar que el timeout se dispare
	time.Sleep(100 * time.Millisecond)
	// No necesitamos verificar `fired` desde Go (es estado JS interno),
	// solo que no paniquee ni genere errores
}

func TestLoader_SetInterval_DoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	// Usamos un interval largo (10s) para que no se dispare durante el test
	// y así evitar que la goroutine acceda a la VM después de que el test termine.
	writePlugin(t, dir,
		`name: IntervalPlugin
version: 1.0.0
main: index.js`,
		`var count = 0;
var id = setInterval(function() { count++; }, 10000);
module = { onEnable: function() {}, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
	// Solo verificamos que setInterval no causa panic al registrarse
	plugins[0].OnDisable()
}

// --- Tests de plugin.info en JS ---

func TestLoader_PluginInfo_AccessibleInJS(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: InfoPlugin
version: 3.0.0
author: InfoAuthor
main: index.js`,
		`if (plugin.name !== "InfoPlugin") throw new Error("nombre incorrecto: " + plugin.name);
if (plugin.version !== "3.0.0") throw new Error("versión incorrecta: " + plugin.version);
if (plugin.author !== "InfoAuthor") throw new Error("autor incorrecto: " + plugin.author);
module = { onEnable: function() {}, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll() falló (plugin.info no accesible): %v", err)
	}
	if len(plugins) != 1 {
		t.Fatalf("expected 1 plugin, got %d", len(plugins))
	}
}

// --- Tests de events.on con evento desconocido ---

func TestLoader_Events_UnknownEvent_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: UnknownEventPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    events.on("EventoQueNoExiste", function(e) {});
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	// No debe panic al registrar evento desconocido
	plugins[0].OnEnable()
}

// --- Tests de commands.register ---

func TestLoader_Commands_Register_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: CommandPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    commands.register("hola", "Dice hola", function(player, args) {
        player.sendMessage("¡Hola!");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	// No debe panic al registrar el comando
	plugins[0].OnEnable()
}

func TestLoader_Commands_Register_WithAliases_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: AliasCommandPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    commands.register("spawn", "Ir al spawn", ["sp", "hub"], function(player, args) {
        player.sendMessage("Teleportando al spawn...");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Commands_Register_WithoutCallback_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: BadCommandPlugin
version: 1.0.0
main: index.js`,
		// Pasar string en lugar de función como callback — no debe panic
		`function onEnable() {
    commands.register("bad", "comando malo", "no soy una función");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("OnEnable() con callback inválido causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_Commands_Register_TooFewArgs_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: FewArgsPlugin
version: 1.0.0
main: index.js`,
		// Llamar commands.register con menos de 3 argumentos — no debe panic
		`function onEnable() {
    commands.register("solo");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("OnEnable() con pocos args causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

// --- Tests de config ---

func TestLoader_Config_SetAndGet(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ConfigPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    config.set("clave", "valor");
    var v = config.get("clave", "default");
    if (v !== "valor") throw new Error("config.get falló: " + v);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Config_Default_WhenMissing(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ConfigDefaultPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var v = config.get("no_existe", "mi_default");
    if (v !== "mi_default") throw new Error("default no funcionó: " + v);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Config_GetString(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ConfigStringPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    config.set("msg", "hola mundo");
    var v = config.getString("msg", "default");
    if (v !== "hola mundo") throw new Error("getString falló: " + v);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Config_GetBool(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ConfigBoolPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    config.set("debug", true);
    var v = config.getBool("debug", false);
    if (v !== true) throw new Error("getBool falló: " + v);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Config_SetDefaults_DoesNotOverwrite(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ConfigDefaultsPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    config.set("prefix", "valor_existente");
    config.setDefaults({ "prefix": "valor_default", "otro": "nuevo" });
    var prefix = config.getString("prefix", "");
    if (prefix !== "valor_existente") throw new Error("setDefaults sobreescribió: " + prefix);
    var otro = config.getString("otro", "");
    if (otro !== "nuevo") throw new Error("setDefaults no agregó clave nueva: " + otro);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Config_Save_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	// writePlugin crea la carpeta "test-plugin" — el dataFolder del plugin es esa ruta
	writePlugin(t, dir,
		`name: ConfigSavePlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    config.set("key", "value");
    config.save();
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()

	// El dataFolder del plugin es dir/test-plugin (nombre de la carpeta en disco, no el name del plugin)
	configPath := filepath.Join(dir, "test-plugin", "config.yml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("config.save() no creó el archivo en: %s", configPath)
	}
}

func TestLoader_Config_Reload_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ConfigReloadPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    config.set("key", "value");
    config.save();
    config.reload();
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("config.reload() causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_Config_GetInt_Default(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ConfigIntPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var v = config.getInt("no_existe", 42);
    if (v !== 42) throw new Error("getInt default falló: " + v);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

// --- Tests de sendTitle ---

func TestLoader_SendTitle_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: TitlePlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    // sendTitle no puede llamarse sin jugador real, pero debe estar disponible como función
    if (typeof sendTitle !== "undefined") throw new Error("sendTitle no debe ser global");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

// --- Tests de efectos de poción ---

func TestLoader_AddEffect_UnknownEffect_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: EffectPlugin
version: 1.0.0
main: index.js`,
		// addEffect con efecto desconocido no debe panic (solo imprime warning)
		`function onEnable() {
    // No tenemos jugador real, solo verificamos que addEffect existe como API global en el playerWrapper
    // La función existe pero no se puede llamar sin jugador
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_PlayerWrapper_HasEffectMethods(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: EffectAPIPlugin
version: 1.0.0
main: index.js`,
		// Verificar que los métodos de efectos se registran en el evento PlayerJoin
		// sin necesitar jugador real — solo registrar el handler no debe fallar
		`function onEnable() {
    events.on("PlayerJoin", function(event) {
        var p = event.getPlayer();
        if (typeof p.addEffect !== "function") throw new Error("addEffect no es función");
        if (typeof p.removeEffect !== "function") throw new Error("removeEffect no es función");
        if (typeof p.clearEffects !== "function") throw new Error("clearEffects no es función");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	// Solo registrar el evento no debe panic
	plugins[0].OnEnable()
}

// --- Tests de armadura ---

func TestLoader_PlayerWrapper_HasArmourMethods(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ArmourAPIPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    events.on("PlayerJoin", function(event) {
        var p = event.getPlayer();
        if (typeof p.setArmour !== "function") throw new Error("setArmour no es función");
        if (typeof p.getArmour !== "function") throw new Error("getArmour no es función");
        if (typeof p.clearArmour !== "function") throw new Error("clearArmour no es función");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_PlayerWrapper_HasTitleMethod(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: TitleAPIPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    events.on("PlayerJoin", function(event) {
        var p = event.getPlayer();
        if (typeof p.sendTitle !== "function") throw new Error("sendTitle no es función");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

// --- Tests de world/server desde comandos (BuildWorldMapFromTx) ---

func TestLoader_Command_WorldGetEntities_IsAvailable(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: CmdWorldPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    commands.register("testworld", "test", function(player, args) {
        // Desde un comando, world debe tener getEntities
        if (typeof world.getEntities !== "function") throw new Error("world.getEntities no es función desde comando");
        if (typeof world.getEntitiesInRadius !== "function") throw new Error("world.getEntitiesInRadius no es función desde comando");
        if (typeof world.removeEntityByUUID !== "function") throw new Error("world.removeEntityByUUID no es función desde comando");
        if (typeof world.setBlock !== "function") throw new Error("world.setBlock no es función desde comando");
        if (typeof world.getBlock !== "function") throw new Error("world.getBlock no es función desde comando");
        if (typeof world.getHighestBlock !== "function") throw new Error("world.getHighestBlock no es función desde comando");
        if (typeof world.spawnEntity !== "function") throw new Error("world.spawnEntity no es función desde comando");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Command_WorldGetEntities_NoTx_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: CmdWorldEmptyPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    commands.register("testentities", "test", function(player, args) {
        var entities = world.getEntities();
        if (!Array.isArray(entities)) throw new Error("getEntities debe retornar array desde comando");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Command_WorldGetEntitiesInRadius_NoTx_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: CmdWorldRadiusPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    commands.register("testradius", "test", function(player, args) {
        var entities = world.getEntitiesInRadius(0, 64, 0, 10);
        if (!Array.isArray(entities)) throw new Error("getEntitiesInRadius debe retornar array desde comando");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Command_WorldRemoveEntityByUUID_NoTx_ReturnsFalse(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: CmdWorldRemovePlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    commands.register("testremove", "test", function(player, args) {
        var removed = world.removeEntityByUUID("00000000-0000-0000-0000-000000000000");
        if (typeof removed !== "boolean") throw new Error("removeEntityByUUID debe retornar boolean desde comando");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Command_ServerGetPlayers_IsAvailable(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: CmdServerPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    commands.register("testserver", "test", function(player, args) {
        if (typeof server.getPlayers !== "function") throw new Error("server.getPlayers no es función desde comando");
        if (typeof server.getPlayer !== "function") throw new Error("server.getPlayer no es función desde comando");
        if (typeof server.broadcast !== "function") throw new Error("server.broadcast no es función desde comando");
        if (typeof server.broadcastTitle !== "function") throw new Error("server.broadcastTitle no es función desde comando");
        if (typeof server.getPlayerCount !== "function") throw new Error("server.getPlayerCount no es función desde comando");
        if (typeof server.getMaxPlayers !== "function") throw new Error("server.getMaxPlayers no es función desde comando");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Command_WorldSpawnEntity_IsAvailable(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: CmdSpawnPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    commands.register("testspawn", "test", function(player, args) {
        if (typeof world.spawnEntity !== "function") throw new Error("world.spawnEntity no es función desde comando");
        // sin tx activo, no debe panic
        world.spawnEntity("lightning", 0, 64, 0);
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.spawnEntity desde comando causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

// --- Tests de Inventory API ---

func TestLoader_World_GetInventory_IsAvailable(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: InvAPIPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    if (typeof world.getInventory !== "function") throw new Error("world.getInventory no es función");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_World_GetInventory_NoServer_ReturnsNull(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: InvNullPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var inv = world.getInventory(0, 64, 0);
    if (inv !== null) throw new Error("sin servidor getInventory debe retornar null");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Player_GetInventory_IsAvailable(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: PlayerInvPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    events.on("PlayerJoin", function(event) {
        var p = event.getPlayer();
        if (typeof p.getInventory !== "function") throw new Error("player.getInventory no es función");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Command_GetInventory_IsAvailable(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: CmdInvPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    commands.register("testinv", "test", function(player, args) {
        if (typeof world.getInventory !== "function") throw new Error("world.getInventory no es función desde comando");
        if (typeof player.getInventory !== "function") throw new Error("player.getInventory no es función desde comando");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Command_WorldGetInventory_NoTx_ReturnsNull(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: CmdInvNullPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    commands.register("testinvnull", "test", function(player, args) {
        var inv = world.getInventory(0, 64, 0);
        // sin tx activo retorna null — no panic
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.getInventory sin tx causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

// --- Tests de Virtual Inventory API ---

func TestLoader_VirtualInventory_CreateMenu_IsAvailable(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: MenuPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    if (typeof inventory === "undefined") throw new Error("inventory global no está definido");
    if (typeof inventory.createMenu !== "function") throw new Error("inventory.createMenu no es función");

    var menu = inventory.createMenu({ title: "Test", type: "chest", size: 27 });
    if (typeof menu.setItems !== "function") throw new Error("menu.setItems no es función");
    if (typeof menu.pattern !== "function") throw new Error("menu.pattern no es función");
    if (typeof menu.onClick !== "function") throw new Error("menu.onClick no es función");
    if (typeof menu.onClose !== "function") throw new Error("menu.onClose no es función");
    if (typeof menu.open !== "function") throw new Error("menu.open no es función");
    if (typeof menu.update !== "function") throw new Error("menu.update no es función");
    if (typeof menu.close !== "function") throw new Error("menu.close no es función");

    menu.setItems([{ slot: 0, name: "minecraft:diamond", count: 1 }]);
    menu.pattern([
        "_________",
        "__xxx____"
    ], {
        x: { name: "minecraft:gold_ingot", count: 3 }
    });
    menu.onClick(function(player, item, click) {});
    menu.onClose(function(player) {});
    menu.open("Notch");
    menu.update("Notch");
    menu.close("Notch");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

// --- Tests de Server API ---

func TestLoader_Server_IsAvailable(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ServerPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    if (typeof server === "undefined") throw new Error("server no está definido");
    if (typeof server.getPlayers !== "function") throw new Error("server.getPlayers no es función");
    if (typeof server.getPlayerCount !== "function") throw new Error("server.getPlayerCount no es función");
    if (typeof server.getMaxPlayers !== "function") throw new Error("server.getMaxPlayers no es función");
    if (typeof server.getPlayer !== "function") throw new Error("server.getPlayer no es función");
    if (typeof server.getPlayerByXUID !== "function") throw new Error("server.getPlayerByXUID no es función");
    if (typeof server.broadcast !== "function") throw new Error("server.broadcast no es función");
    if (typeof server.broadcastTitle !== "function") throw new Error("server.broadcastTitle no es función");
    if (typeof server.getName !== "function") throw new Error("server.getName no es función");
    if (typeof server.shutdown !== "function") throw new Error("server.shutdown no es función");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Server_GetPlayers_NoServer_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ServerPlayersPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var players = server.getPlayers();
    if (!Array.isArray(players)) throw new Error("getPlayers debe retornar array");
    if (players.length !== 0) throw new Error("sin servidor debe retornar array vacío");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Server_GetPlayerCount_NoServer_ReturnsZero(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ServerCountPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var count = server.getPlayerCount();
    if (typeof count !== "number") throw new Error("getPlayerCount debe retornar number");
    if (count !== 0) throw new Error("sin servidor debe retornar 0, got: " + count);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Server_GetMaxPlayers_NoServer_ReturnsZero(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ServerMaxPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var max = server.getMaxPlayers();
    if (typeof max !== "number") throw new Error("getMaxPlayers debe retornar number");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Server_GetPlayer_NoServer_ReturnsNull(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ServerGetPlayerPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var p = server.getPlayer("Notch");
    if (p !== null) throw new Error("sin servidor getPlayer debe retornar null");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Server_GetPlayerByXUID_NoServer_ReturnsNull(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ServerGetXUIDPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var p = server.getPlayerByXUID("123456789");
    if (p !== null) throw new Error("sin servidor getPlayerByXUID debe retornar null");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Server_Broadcast_NoServer_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ServerBroadcastPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    server.broadcast("§aHola a todos!");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("server.broadcast sin servidor causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_Server_BroadcastTitle_NoServer_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ServerBroadcastTitlePlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    server.broadcastTitle("§aEvento!", "§7Comenzando ahora");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("server.broadcastTitle sin servidor causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_Server_GetName_ReturnsPluginName(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: MiPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var name = server.getName();
    if (typeof name !== "string") throw new Error("getName debe retornar string");
    if (name !== "MiPlugin") throw new Error("getName debe retornar 'MiPlugin', got: " + name);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Server_Shutdown_NoServer_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: ServerShutdownPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    server.shutdown(); // sin servidor real, debe no hacer nada
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("server.shutdown sin servidor causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

// --- Tests de Entity API (world.getEntities, getEntitiesInRadius, removeEntityByUUID) ---

func TestLoader_World_GetEntities_IsAvailable(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: EntityAPIPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    if (typeof world.getEntities !== "function") throw new Error("world.getEntities no es función");
    if (typeof world.getEntitiesInRadius !== "function") throw new Error("world.getEntitiesInRadius no es función");
    if (typeof world.removeEntityByUUID !== "function") throw new Error("world.removeEntityByUUID no es función");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_World_GetEntities_NoServer_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: EntityEmptyPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var entities = world.getEntities();
    if (!Array.isArray(entities)) throw new Error("getEntities debe retornar array");
    if (entities.length !== 0) throw new Error("sin servidor debe retornar array vacío");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_World_GetEntitiesInRadius_NoServer_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: EntityRadiusPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var entities = world.getEntitiesInRadius(0, 64, 0, 10);
    if (!Array.isArray(entities)) throw new Error("getEntitiesInRadius debe retornar array");
    if (entities.length !== 0) throw new Error("sin servidor debe retornar array vacío");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_World_RemoveEntityByUUID_NoServer_ReturnsFalse(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: EntityRemovePlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var removed = world.removeEntityByUUID("00000000-0000-0000-0000-000000000000");
    if (typeof removed !== "boolean") throw new Error("removeEntityByUUID debe retornar boolean");
    if (removed !== false) throw new Error("sin servidor debe retornar false");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_Entity_HasBaseMethods_ViaPlayerJoin(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: EntityWrapperPlugin
version: 1.0.0
main: index.js`,
		// Verificar que el entityWrapper tiene los métodos base
		// Lo hacemos via PlayerJoin que expone un playerWrapper (que también es una entity)
		`function onEnable() {
    events.on("PlayerJoin", function(event) {
        var p = event.getPlayer();
        // El playerWrapper tiene todos los métodos del entityWrapper + los de player
        if (typeof p.getX !== "function") throw new Error("getX no es función");
        if (typeof p.getY !== "function") throw new Error("getY no es función");
        if (typeof p.getZ !== "function") throw new Error("getZ no es función");
        if (typeof p.getUUID !== "function") throw new Error("getUUID no es función");
    });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

// --- Tests de World API ---

func TestLoader_World_IsAvailable(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldPlugin
version: 1.0.0
main: index.js`,
		// Verificar que el objeto world existe y tiene las funciones esperadas
		`function onEnable() {
    if (typeof world === "undefined") throw new Error("world no está definido");
    if (typeof world.setBlock !== "function") throw new Error("world.setBlock no es función");
    if (typeof world.getBlock !== "function") throw new Error("world.getBlock no es función");
    if (typeof world.getHighestBlock !== "function") throw new Error("world.getHighestBlock no es función");
    if (typeof world.spawnEntity !== "function") throw new Error("world.spawnEntity no es función");
    if (typeof world.spawnParticle !== "function") throw new Error("world.spawnParticle no es función");
    if (typeof world.getPlayers !== "function") throw new Error("world.getPlayers no es función");
    if (typeof world.getPlayerCount !== "function") throw new Error("world.getPlayerCount no es función");
    if (typeof world.broadcast !== "function") throw new Error("world.broadcast no es función");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_World_SetBlock_NoServer_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldSetBlockPlugin
version: 1.0.0
main: index.js`,
		// Sin servidor real, setBlock debe retornar silenciosamente (no panic)
		`function onEnable() {
    world.setBlock(0, 64, 0, "minecraft:stone");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.setBlock sin servidor causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_World_GetBlock_NoServer_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldGetBlockPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var b = world.getBlock(0, 64, 0);
    if (typeof b !== "string") throw new Error("getBlock debe retornar string, got: " + typeof b);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_World_GetHighestBlock_NoServer_ReturnsZero(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldHighestPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var y = world.getHighestBlock(0, 0);
    if (typeof y !== "number") throw new Error("getHighestBlock debe retornar number, got: " + typeof y);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_World_SpawnEntity_Lightning_NoServer_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldLightningPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    world.spawnEntity("lightning", 0, 64, 0);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.spawnEntity lightning sin servidor causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_World_SpawnEntity_TNT_NoServer_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldTNTPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    world.spawnEntity("tnt", 0, 64, 0, { fuse: 4 });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.spawnEntity tnt sin servidor causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_World_SpawnEntity_Text_NoServer_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldTextPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    world.spawnEntity("text", 0, 64, 0, { text: "§aHola Mundo" });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.spawnEntity text sin servidor causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_World_SpawnEntity_ExperienceOrb_NoServer_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldXPPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    world.spawnEntity("experience_orb", 0, 64, 0, { amount: 10 });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.spawnEntity experience_orb sin servidor causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_World_SpawnEntity_Item_NoServer_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldItemPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    world.spawnEntity("item", 0, 64, 0, { item: "minecraft:diamond", count: 5 });
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.spawnEntity item sin servidor causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_World_SpawnEntity_Unknown_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldUnknownEntityPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    // tipo desconocido — solo imprime warning, no panic
    world.spawnEntity("entidad_que_no_existe", 0, 64, 0);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.spawnEntity tipo desconocido causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_World_SpawnEntity_NoOptions_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldNoOptsPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    // sin opciones — debe usar defaults
    world.spawnEntity("lightning", 0, 64, 0);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.spawnEntity sin opciones causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_World_SpawnParticle_Unknown_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldParticlePlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    // partícula desconocida — solo imprime warning, no panic
    world.spawnParticle(0, 64, 0, "partícula_que_no_existe");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.spawnParticle desconocida causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_World_GetPlayers_NoServer_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldPlayersPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var players = world.getPlayers();
    if (!Array.isArray(players)) throw new Error("getPlayers debe retornar array");
    if (players.length !== 0) throw new Error("sin servidor debe retornar array vacío");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_World_GetPlayerCount_NoServer_ReturnsZero(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldCountPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    var count = world.getPlayerCount();
    if (typeof count !== "number") throw new Error("getPlayerCount debe retornar number");
    if (count !== 0) throw new Error("sin servidor debe retornar 0, got: " + count);
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	plugins[0].OnEnable()
}

func TestLoader_World_Broadcast_NoServer_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldBroadcastPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    world.broadcast("§aHola a todos!");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.broadcast sin servidor causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

func TestLoader_World_SetBlock_UnknownBlock_NoError(t *testing.T) {
	dir := t.TempDir()
	writePlugin(t, dir,
		`name: WorldUnknownBlockPlugin
version: 1.0.0
main: index.js`,
		`function onEnable() {
    // bloque desconocido — solo imprime warning, no panic
    world.setBlock(0, 64, 0, "minecraft:bloque_que_no_existe");
}
module = { onEnable: onEnable, onDisable: function() {} };`,
	)

	ldr := newTestLoader(t, dir)
	plugins, err := ldr.LoadAll()
	if err != nil || len(plugins) != 1 {
		t.Fatalf("LoadAll() falló: %v", err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("world.setBlock con bloque desconocido causó panic: %v", r)
		}
	}()
	plugins[0].OnEnable()
}

// --- Tests de GetPlugins ---

func TestLoader_GetPlugins(t *testing.T) {
	dir := t.TempDir()

	// Crear dos plugins
	for _, name := range []string{"PluginA", "PluginB"} {
		pluginDir := filepath.Join(dir, name)
		os.MkdirAll(pluginDir, 0755)
		os.WriteFile(filepath.Join(pluginDir, "plugin.yml"), []byte("name: "+name+"\nversion: 1.0.0\nmain: index.js"), 0644)
		os.WriteFile(filepath.Join(pluginDir, "index.js"), []byte(`module = { onEnable: function(){}, onDisable: function(){} };`), 0644)
	}

	ldr := newTestLoader(t, dir)
	_, err := ldr.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll() falló: %v", err)
	}

	jsPlugins := ldr.GetPlugins()
	if len(jsPlugins) != 2 {
		t.Errorf("GetPlugins() expected 2, got %d", len(jsPlugins))
	}
}
