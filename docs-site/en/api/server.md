# Server Object

The `server` object is available globally in every plugin. It allows interacting with the server at a global level: getting players from all worlds, sending messages to everyone, and querying server information. It is the equivalent of Bukkit's `Server` object.

> **Difference from `world`:** The `world` object interacts with a specific world (blocks, entities, particles). The `server` object operates at the global server level (all players, regardless of which world they are in).

## Players

| Method | Returns | Description |
|---|---|---|
| `getPlayers()` | `player[]` | All players connected to the server |
| `getPlayerCount()` | `number` | Number of connected players |
| `getMaxPlayers()` | `number` | Maximum player limit configured |
| `getPlayer(name)` | `player\|null` | Find a player by name. Returns `null` if not online |
| `getPlayerByXUID(xuid)` | `player\|null` | Find a player by Xbox XUID. Returns `null` if not online |

```js
// Player count
var count = server.getPlayerCount();
var max = server.getMaxPlayers();
console.log("Players: " + count + "/" + max);

// Find a specific player
var player = server.getPlayer("Notch");
if (player !== null) {
    player.sendMessage("§aHello Notch!");
}

// Iterate over all players
var players = server.getPlayers();
for (var i = 0; i < players.length; i++) {
    players[i].sendMessage("§eMessage for everyone!");
}
```

## Messages

| Method | Description |
|---|---|
| `broadcast(msg)` | Send a chat message to all connected players |
| `broadcastTitle(text, subtitle)` | Send a large title to all connected players |

```js
server.broadcast("§a[Announcement] §fThe event starts in 1 minute!");
server.broadcastTitle("§c§lEVENT!", "§7Get ready to play");
```

## Server Information

| Method | Returns | Description |
|---|---|---|
| `getName()` | `string` | Name of the plugin making the call |
| `shutdown()` | — | Shut down the server gracefully |

```js
console.log("Plugin: " + server.getName());

// Shut down the server (use with caution)
server.shutdown();
```
