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

    // --- Comandos de test World API ---

    commands.register("testbloque", "Coloca y lee un bloque en tu posición", function(player, args) {
        var x = Math.floor(player.getX());
        var y = Math.floor(player.getY()) - 1; // bloque bajo los pies
        var z = Math.floor(player.getZ());
        var bloque = world.getBlock(x, y, z);
        console.log("[TEST] getBlock en " + x + "," + y + "," + z + " = " + bloque);
        player.sendMessage("§eBloque bajo tus pies: §f" + bloque);

        // Colocar un bloque de diamante 3 bloques adelante
        world.setBlock(x + 3, y + 1, z, "minecraft:diamond_block");
        player.sendMessage("§bBloque de diamante colocado en §f" + (x+3) + "," + (y+1) + "," + z);
        console.log("[TEST] setBlock diamond_block en " + (x+3) + "," + (y+1) + "," + z);
    });

    commands.register("testaltura", "Muestra la altura del terreno en tu posición", function(player, args) {
        var x = Math.floor(player.getX());
        var z = Math.floor(player.getZ());
        var y = world.getHighestBlock(x, z);
        console.log("[TEST] getHighestBlock en " + x + "," + z + " = " + y);
        player.sendMessage("§eBloque más alto en §f" + x + "," + z + "§e: §fY=" + y);
        player.sendTitle("§eAltitud", "§fY=" + y + " en X=" + x + " Z=" + z);
    });

    commands.register("testrayo", "Invoca un rayo en tu posición", function(player, args) {
        var x = player.getX();
        var y = player.getY();
        var z = player.getZ();
        world.spawnEntity("lightning", x, y, z);
        player.sendMessage("§e⚡ Rayo invocado!");
        console.log("[TEST] spawnEntity lightning en " + x + "," + y + "," + z);
    });

    commands.register("testtnt", "Coloca un TNT cerca de ti (4s mecha)", function(player, args) {
        var x = player.getX() + 3;
        var y = player.getY();
        var z = player.getZ();
        world.spawnEntity("tnt", x, y, z, { fuse: 4 });
        player.sendMessage("§c💣 TNT activado a 3 bloques! (4 segundos)");
        console.log("[TEST] spawnEntity tnt en " + x + "," + y + "," + z);
    });

    commands.register("testtexto", "Crea un texto flotante en tu posición", function(player, args) {
        var x = player.getX();
        var y = player.getY() + 2;
        var z = player.getZ();
        var texto = args.length > 0 ? args.join(" ") : "§a§lTexto flotante!";
        world.spawnEntity("text", x, y, z, { text: texto });
        player.sendMessage("§aTexto flotante creado!");
        console.log("[TEST] spawnEntity text '" + texto + "' en " + x + "," + y + "," + z);
    });

    commands.register("testxp", "Genera orbes de XP en tu posición", function(player, args) {
        var cantidad = args.length > 0 ? parseInt(args[0]) : 100;
        if (isNaN(cantidad) || cantidad < 1) cantidad = 100;
        world.spawnEntity("experience_orb", player.getX(), player.getY(), player.getZ(), { amount: cantidad });
        player.sendMessage("§b✨ " + cantidad + " XP generados!");
        console.log("[TEST] spawnEntity experience_orb " + cantidad + " en pos del jugador");
    });

    commands.register("testparticula", "Genera partículas en tu posición", function(player, args) {
        var tipo = args.length > 0 ? args[0] : "flame";
        var x = player.getX();
        var y = player.getY() + 1;
        var z = player.getZ();
        world.spawnParticle(x, y, z, tipo);
        player.sendMessage("§dPartícula §f'" + tipo + "' §dgenerada!");
        console.log("[TEST] spawnParticle '" + tipo + "' en " + x + "," + y + "," + z);
    });

    commands.register("testjugadores", "Lista todos los jugadores conectados", function(player, args) {
        var jugadores = world.getPlayers();
        var cantidad = world.getPlayerCount();
        player.sendMessage("§eJugadores conectados: §f" + cantidad);
        for (var i = 0; i < jugadores.length; i++) {
            player.sendMessage("§7- §f" + jugadores[i].getName() + " §7(" + Math.floor(jugadores[i].getX()) + "," + Math.floor(jugadores[i].getY()) + "," + Math.floor(jugadores[i].getZ()) + ")");
        }
        console.log("[TEST] getPlayers: " + cantidad + " jugadores");
    });

    commands.register("broadcast", "Envía un mensaje a todos los jugadores", function(player, args) {
        if (args.length < 1) {
            player.sendMessage("§cUso: /broadcast <mensaje>");
            return;
        }
        var msg = "§6[Broadcast] §f" + args.join(" ");
        world.broadcast(msg);
        console.log("[TEST] broadcast: " + msg);
    });

    // --- Comandos de test Inventory API ---

    commands.register("verinv", "Ver tu inventario", function(player, args) {
        var inv = player.getInventory();
        var items = inv.getItems();
        player.sendMessage("§eInventario (§f" + items.length + "§e slots ocupados de §f" + inv.getSize() + "§e):");
        for (var i = 0; i < items.length; i++) {
            player.sendMessage("§7Slot §f" + items[i].slot + "§7: §f" + items[i].name + " x" + items[i].count);
        }
        console.log("[TEST] player.getInventory: " + items.length + " items");
    });

    commands.register("contaritems", "Contar items en tu inventario", function(player, args) {
        if (args.length < 1) {
            player.sendMessage("§cUso: /contaritems <nombre>");
            return;
        }
        var nombre = args[0];
        var inv = player.getInventory();
        var cantidad = inv.count(nombre);
        var tiene = inv.contains(nombre);
        player.sendMessage("§eItem: §f" + nombre);
        player.sendMessage("§eCantidad total: §f" + cantidad);
        player.sendMessage("§eContiene: §f" + (tiene ? "§aSí" : "§cNo"));
        console.log("[TEST] inv.count(" + nombre + ")=" + cantidad);
    });

    commands.register("daritems", "Agrega items a tu inventario", function(player, args) {
        if (args.length < 1) {
            player.sendMessage("§cUso: /daritems <nombre> [cantidad]");
            return;
        }
        var nombre = args[0];
        var cantidad = args.length > 1 ? parseInt(args[1]) : 1;
        var inv = player.getInventory();
        var sobrante = inv.addItem(nombre, cantidad);
        if (sobrante === 0) {
            player.sendMessage("§a" + cantidad + "x §f" + nombre + " §aagregado!");
        } else {
            player.sendMessage("§eSolo se agregaron §f" + (cantidad - sobrante) + "§e (inventario lleno, sobrante: §f" + sobrante + "§e)");
        }
        console.log("[TEST] inv.addItem(" + nombre + ", " + cantidad + ") sobrante=" + sobrante);
    });

    commands.register("quitaritems", "Quita items de tu inventario", function(player, args) {
        if (args.length < 1) {
            player.sendMessage("§cUso: /quitaritems <nombre> [cantidad]");
            return;
        }
        var nombre = args[0];
        var cantidad = args.length > 1 ? parseInt(args[1]) : 1;
        var inv = player.getInventory();
        var ok = inv.removeItem(nombre, cantidad);
        if (ok) {
            player.sendMessage("§a" + cantidad + "x §f" + nombre + " §aremovido!");
        } else {
            player.sendMessage("§cNo se pudo remover §f" + nombre + " §c(no tienes suficientes)");
        }
        console.log("[TEST] inv.removeItem(" + nombre + ", " + cantidad + ")=" + ok);
    });

    commands.register("limpiartodo", "Vacía tu inventario", function(player, args) {
        var inv = player.getInventory();
        inv.clear();
        player.sendMessage("§aInventario vaciado!");
        console.log("[TEST] inv.clear()");
    });

    commands.register("vercofre", "Ver contenido del cofre frente a ti", function(player, args) {
        var x = Math.floor(player.getX());
        var y = Math.floor(player.getY());
        var z = Math.floor(player.getZ() + 1);
        var inv = world.getInventory(x, y, z);
        if (inv === null) {
            player.sendMessage("§cNo hay un contenedor en §f" + x + "," + y + "," + z);
            console.log("[TEST] world.getInventory: null en " + x + "," + y + "," + z);
            return;
        }
        var items = inv.getItems();
        player.sendMessage("§eCofre en §f" + x + "," + y + "," + z + " §e(§f" + items.length + "§e/§f" + inv.getSize() + " slots):");
        for (var i = 0; i < items.length; i++) {
            player.sendMessage("§7Slot §f" + items[i].slot + "§7: §f" + items[i].name + " x" + items[i].count);
        }
        console.log("[TEST] world.getInventory: " + items.length + " items en cofre");
    });

    commands.register("llenarcofre", "Llena cofre frente a ti con diamantes", function(player, args) {
        var x = Math.floor(player.getX());
        var y = Math.floor(player.getY());
        var z = Math.floor(player.getZ() + 1);
        var inv = world.getInventory(x, y, z);
        if (inv === null) {
            player.sendMessage("§cNo hay un contenedor ahí.");
            return;
        }
        inv.clear();
        var sobrante = inv.addItem("minecraft:diamond", 1000);
        player.sendMessage("§bCofre llenado con diamantes! Sobrante: §f" + sobrante);
        console.log("[TEST] inv.addItem diamonds sobrante=" + sobrante);
    });

    commands.register("vaciarcofre", "Vacía el cofre frente a ti", function(player, args) {
        var x = Math.floor(player.getX());
        var y = Math.floor(player.getY());
        var z = Math.floor(player.getZ() + 1);
        var inv = world.getInventory(x, y, z);
        if (inv === null) {
            player.sendMessage("§cNo hay un contenedor ahí.");
            return;
        }
        inv.clear();
        player.sendMessage("§aCofre vaciado!");
        console.log("[TEST] cofre clear()");
    });

    commands.register("tipocofre", "Muestra el tipo del contenedor frente a ti", function(player, args) {
        var x = Math.floor(player.getX());
        var y = Math.floor(player.getY());
        var z = Math.floor(player.getZ() + 1);
        var inv = world.getInventory(x, y, z);
        if (inv === null) {
            player.sendMessage("§cNo hay un contenedor en §f" + x + "," + y + "," + z);
            return;
        }
        var tipo = inv.getType();
        player.sendMessage("§eTipo de contenedor: §f" + tipo);
        player.sendMessage("§eSlots: §f" + inv.getSize());
        player.sendTitle("§e" + tipo, "§f" + inv.getSize() + " slots");
        console.log("[TEST] inv.getType()=" + tipo + " size=" + inv.getSize());
    });

    commands.register("tipoinv", "Muestra el tipo de tu inventario", function(player, args) {
        var inv = player.getInventory();
        player.sendMessage("§eTipo: §f" + inv.getType() + " §e(§f" + inv.getSize() + " slots§e)");
        console.log("[TEST] player inv.getType()=" + inv.getType());
    });

    commands.register("setcontenido", "Llena cofre con contenido predefinido", function(player, args) {
        var x = Math.floor(player.getX());
        var y = Math.floor(player.getY());
        var z = Math.floor(player.getZ() + 1);
        var inv = world.getInventory(x, y, z);
        if (inv === null) {
            player.sendMessage("§cNo hay un contenedor ahí.");
            return;
        }
        inv.setContents([
            { slot: 0, name: "minecraft:diamond", count: 64 },
            { slot: 1, name: "minecraft:gold_ingot", count: 32 },
            { slot: 2, name: "minecraft:iron_ingot", count: 16 },
            { slot: 3, name: "minecraft:emerald", count: 8 },
            { slot: 4, name: "minecraft:netherite_ingot", count: 4 }
        ]);
        player.sendMessage("§aContenido predefinido establecido en el cofre!");
        console.log("[TEST] inv.setContents() 5 items");
    });

    commands.register("getslot", "Ver item en un slot especifico del cofre", function(player, args) {
        if (args.length < 1) {
            player.sendMessage("§cUso: /getslot <slot>");
            return;
        }
        var slot = parseInt(args[0]);
        var x = Math.floor(player.getX());
        var y = Math.floor(player.getY());
        var z = Math.floor(player.getZ() + 1);
        var inv = world.getInventory(x, y, z);
        if (inv === null) {
            player.sendMessage("§cNo hay un contenedor ahí.");
            return;
        }
        var item = inv.getItem(slot);
        if (item === null) {
            player.sendMessage("§7Slot §f" + slot + "§7: §cvacío");
        } else {
            player.sendMessage("§7Slot §f" + slot + "§7: §f" + item.name + " x" + item.count);
        }
        console.log("[TEST] inv.getItem(" + slot + ")=" + JSON.stringify(item));
    });

    commands.register("setslot", "Poner item en un slot del cofre", function(player, args) {
        if (args.length < 3) {
            player.sendMessage("§cUso: /setslot <slot> <nombre> <cantidad>");
            return;
        }
        var slot = parseInt(args[0]);
        var nombre = args[1];
        var cantidad = parseInt(args[2]);
        var x = Math.floor(player.getX());
        var y = Math.floor(player.getY());
        var z = Math.floor(player.getZ() + 1);
        var inv = world.getInventory(x, y, z);
        if (inv === null) {
            player.sendMessage("§cNo hay un contenedor ahí.");
            return;
        }
        var ok = inv.setItem(slot, nombre, cantidad);
        if (ok) {
            player.sendMessage("§aItem §f" + nombre + " x" + cantidad + " §acolocado en slot §f" + slot);
        } else {
            player.sendMessage("§cNo se pudo colocar el item (nombre inválido o slot fuera de rango)");
        }
        console.log("[TEST] inv.setItem(" + slot + ", " + nombre + ", " + cantidad + ")=" + ok);
    });

    // --- Comandos de test Entity API ---

    commands.register("entidades", "Lista todas las entidades del mundo", function(player, args) {
        var entities = world.getEntities();
        player.sendMessage("§eEntidades en el mundo: §f" + entities.length);
        for (var i = 0; i < entities.length; i++) {
            var e = entities[i];
            var info = e.getType() + " §7(" + Math.floor(e.getX()) + "," + Math.floor(e.getY()) + "," + Math.floor(e.getZ()) + ")";
            if (typeof e.getHealth === "function") {
                info += " §cvida=" + e.getHealth().toFixed(1);
            }
            if (typeof e.getName === "function") {
                info += " §a[" + e.getName() + "]";
            }
            player.sendMessage("§7- §f" + info);
        }
        console.log("[TEST] world.getEntities: " + entities.length + " entidades");
    });

    commands.register("entidadescercanas", "Lista entidades en radio de 10 bloques", function(player, args) {
        var radio = args.length > 0 ? parseFloat(args[0]) : 10;
        if (isNaN(radio) || radio < 1) radio = 10;
        var entities = world.getEntitiesInRadius(player.getX(), player.getY(), player.getZ(), radio);
        player.sendMessage("§eEntidades en radio §f" + radio + "§e bloques: §f" + entities.length);
        for (var i = 0; i < entities.length; i++) {
            var e = entities[i];
            var info = e.getType();
            if (typeof e.getHealth === "function") info += " §cvida=" + e.getHealth().toFixed(1);
            if (typeof e.getName === "function") info += " §a[" + e.getName() + "]";
            player.sendMessage("§7- §f" + info);
        }
        player.sendTitle("§eEntidades cercanas", "§f" + entities.length + " en radio " + radio);
        console.log("[TEST] world.getEntitiesInRadius radio=" + radio + ": " + entities.length + " entidades");
    });

    commands.register("removeritem", "Remueve todos los items del suelo cercanos (radio 20)", function(player, args) {
        var entities = world.getEntitiesInRadius(player.getX(), player.getY(), player.getZ(), 20);
        var count = 0;
        for (var i = 0; i < entities.length; i++) {
            if (entities[i].getType() === "minecraft:item") {
                entities[i].remove();
                count++;
            }
        }
        player.sendMessage("§a" + count + " items removidos del suelo.");
        player.sendTitle("§aLimpieza", "§f" + count + " items removidos");
        console.log("[TEST] entity.remove() items: " + count + " removidos");
    });

    commands.register("removertnt", "Remueve todos los TNT del mundo", function(player, args) {
        var entities = world.getEntities();
        var count = 0;
        for (var i = 0; i < entities.length; i++) {
            if (entities[i].getType() === "minecraft:tnt") {
                entities[i].remove();
                count++;
            }
        }
        player.sendMessage("§a" + count + " TNT removidos.");
        console.log("[TEST] entity.remove() tnt: " + count + " removidos");
    });

    commands.register("spawnitem", "Spawnea un item en tu posición", function(player, args) {
        var itemName = args.length > 0 ? args[0] : "minecraft:diamond";
        var count = args.length > 1 ? parseInt(args[1]) : 1;
        world.spawnEntity("item", player.getX(), player.getY() + 1, player.getZ(), { item: itemName, count: count });
        player.sendMessage("§aItem spawneado: §f" + itemName + " x" + count);
        console.log("[TEST] spawnEntity item: " + itemName + " x" + count);
    });

    // --- Comandos de test Server API ---

    commands.register("serverinfo", "Muestra info del servidor", function(player, args) {
        var count = server.getPlayerCount();
        var max = server.getMaxPlayers();
        var name = server.getName();
        player.sendMessage("§ePlugin: §f" + name);
        player.sendMessage("§eJugadores: §f" + count + "§e/§f" + max);
        player.sendTitle("§e" + name, "§f" + count + "/" + max + " jugadores");
        console.log("[TEST] server.getName=" + name + " count=" + count + "/" + max);
    });

    commands.register("serverjugadores", "Lista jugadores via server.getPlayers()", function(player, args) {
        var jugadores = server.getPlayers();
        var count = server.getPlayerCount();
        player.sendMessage("§eJugadores conectados (server API): §f" + count);
        for (var i = 0; i < jugadores.length; i++) {
            var p = jugadores[i];
            player.sendMessage("§7- §f" + p.getName() + " §7x=" + Math.floor(p.getX()) + " y=" + Math.floor(p.getY()) + " z=" + Math.floor(p.getZ()));
        }
        console.log("[TEST] server.getPlayers: " + count + " jugadores");
    });

    commands.register("buscarjugador", "Busca un jugador por nombre via server.getPlayer()", function(player, args) {
        if (args.length < 1) {
            player.sendMessage("§cUso: /buscarjugador <nombre>");
            return;
        }
        var nombre = args[0];
        var target = server.getPlayer(nombre);
        if (target === null) {
            player.sendMessage("§cJugador §f'" + nombre + "§c' no encontrado.");
            console.log("[TEST] server.getPlayer('" + nombre + "') = null");
        } else {
            player.sendMessage("§aJugador encontrado: §f" + target.getName());
            player.sendMessage("§ePos: §f" + Math.floor(target.getX()) + "," + Math.floor(target.getY()) + "," + Math.floor(target.getZ()));
            player.sendMessage("§eVida: §f" + target.getHealth().toFixed(1) + "/" + target.getMaxHealth().toFixed(1));
            console.log("[TEST] server.getPlayer('" + nombre + "') = encontrado");
        }
    });

    commands.register("serverbroadcast", "Broadcast con server API (chat + título)", function(player, args) {
        if (args.length < 1) {
            player.sendMessage("§cUso: /serverbroadcast <mensaje>");
            return;
        }
        var msg = args.join(" ");
        server.broadcast("§6§l[Server] §r§f" + msg);
        server.broadcastTitle("§6Anuncio", "§f" + msg);
        console.log("[TEST] server.broadcast + broadcastTitle: " + msg);
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
