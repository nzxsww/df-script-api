# Full Plugin Example

This is the example plugin included in `plugins/example/`. It was verified with a real manual test on Minecraft Bedrock 1.21.130 - 1.21.132 and demonstrates all available systems.

## Files

### `plugin.yml`

```yaml
name: ExamplePlugin
version: 1.0.0
author: nzxsww
description: Example plugin for Dragonfly Script API
main: index.js
api-version: 1.0.0
```

### `index.js`

```js
console.log("Loading plugin: " + plugin.name + " v" + plugin.version);

config.setDefaults({
    "welcome-message": "Welcome to the server!",
    "prefix":          "§a[Server]§r",
    "protected-y":     10,
    "chat-format":     "§7[§f{name}§7] §r{msg}"
});

function onEnable() {
    console.log("Plugin enabled!");

    var welcomeMsg  = config.getString("welcome-message", "Welcome!");
    var prefix      = config.getString("prefix", "§f[Server]§r");
    var protectedY  = config.getInt("protected-y", 10);
    var chatFormat  = config.getString("chat-format", "§7[§f{name}§7] §r{msg}");

    events.on("PlayerJoin", function(event) {
        var player = event.getPlayer();
        player.sendMessage(prefix + " " + welcomeMsg);
        player.playSound("click");
        event.setJoinMessage("§e" + player.getName() + " §fjoined the server.");
    });

    events.on("PlayerQuit", function(event) {
        var player = event.getPlayer();
        event.setQuitMessage("§c" + player.getName() + " §fleft the server.");
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
            event.getPlayer().sendMessage("§cCan't break blocks here!");
            event.getPlayer().playSound("deny");
        }
    });

    events.on("BlockPlace", function(event) {
        if (event.getBlockY() < protectedY) {
            event.setCancelled(true);
            event.getPlayer().sendMessage("§cCan't place blocks here!");
            event.getPlayer().playSound("deny");
        }
    });

    commands.register("hello", "Greets you by name", function(player, args) {
        player.sendMessage("§aHello, " + player.getName() + "!");
        player.playSound("pop");
    });

    commands.register("pos", "Shows your current position", function(player, args) {
        var x = Math.floor(player.getX());
        var y = Math.floor(player.getY());
        var z = Math.floor(player.getZ());
        player.sendMessage("§7Position: §f" + x + ", " + y + ", " + z);
    });

    commands.register("health", "Shows your current health", function(player, args) {
        var hp    = Math.floor(player.getHealth());
        var maxHp = Math.floor(player.getMaxHealth());
        player.sendMessage("§cHealth: §f" + hp + "/" + maxHp);
    });

    commands.register("give", "Give an item. Usage: /give <item> [count]", function(player, args) {
        if (args.length < 1) {
            player.sendMessage("§cUsage: /give <item> [count]");
            return;
        }
        var item  = args[0];
        var count = args.length >= 2 ? parseInt(args[1]) : 1;
        if (isNaN(count) || count < 1) count = 1;
        if (count > 64) count = 64;

        var ok = player.giveItem(item, count);
        if (ok) {
            player.sendMessage("§aGiven §f" + count + "x §a" + item);
            player.playSound("pop");
        } else {
            player.sendMessage("§cUnknown item or full inventory: " + item);
            player.playSound("deny");
        }
    });

    commands.register("creative", "Switch to creative mode", function(player, args) {
        player.setGameMode("creative");
        player.sendMessage("§aCreative mode enabled.");
        player.playSound("levelup");
    });

    commands.register("survival", "Switch to survival mode", function(player, args) {
        player.setGameMode("survival");
        player.sendMessage("§aSurvival mode enabled.");
    });

    commands.register("fly", "Toggle flight", function(player, args) {
        var mode = player.getGameMode();
        if (mode === "survival" || mode === "adventure") {
            player.setGameMode("creative");
            player.startFlying();
            player.sendMessage("§aFlight enabled. §7(creative mode)");
        } else if (player.isFlying()) {
            player.stopFlying();
            player.sendMessage("§cFlight disabled.");
        } else {
            player.startFlying();
            player.sendMessage("§aFlight enabled.");
        }
    });

    setTimeout(function() {
        console.log("Plugin has been active for 30 seconds.");
    }, 30000);
}

function onDisable() {
    config.save();
    console.log("Plugin disabled — config saved.");
}

module = {
    onEnable:  onEnable,
    onDisable: onDisable
};
```
