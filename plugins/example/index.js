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

    commands.register("armadura", "Equipa armadura de diamante completa", function(player, args) {
        player.setArmour(0, "minecraft:diamond_helmet");
        player.setArmour(1, "minecraft:diamond_chestplate");
        player.setArmour(2, "minecraft:diamond_leggings");
        player.setArmour(3, "minecraft:diamond_boots");
        player.sendMessage("§bArmadura de diamante equipada!");
        player.sendTitle("§b§lArmadura", "§7diamond completa");
        player.playSound("anvil_land");
        console.log("[TEST] setArmour: " + player.getName() + " equipó armadura de diamante");
    });

    commands.register("quitararmadura", "Quita toda la armadura", function(player, args) {
        player.clearArmour();
        player.sendMessage("§7Armadura quitada.");
        console.log("[TEST] clearArmour: " + player.getName());
    });

    commands.register("verarmadura", "Muestra la armadura equipada", function(player, args) {
        var casco     = player.getArmour(0) || "ninguno";
        var pechera   = player.getArmour(1) || "ninguna";
        var pantalon  = player.getArmour(2) || "ninguno";
        var botas     = player.getArmour(3) || "ninguna";
        player.sendMessage("§eCasco: §f"    + casco);
        player.sendMessage("§ePechera: §f"  + pechera);
        player.sendMessage("§ePantalón: §f" + pantalon);
        player.sendMessage("§eBotas: §f"    + botas);
    });

    commands.register("efecto", "Aplica speed II + regeneration I por 30s", function(player, args) {
        player.addEffect("speed", 2, 30);
        player.addEffect("regeneration", 1, 30);
        player.sendMessage("§aEfectos aplicados: §fSpeed II + Regeneration I (30s)");
        player.sendTitle("§a§lEfectos!", "§7Speed II + Regen I");
        player.playSound("levelup");
        console.log("[TEST] addEffect: " + player.getName() + " speed+regen");
    });

    commands.register("quitarefectos", "Quita todos los efectos activos", function(player, args) {
        player.clearEffects();
        player.sendMessage("§7Todos los efectos eliminados.");
        console.log("[TEST] clearEffects: " + player.getName());
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

    // --- TEST: PlayerMove ---
    events.on("PlayerMove", function(event) {
        // Solo loguear en servidor, demasiado frecuente para title
        var p = event.getPlayer();
        // console.log("[TEST] PlayerMove: " + p.getName() + " -> " + Math.floor(event.getToX()) + "," + Math.floor(event.getToY()) + "," + Math.floor(event.getToZ()));
    });

    // --- TEST: PlayerJump ---
    events.on("PlayerJump", function(event) {
        var p = event.getPlayer();
        console.log("[TEST] PlayerJump: " + p.getName());
        p.sendTitle("§eSaltaste!", "§7PlayerJump OK");
    });

    // --- TEST: PlayerToggleSprint ---
    events.on("PlayerToggleSprint", function(event) {
        var p = event.getPlayer();
        var sprinting = event.isSprinting();
        console.log("[TEST] PlayerToggleSprint: " + p.getName() + " sprinting=" + sprinting);
        p.sendTitle(sprinting ? "§aCorriendo" : "§7Caminando", "§7PlayerToggleSprint OK");
    });

    // --- TEST: PlayerToggleSneak ---
    events.on("PlayerToggleSneak", function(event) {
        var p = event.getPlayer();
        var sneaking = event.isSneaking();
        console.log("[TEST] PlayerToggleSneak: " + p.getName() + " sneaking=" + sneaking);
        p.sendTitle(sneaking ? "§6Agachado" : "§7De pie", "§7PlayerToggleSneak OK");
    });

    // --- TEST: PlayerHurt ---
    events.on("PlayerHurt", function(event) {
        var p = event.getPlayer();
        var dmg = event.getDamage();
        console.log("[TEST] PlayerHurt: " + p.getName() + " daño=" + dmg);
        p.sendTitle("§cDaño recibido", "§f" + dmg.toFixed(1) + " pts — PlayerHurt OK");
    });

    // --- TEST: PlayerHeal ---
    events.on("PlayerHeal", function(event) {
        var p = event.getPlayer();
        var hp = event.getHealth();
        console.log("[TEST] PlayerHeal: " + p.getName() + " curación=" + hp);
        p.sendTitle("§aRecuperando vida", "§f+" + hp.toFixed(1) + " — PlayerHeal OK");
    });

    // --- TEST: PlayerDeath ---
    events.on("PlayerDeath", function(event) {
        var p = event.getPlayer();
        console.log("[TEST] PlayerDeath: " + p.getName());
        // Activar keepInventory para el test
        event.setKeepInventory(true);
        p.sendMessage("§c[TEST] Moriste — keepInventory activado. PlayerDeath OK");
    });

    // --- TEST: PlayerRespawn ---
    events.on("PlayerRespawn", function(event) {
        var p = event.getPlayer();
        console.log("[TEST] PlayerRespawn: " + p.getName() + " pos=" + event.getX().toFixed(1) + "," + event.getY().toFixed(1) + "," + event.getZ().toFixed(1));
        p.sendTitle("§aRespawn", "§7PlayerRespawn OK");
    });

    // --- TEST: PlayerFoodLoss ---
    events.on("PlayerFoodLoss", function(event) {
        var p = event.getPlayer();
        console.log("[TEST] PlayerFoodLoss: " + p.getName() + " de=" + event.getFrom() + " a=" + event.getTo());
        p.sendTitle("§6Hambre", "§fde " + event.getFrom() + " a " + event.getTo() + " — PlayerFoodLoss OK");
    });

    // --- TEST: PlayerExperienceGain ---
    events.on("PlayerExperienceGain", function(event) {
        var p = event.getPlayer();
        console.log("[TEST] PlayerExperienceGain: " + p.getName() + " xp=" + event.getAmount());
        p.sendTitle("§bXP ganada", "§f+" + event.getAmount() + " — PlayerExperienceGain OK");
    });

    // --- TEST: PlayerTeleport ---
    events.on("PlayerTeleport", function(event) {
        var p = event.getPlayer();
        console.log("[TEST] PlayerTeleport: " + p.getName() + " a=" + event.getX().toFixed(1) + "," + event.getY().toFixed(1) + "," + event.getZ().toFixed(1));
        p.sendTitle("§dTeleport", "§7PlayerTeleport OK");
    });

    // --- TEST: PlayerAttackEntity ---
    events.on("PlayerAttackEntity", function(event) {
        var p = event.getPlayer();
        var crit = event.isCritical();
        console.log("[TEST] PlayerAttackEntity: " + p.getName() + " fuerza=" + event.getForce().toFixed(2) + " critico=" + crit);
        p.sendTitle(crit ? "§c¡Crítico!" : "§eAtaque", "§ffuerza=" + event.getForce().toFixed(2) + " — PlayerAttackEntity OK");
    });

    // --- TEST: PlayerItemUse ---
    events.on("PlayerItemUse", function(event) {
        var p = event.getPlayer();
        console.log("[TEST] PlayerItemUse: " + p.getName());
        p.sendTitle("§aItem usado", "§7PlayerItemUse OK");
    });

    // --- TEST: PlayerItemDrop ---
    events.on("PlayerItemDrop", function(event) {
        var p = event.getPlayer();
        console.log("[TEST] PlayerItemDrop: " + p.getName() + " cantidad=" + event.getItemCount());
        p.sendTitle("§eItem tirado", "§f" + event.getItemCount() + "x — PlayerItemDrop OK");
    });

    // --- TEST: PlayerItemPickup ---
    events.on("PlayerItemPickup", function(event) {
        var p = event.getPlayer();
        console.log("[TEST] PlayerItemPickup: " + p.getName() + " cantidad=" + event.getItemCount());
        p.sendTitle("§aItem recogido", "§f" + event.getItemCount() + "x — PlayerItemPickup OK");
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
