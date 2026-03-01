# Configuración YAML

Cada plugin tiene su propia configuración en un archivo YAML ubicado en `plugins/nombre-plugin/config.yml`. El archivo se crea automáticamente cuando el plugin llama `config.save()`.

## Flujo recomendado

```js
// 1. Definir defaults al inicio del script (fuera de onEnable)
config.setDefaults({
    "prefix":       "§a[MiPlugin]§r",
    "max-jugadores": 20,
    "modo-debug":   false,
    "multiplicador": 1.5
});

function onEnable() {
    // 2. Leer los valores (usa los defaults si no hay config.yml)
    var prefix = config.getString("prefix", "§f[?]");
    console.log("Prefix: " + prefix);
}

function onDisable() {
    // 3. Guardar al cerrar
    config.save();
}
```

## Métodos disponibles

### `config.setDefaults(objeto)`

Define los valores por defecto. **No sobreescribe** valores existentes — solo aplica si la clave no existe todavía.

```js
config.setDefaults({
    "welcome": "¡Bienvenido!",
    "max":     100,
    "debug":   false
});
```

### `config.get(clave, default)`

Obtiene cualquier valor. Retorna `null` si no existe y no se pasó default.

```js
var val = config.get("welcome", "Hola");
```

### `config.getString(clave, default)`

Obtiene un valor como string.

```js
var prefix = config.getString("prefix", "§f[Server]");
```

### `config.getInt(clave, default)`

Obtiene un valor como número entero.

```js
var max = config.getInt("max-jugadores", 20);
```

### `config.getBool(clave, default)`

Obtiene un valor como booleano.

```js
var debug = config.getBool("modo-debug", false);
if (debug) {
    console.log("Modo debug activado");
}
```

### `config.getFloat(clave, default)`

Obtiene un valor como número decimal.

```js
var mult = config.getFloat("multiplicador", 1.0);
```

### `config.set(clave, valor)`

Escribe un valor en memoria. No guarda al disco hasta llamar `save()`.

```js
config.set("ultimo-jugador", player.getName());
config.set("contador", 42);
config.set("activo", true);
```

### `config.save()`

Guarda la configuración en disco (`config.yml`). Crea el archivo si no existe.

```js
config.save();
```

### `config.reload()`

Recarga la configuración desde disco. Útil si alguien editó el archivo manualmente.

```js
config.reload();
```

## Ejemplo completo

```js
config.setDefaults({
    "prefix":           "§a[Servidor]§r",
    "welcome-message":  "¡Bienvenido al servidor!",
    "protected-y":      10,
    "chat-format":      "§7[§f{name}§7] §r{msg}"
});

function onEnable() {
    var prefix     = config.getString("prefix", "§f[?]");
    var welcome    = config.getString("welcome-message", "Bienvenido");
    var protectedY = config.getInt("protected-y", 10);
    var chatFormat = config.getString("chat-format", "{name}: {msg}");

    events.on("PlayerJoin", function(event) {
        var player = event.getPlayer();
        player.sendMessage(prefix + " " + welcome);
    });

    events.on("PlayerChat", function(event) {
        var player = event.getPlayer();
        var msg = chatFormat
            .replace("{name}", player.getName())
            .replace("{msg}", event.getMessage());
        event.setMessage(msg);
    });

    events.on("BlockBreak", function(event) {
        if (event.getBlockY() < protectedY) {
            event.setCancelled(true);
            event.getPlayer().sendMessage("§cZona protegida.");
        }
    });
}

function onDisable() {
    config.save();
}

module = { onEnable: onEnable, onDisable: onDisable };
```

El `config.yml` resultante después de `config.save()`:

```yaml
prefix: §a[Servidor]§r
welcome-message: ¡Bienvenido al servidor!
protected-y: 10
chat-format: §7[§f{name}§7] §r{msg}
```
