# Plugin completo de ejemplo

Este es el plugin de ejemplo incluido en `plugins/example/`. Fue verificado con un test manual real en Minecraft Bedrock 1.21.50 y demuestra todos los sistemas disponibles.

## Archivos

### `plugin.yml`

```yaml
name: ExamplePlugin
version: 1.0.0
author: nzxsww
description: Plugin de ejemplo para Dragonfly Script API
main: index.js
api-version: 1.0.0
```

### `index.js`

```js
// Plugin de ejemplo — Dragonfly Script API
// Demuestra: eventos, comandos, config YAML, inventario y sonidos.

console.log("Cargando plugin: " + plugin.name + " v" + plugin.version);

// =============================================
// Configuración por defecto
// Se guarda en plugins/example/config.yml
// =============================================
config.setDefaults({
    "welcome-message": "¡Bienvenido al servidor!",
    "prefix":          "§a[Servidor]§r",
    "protected-y":     10,
    "chat-format":     "§7[§f{name}§7] §r{msg}"
});

// =============================================
// Ciclo de vida
// =============================================

function onEnable() {
    console.log("Plugin habilitado!");

    // Leer config (usa defaults si no hay config.yml)
    var welcomeMsg  = config.getString("welcome-message", "¡Bienvenido!");
    var prefix      = config.getString("prefix", "§f[Servidor]§r");
    var protectedY  = config.getInt("protected-y", 10);
    var chatFormat  = config.getString("chat-format", "§7[§f{name}§7] §r{msg}");

    console.log("Config cargada — protected-y: " + protectedY);

    // =============================================
    // Eventos
    // =============================================

    // Jugador entra al servidor
    events.on("PlayerJoin", function(event) {
        var player = event.getPlayer();
        console.log(player.getName() + " se unió al servidor");
        player.sendMessage(prefix + " " + welcomeMsg);
        player.playSound("click");
        event.setJoinMessage("§e" + player.getName() + " §fentra al servidor.");
    });

    // Jugador sale del servidor
    events.on("PlayerQuit", function(event) {
        var player = event.getPlayer();
        console.log(player.getName() + " salió del servidor");
        event.setQuitMessage("§c" + player.getName() + " §fsale del servidor.");
    });

    // Jugador escribe en el chat — formato personalizado
    events.on("PlayerChat", function(event) {
        var player = event.getPlayer();
        var msg = chatFormat
            .replace("{name}", player.getName())
            .replace("{msg}", event.getMessage());
        event.setMessage(msg);
    });

    // Jugador rompe un bloque — proteger zona baja
    events.on("BlockBreak", function(event) {
        var player = event.getPlayer();
        if (event.getBlockY() < protectedY) {
            event.setCancelled(true);
            player.sendMessage("§cNo puedes romper bloques tan abajo!");
            player.playSound("deny");
        }
    });

    // Jugador coloca un bloque — proteger zona baja
    events.on("BlockPlace", function(event) {
        var player = event.getPlayer();
        if (event.getBlockY() < protectedY) {
            event.setCancelled(true);
            player.sendMessage("§cNo puedes colocar bloques aquí!");
            player.playSound("deny");
        }
    });

    // =============================================
    // Comandos
    // =============================================

    commands.register("hola", "Te saluda por tu nombre", function(player, args) {
        player.sendMessage("§aHola, " + player.getName() + "!");
        player.playSound("pop");
    });

    commands.register("pos", "Muestra tu posición actual", function(player, args) {
        var x = Math.floor(player.getX());
        var y = Math.floor(player.getY());
        var z = Math.floor(player.getZ());
        player.sendMessage("§7Posición: §f" + x + ", " + y + ", " + z);
    });

    commands.register("salud", "Muestra tu vida actual", function(player, args) {
        var hp    = Math.floor(player.getHealth());
        var maxHp = Math.floor(player.getMaxHealth());
        player.sendMessage("§cVida: §f" + hp + "/" + maxHp);
    });

    commands.register("give", "Da un item. Uso: /give <item> [cantidad]", function(player, args) {
        if (args.length < 1) {
            player.sendMessage("§cUso: /give <item> [cantidad]");
            return;
        }
        var item  = args[0];
        var count = args.length >= 2 ? parseInt(args[1]) : 1;
        if (isNaN(count) || count < 1) count = 1;
        if (count > 64) count = 64;

        var ok = player.giveItem(item, count);
        if (ok) {
            player.sendMessage("§aSe te dieron §f" + count + "x §a" + item);
            player.playSound("pop");
        } else {
            player.sendMessage("§cItem desconocido o inventario lleno: " + item);
            player.playSound("deny");
        }
    });

    commands.register("creative", "Cambia a modo creativo", function(player, args) {
        player.setGameMode("creative");
        player.sendMessage("§aModo creativo activado.");
        player.playSound("levelup");
    });

    commands.register("survival", "Cambia a modo supervivencia", function(player, args) {
        player.setGameMode("survival");
        player.sendMessage("§aModo supervivencia activado.");
    });

    commands.register("fly", "Activa o desactiva el vuelo", function(player, args) {
        var mode = player.getGameMode();
        if (mode === "survival" || mode === "adventure") {
            // startFlying solo funciona en creative/spectator
            player.setGameMode("creative");
            player.startFlying();
            player.sendMessage("§aVuelo activado. §7(modo creativo)");
        } else if (player.isFlying()) {
            player.stopFlying();
            player.sendMessage("§cVuelo desactivado.");
        } else {
            player.startFlying();
            player.sendMessage("§aVuelo activado.");
        }
    });

    // Anuncio de que el plugin lleva 30 segundos activo
    setTimeout(function() {
        console.log("El plugin lleva 30 segundos activo.");
    }, 30000);
}

function onDisable() {
    config.save(); // guardar config al cerrar
    console.log("Plugin deshabilitado — config guardada.");
}

// Exportar ciclo de vida (obligatorio)
module = {
    onEnable:  onEnable,
    onDisable: onDisable
};
```

## Qué demuestra este plugin

| Característica | Dónde se ve |
|---|---|
| Config YAML con defaults | `config.setDefaults()` al inicio |
| Leer config en onEnable | `config.getString()`, `config.getInt()` |
| Evento PlayerJoin | Mensaje de bienvenida + sonido |
| Evento PlayerQuit | Mensaje de salida personalizado |
| Evento PlayerChat | Formato de chat personalizado |
| Evento BlockBreak/Place | Zona protegida bajo Y=10 |
| Comandos simples | `/hola`, `/pos`, `/salud` |
| Comando con argumentos | `/give <item> [cantidad]` |
| Cambio de gamemode | `/creative`, `/survival` |
| Control de vuelo | `/fly` (compatible con survival) |
| Sonidos | `playSound()` en varios eventos |
| setTimeout | Log a los 30 segundos |
| Guardar config al cerrar | `config.save()` en onDisable |

## `config.yml` generado

Al cerrar el servidor por primera vez, se crea automáticamente:

```yaml
welcome-message: ¡Bienvenido al servidor!
prefix: §a[Servidor]§r
protected-y: 10
chat-format: §7[§f{name}§7] §r{msg}
```

Podés editar este archivo y los cambios se aplican al reiniciar el servidor.
