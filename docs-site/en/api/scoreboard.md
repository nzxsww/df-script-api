# Scoreboard API

The Scoreboard API lets you display information on the right side of the player's screen (the Minecraft Bedrock scoreboard). It supports three modes: **static**, **player-managed** (ScoreboardManager), and **auto-updating** (Live Scoreboard).

## Static Scoreboard

The simplest mode. Create a scoreboard, set its lines, and send it to a player.

```js
var sb = scoreboard.create("§6§lMy Server");

// Recommended way — setLines replaces all content at once
sb.setLines([
    "§7Players: 5",
    "§7Map: Lobby",
    "",
    "§fplay.myserver.com"
]);

player.sendScoreboard(sb);
player.removeScoreboard();
```

You can also build it conditionally with `push`:

```js
var sb = scoreboard.create("§6§lMy Server");
var lines = [
    "§7Players: " + server.getPlayerCount(),
    "§7Mode: " + player.getGameMode()
];

if (player.getGameMode() === "creative") {
    lines.push("§bCreative mode active");
}

lines.push("");
lines.push("§fplay.myserver.com");

sb.setLines(lines);
player.sendScoreboard(sb);
```

::: warning
Modifying a scoreboard after sending it **does not update** what the player sees automatically. You must call `player.sendScoreboard(sb)` again, or use a [Live Scoreboard](#live-scoreboard).
:::

### Scoreboard methods

| Method | Returns | Description |
|---|---|---|
| `scoreboard.create(title)` | `Scoreboard` | Creates a new scoreboard with the given title |
| `sb.setLines(array)` | — | **Replaces all content** with the given array. The recommended way for dynamic content |
| `sb.setLine(index, text)` | `boolean` | Sets text on line `index` (0–14). Fills empty lines if index exceeds current length |
| `sb.addLine(text)` | `boolean` | Appends a line at the end. Returns `false` if the 15-line limit is exceeded |
| `sb.removeLine(index)` | `boolean` | Removes the line at the given index |
| `sb.getLines()` | `string[]` | Returns all current lines as an array |
| `sb.getLineCount()` | `number` | Number of current lines |
| `sb.getTitle()` | `string` | Title of the scoreboard |
| `sb.setDescending()` | — | Reverses the visual order of lines |
| `sb.isDescending()` | `boolean` | Whether the scoreboard is in descending order |
| `sb.removePadding()` | — | Removes the automatic space padding added to each line |
| `sb.sendTo(player)` | — | Shortcut: sends this scoreboard to the given player |

### Player methods

| Method | Description |
|---|---|
| `player.sendScoreboard(sb)` | Sends the scoreboard to the player |
| `player.removeScoreboard()` | Removes the scoreboard from the player's screen |

---

## ScoreboardManager

The `ScoreboardManager` tracks which scoreboard is assigned to each player. It's ideal when you have multiple players with different scoreboards, or when you need to know at any time what a player is viewing.

Each plugin has its **own isolated manager** — `scoreboard.getManager()` always returns the same manager for that plugin.

```js
var mgr = scoreboard.getManager();

events.on("PlayerJoin", function(event) {
    var player = event.getPlayer();
    var sb = scoreboard.create("§6My Server");
    sb.setLine(0, "§7Welcome, §f" + player.getName());
    sb.setLine(1, "§7Players: §f" + server.getPlayerCount());
    mgr.set(player, sb);
});

events.on("PlayerQuit", function(event) {
    mgr.remove(event.getPlayer());
});
```

### ScoreboardManager methods

| Method | Returns | Description |
|---|---|---|
| `scoreboard.getManager()` | `ScoreboardManager` | Returns the plugin's manager |
| `mgr.set(player, sb)` | `boolean` | Assigns and immediately sends the scoreboard to the player |
| `mgr.remove(player)` | — | Removes the assigned scoreboard from the player |
| `mgr.get(player)` | `Scoreboard \| null` | Returns the player's active scoreboardWrapper, or `null` |
| `mgr.hasScoreboard(player)` | `boolean` | `true` if the player has an assigned scoreboard |
| `mgr.getAssignedCount()` | `number` | Number of players with an assigned scoreboard |
| `mgr.clearAll()` | — | Unregisters all players from the manager |

---

## Live Scoreboard

The Live Scoreboard **auto-updates** on a defined interval. On each tick, the callback receives a fresh scoreboard to fill in, which is then sent to all assigned players.

```js
var live = scoreboard.createLive("§aMy Server", function(sb, player) {
    var lines = [
        "§7Player: §f" + player.getName(),
        "§7Health: §c" + Math.floor(player.getHealth()) + "§7/§c" + Math.floor(player.getMaxHealth()),
        "§7Coords: §b" + Math.floor(player.getX()) + " §7/ §b" + Math.floor(player.getY()) + " §7/ §b" + Math.floor(player.getZ()),
        ""
    ];

    var gm = player.getGameMode();
    if (gm === "creative") {
        lines.push("§bMode: §3Creative");
        lines.push("§7Flying: " + (player.isFlying() ? "§aYes" : "§cNo"));
    } else if (gm === "survival") {
        lines.push("§aMode: §2Survival");
        lines.push("§7Exp: §a" + player.getExperienceLevel() + " §7levels");
    } else {
        lines.push("§7Mode: §f" + gm);
    }

    lines.push("");
    lines.push("§7play.myserver.com");
    sb.setLines(lines);
}, 2000); // updates every 2 seconds

events.on("PlayerJoin", function(event) {
    live.addPlayer(event.getPlayer());
});

events.on("PlayerQuit", function(event) {
    live.removePlayer(event.getPlayer());
});

function onDisable() {
    live.stop(); // stop the ticker when the plugin shuts down
}
```

::: tip
The minimum interval is **50ms**. For HUDs with smooth coordinates, **100ms** is recommended. For static lobby scoreboards, **1000–5000ms** is enough.
:::

::: warning
The live callback receives `(sb, playerData)` — the scoreboard and an object with precomputed getters (`getName`, `getHealth`, `getX`, `getGameMode`, etc.). It is called **once per player** registered in the live. **Do not use `server.getPlayers()`** inside the callback — it will cause a transaction deadlock and a server panic.
:::

### Live Scoreboard methods

| Method | Returns | Description |
|---|---|---|
| `scoreboard.createLive(title, fn, intervalMs)` | `LiveScoreboard` | Creates and starts a live scoreboard |
| `live.addPlayer(player)` | `boolean` | Adds the player — they will start receiving automatic updates |
| `live.removePlayer(player)` | `boolean` | Removes the player from the live and removes their scoreboard |
| `live.hasPlayer(player)` | `boolean` | `true` if the player is receiving updates |
| `live.getPlayerCount()` | `number` | Number of players in the live |
| `live.clearPlayers()` | — | Removes all players without stopping the live |
| `live.stop()` | — | Stops the auto-update and removes the scoreboard from all players |

---

## Full Example

A lobby plugin with a real-time scoreboard for all players:

```js
config.setDefaults({
    "title": "§6§lMy Server",
    "interval": 2000
});

var live;

function onEnable() {
    var title = config.getString("title", "§6§lMy Server");
    var interval = config.getInt("interval", 2000);

    live = scoreboard.createLive(title, function(sb) {
        sb.setLine(0, "§7Players: §f" + server.getPlayerCount() + "§7/§f" + server.getMaxPlayers());
        sb.setLine(1, "§7Mode: §eLobby");
        sb.setLine(2, "");
        sb.setLine(3, "§7play.myserver.com");
    }, interval);

    // Add already connected players
    var players = server.getPlayers();
    for (var i = 0; i < players.length; i++) {
        live.addPlayer(players[i]);
    }

    events.on("PlayerJoin", function(event) {
        live.addPlayer(event.getPlayer());
    });

    events.on("PlayerQuit", function(event) {
        live.removePlayer(event.getPlayer());
    });
}

function onDisable() {
    if (live) live.stop();
    config.save();
}

module = { onEnable: onEnable, onDisable: onDisable };
```
