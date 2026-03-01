// Plugin de ejemplo — Dragonfly Script API
// Demuestra el uso de: eventos, comandos, config YAML, inventario y sonidos.

console.log("Cargando plugin: " + plugin.name + " v" + plugin.version);

// === Valores por defecto de la configuración ===
// Estos se escriben en config.yml si no existen todavía.
// El archivo se crea en plugins/example/config.yml
config.setDefaults({
    "welcome-message": "¡Bienvenido al servidor!",
    "prefix":          "§a[Servidor]§r",
    "protected-y":     10,
    "chat-format":     "§7[§f{name}§7] §r{msg}"
});

// === Ciclo de vida del plugin ===

function onEnable() {
    console.log("Plugin habilitado!");

    // Leer config desde el YAML (usa los defaults si no hay config.yml aún)
    var welcomeMsg  = config.getString("welcome-message", "¡Bienvenido!");
    var prefix      = config.getString("prefix", "§f[Servidor]§r");
    var protectedY  = config.getInt("protected-y", 10);
    var chatFormat  = config.getString("chat-format", "§7[§f{name}§7] §r{msg}");

    console.log("Config cargada — protected-y: " + protectedY);

    // --- Evento: jugador entra al servidor ---
    events.on("PlayerJoin", function(event) {
        var player = event.getPlayer();
        console.log(player.getName() + " se unió al servidor");
        player.sendMessage(prefix + " " + welcomeMsg);
        player.playSound("click");
        event.setJoinMessage("§e" + player.getName() + " §fentra al servidor.");
    });

    // --- Evento: jugador sale del servidor ---
    events.on("PlayerQuit", function(event) {
        var player = event.getPlayer();
        console.log(player.getName() + " salió del servidor");
        event.setQuitMessage("§c" + player.getName() + " §fsale del servidor.");
    });

    // --- Evento: jugador escribe en el chat ---
    events.on("PlayerChat", function(event) {
        var player = event.getPlayer();
        var msg = event.getMessage();
        // Aplicar el formato del chat desde config
        var formatted = chatFormat
            .replace("{name}", player.getName())
            .replace("{msg}", msg);
        event.setMessage(formatted);
    });

    // --- Evento: jugador rompe un bloque ---
    events.on("BlockBreak", function(event) {
        var player = event.getPlayer();
        var y = event.getBlockY();
        if (y < protectedY) {
            event.setCancelled(true);
            player.sendMessage("§cNo puedes romper bloques tan abajo!");
            player.playSound("deny");
        }
    });

    // --- Evento: jugador coloca un bloque ---
    events.on("BlockPlace", function(event) {
        var player = event.getPlayer();
        var y = event.getBlockY();
        if (y < protectedY) {
            event.setCancelled(true);
            player.sendMessage("§cNo puedes colocar bloques aquí!");
            player.playSound("deny");
        }
    });

    // --- Comandos ---

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

    commands.register("give", "Da un item al jugador. Uso: /give <item> [cantidad]", function(player, args) {
        if (args.length < 1) {
            player.sendMessage("§cUso: /give <item> [cantidad]");
            return;
        }
        var itemName = args[0];
        var count    = args.length >= 2 ? parseInt(args[1]) : 1;
        if (isNaN(count) || count < 1) count = 1;
        if (count > 64) count = 64;

        var ok = player.giveItem(itemName, count);
        if (ok) {
            player.sendMessage("§aSe te dieron §f" + count + "x §a" + itemName);
            player.playSound("pop");
        } else {
            player.sendMessage("§cItem desconocido o inventario lleno: " + itemName);
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
        // StartFlying/StopFlying solo funciona en gamemodes que permiten volar (creative/spectator).
        // Si el jugador está en survival o adventure, primero lo pasamos a creative.
        var mode = player.getGameMode();
        if (mode === "survival" || mode === "adventure") {
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

    // Log para confirmar que el plugin arrancó correctamente
    setTimeout(function() {
        console.log("El plugin lleva 30 segundos activo.");
    }, 30000);
}

function onDisable() {
    // Guardar config al desactivar (por si se modificó algo en runtime)
    config.save();
    console.log("Plugin deshabilitado — config guardada.");
}

// Exportar ciclo de vida (obligatorio)
module = {
    onEnable:  onEnable,
    onDisable: onDisable
};
